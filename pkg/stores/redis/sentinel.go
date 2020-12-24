package redis

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
	"yumi/pkg/log"

	"github.com/gomodule/redigo/redis"
)

// Sentinel provides a way to add high availability (HA) to Redis Pool using
// preconfigured addresses of Sentinel servers and name of master which Sentinels
// monitor. It works with Redis >= 2.8.12 (mostly because of ROLE command that
// was introduced in that version, it's possible though to support old versions
// using INFO command).
type Sentinel struct {
	// Addrs is a slice with known Sentinel addresses.
	addrs           []string
	actionAddrs     map[string]int
	connAddr        map[redis.Conn]string
	actionAddrConns map[string][]redis.Conn

	failureAddrConns map[string][]redis.Conn

	// MasterName is a name of Redis master Sentinel servers monitor.
	masterName string

	masterAddr string
	slaveAddrs map[string]int

	muxtex sync.Mutex
}

// NoSentinelsAvailable is returned when all sentinels in the list are exhausted
// (or none configured), and contains the last error returned by Dial (which
// may be nil)
type NoSentinelsAvailable struct {
	lastError error
}

func (ns NoSentinelsAvailable) Error() string {
	if ns.lastError != nil {
		return fmt.Sprintf("redigo: no sentinels available; last error: %s", ns.lastError.Error())
	}
	return fmt.Sprintf("redigo: no sentinels available")
}

// NewSentinel create a sentinel
func NewSentinel(addrs []string, masterName string) *Sentinel {
	s := &Sentinel{
		addrs:            addrs,
		actionAddrs:      make(map[string]int),
		connAddr:   make(map[redis.Conn]string),
		actionAddrConns:  make(map[string][]redis.Conn),
		failureAddrConns: make(map[string][]redis.Conn),
		masterName:       masterName,
		slaveAddrs:       make(map[string]int),
	}
	for _, addr := range addrs {
		s.actionAddrs[addr] = 0
	}
	return s
}

// NewSentinelPool creat a pool of current Redis sentinel instance.
func (s *Sentinel) NewSentinelPool() (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (c redis.Conn, err error) {
			s.muxtex.Lock()
			defer s.muxtex.Unlock()

			addr := s.nextSentinelAddr()
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
				s.failureAddrConns[addr] = append(s.failureAddrConns[addr], s.actionAddrConns[addr]...)
				for _, failconn := range s.failureAddrConns[addr] {
					s.connAddr[failconn] = "failure"
				}
				delete(s.actionAddrConns, addr)
				delete(s.actionAddrs, addr)
			} else {
				s.actionAddrs[addr]++
				s.connAddr[c] = addr
				s.actionAddrConns[addr] = append(s.actionAddrConns[addr], c)
			}
			return
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			addr := s.connAddr[c]
			if addr == "" {
				log.Warning(fmt.Errorf("redis: redis conn leak"))
			} else if addr == "failure" {
				return errors.New("redis: conn failure")
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return pool, nil
}

// NewMasterPool returns an pool of current Redis master instance.
func (s *Sentinel) NewMasterPool() (pool *redis.Pool, err error) {
	return
}

// NewSlavePool returns a pool of current Redis slave instance.
func (s *Sentinel) NewSlavePool() (pool *redis.Pool, err error) {
	return
}

// Discover watch server and update status
func (s *Sentinel) Discover() {

}

func (s *Sentinel) nextSentinelAddr() string {
	nextAddr := ""
	min := 0
	for addr, n := range s.actionAddrs {
		if n == 0 {
			return addr
		}
		if min == 0 || min > n {
			nextAddr = addr
			min = n
		}
	}

	return nextAddr
}

func (s *Sentinel) nextSlaveAddr() string {
	nextAddr := ""
	min := 0
	for addr, n := range s.slaveAddrs {
		if n == 0 {
			return addr
		}
		if min == 0 || min > n {
			nextAddr = addr
			min = n
		}
	}

	return nextAddr
}

// Slave represents a Redis slave instance which is known by Sentinel.
type Slave struct {
	ip    string
	port  string
	flags string
}

// Addr returns an address of slave.
func (s *Slave) Addr() string {
	return net.JoinHostPort(s.ip, s.port)
}

// Available returns if slave is in working state at moment based on information in slave flags.
func (s *Slave) Available() bool {
	return !strings.Contains(s.flags, "disconnected") && !strings.Contains(s.flags, "s_down")
}

// TestRole wraps GetRole in a test to verify if the role matches an expected
// role string. If there was any error in querying the supplied connection,
// the function returns false. Works with Redis >= 2.8.12.
// It's not goroutine safe, but if you call this method on pooled connections
// then you are OK.
func TestRole(c redis.Conn, expectedRole string) bool {
	role, err := getRole(c)
	if err != nil || role != expectedRole {
		return false
	}
	return true
}

// getRole is a convenience function supplied to query an instance (master or
// slave) for its role. It attempts to use the ROLE command introduced in
// redis 2.8.12.
func getRole(c redis.Conn) (string, error) {
	res, err := c.Do("ROLE")
	if err != nil {
		return "", err
	}
	rres, ok := res.([]interface{})
	if ok {
		return redis.String(rres[0], nil)
	}
	return "", errors.New("redigo: can not transform ROLE reply to string")
}

func queryConnectedClients(conn redis.Conn)

func queryForMaster(conn redis.Conn, masterName string) (string, error) {
	res, err := redis.Strings(conn.Do("SENTINEL", "get-master-addr-by-name", masterName))
	if err != nil {
		return "", err
	}
	if len(res) < 2 {
		return "", errors.New("redigo: malformed get-master-addr-by-name reply")
	}
	masterAddr := net.JoinHostPort(res[0], res[1])
	return masterAddr, nil
}

func queryForSlaveAddrs(conn redis.Conn, masterName string) ([]string, error) {
	slaves, err := queryForSlaves(conn, masterName)
	if err != nil {
		return nil, err
	}
	slaveAddrs := make([]string, 0)
	for _, slave := range slaves {
		slaveAddrs = append(slaveAddrs, slave.Addr())
	}
	return slaveAddrs, nil
}

func queryForSlaves(conn redis.Conn, masterName string) ([]*Slave, error) {
	res, err := redis.Values(conn.Do("SENTINEL", "slaves", masterName))
	if err != nil {
		return nil, err
	}
	slaves := make([]*Slave, 0)
	for _, a := range res {
		sm, err := redis.StringMap(a, err)
		if err != nil {
			return slaves, err
		}
		slave := &Slave{
			ip:    sm["ip"],
			port:  sm["port"],
			flags: sm["flags"],
		}
		slaves = append(slaves, slave)
	}
	return slaves, nil
}

func queryForSentinels(conn redis.Conn, masterName string) ([]string, error) {
	res, err := redis.Values(conn.Do("SENTINEL", "sentinels", masterName))
	if err != nil {
		return nil, err
	}
	sentinels := make([]string, 0)
	for _, a := range res {
		sm, err := redis.StringMap(a, err)
		if err != nil {
			return sentinels, err
		}
		sentinels = append(sentinels, fmt.Sprintf("%s:%s", sm["ip"], sm["port"]))
	}
	return sentinels, nil
}

func stringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
