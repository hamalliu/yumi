package redis

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gomodule/redigo/redis"

	"yumi/pkg/log"
)

// PoolMulAddr containing multiple addresses
type PoolMulAddr struct {
	// Addrs is a slice with known Sentinel addresses.
	addrs map[string]Config
	// type is map[string]int32
	actionAddrs sync.Map
	// type is map[redis.Conn]*poolConn
	conns sync.Map
	// type is map[string][]redis.Conn
	actionAddrConns sync.Map

	muxtex sync.Mutex
}

type poolConn struct {
	status  int
	addr    string
	created time.Time
	t       time.Time
}

// NewPoolMulAddr create a PoolMulAddr object
func NewPoolMulAddr(confs []Config) *PoolMulAddr {
	pma := &PoolMulAddr{}

	for i := range confs {
		var zero int32
		pma.addrs[confs[i].Addr] = confs[i]
		pma.actionAddrs.Store(confs[i].Addr, &zero)
	}

	return pma
}

// New creat a pool containing multiple addresses.
func (pma *PoolMulAddr) New(maxIdle, maxActive, idleTimeoutSecond int, wait bool) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		Wait:        wait,
		IdleTimeout: time.Duration(idleTimeoutSecond) * time.Second,
		Dial: func() (c redis.Conn, err error) {
			pma.muxtex.Lock()
			defer pma.muxtex.Unlock()

			addr := pma.nextAddr()
			// 重试次数
			retryTimes := 3
			for i := 0; i < retryTimes; i++ {
				c, err = redis.Dial("tcp", addr)
				if err != nil {
					log.Warning("redis: %s", err.Error())
				} else {
					break
				}
			}
			if err != nil {
				// 认为addr宕掉，标记状态
				pma.actionAddrConns.Range(func(key, value interface{}) bool {
					if _, ok := pma.conns.Load(value); ok {
						pma.conns.Store(value, poolConnStatusFailure)
					}
					return true
				})
				pma.actionAddrConns.Delete(addr)
				pma.actionAddrs.Delete(addr)
			} else {
				count, _ := pma.actionAddrs.Load(addr)
				atomic.AddInt32(count.(*int32), 1)

				pma.conns.Store(c, &poolConn{
					status:  poolConnStatusAction,
					addr:    addr,
					created: time.Now(),
				})

				conns, _ := pma.actionAddrConns.Load(addr)
				if conns == nil {
					conns = []redis.Conn{}
				}
				conns = append(conns.([]redis.Conn), c)
			}
			return
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			v, ok := pma.conns.Load(c)
			if !ok {
				log.Warning(fmt.Errorf("redis: redis conn leak"))
			}
			pc := v.(*poolConn)

			if pc.status == poolConnStatusFailure {
				pma.conns.Delete(c)
				count, _ := pma.actionAddrs.Load(pc.addr)
				atomic.AddInt32(count.(*int32), -1)
				return errors.New("redis: conn failure")
			}

			_, err := c.Do("PING")
			if err != nil {
				pma.conns.Delete(c)
				count, _ := pma.actionAddrs.Load(pc.addr)
				atomic.AddInt32(count.(*int32), -1)
			}

			pc.t = t
			return err
		},
	}

	return pool, nil
}

func (pma *PoolMulAddr) connGC() {
	// TODO:
	// pma.conns.Range()
}

func (pma *PoolMulAddr) nextAddr() string {
	nextAddr := ""
	var min int32
	pma.actionAddrs.Range(func(key, value interface{}) bool {
		n := atomic.LoadInt32(value.(*int32))
		addr := key.(string)
		if n == 0 {
			nextAddr = addr
		}
		if min == 0 || min > n {
			nextAddr = addr
			min = n
		}
		return false
	})

	return nextAddr
}
