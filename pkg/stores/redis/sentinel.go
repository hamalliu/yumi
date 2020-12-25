package redis

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Sentinel provides a way to add high availability (HA) to Redis Pool using
// preconfigured addresses of Sentinel servers and name of master which Sentinels
// monitor. It works with Redis >= 2.8.12 (mostly because of ROLE command that
// was introduced in that version, it's possible though to support old versions
// using INFO command).
type Sentinel struct {
	pool *PoolMulAddr

	// MasterName is a name of Redis master Sentinel servers monitor.
	masterName string

	masterPool string
	// type is map[string]int32
	slaveAddrs sync.Map

	muxtex sync.Mutex
}

const (
	poolConnStatusAction = iota
	poolConnStatusFailure
)

type poolConn struct {
	status  int
	addr    string
	created time.Time
	t       time.Time
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
func NewSentinel(masterName string) *Sentinel {
	s := &Sentinel{masterName: masterName}

	return s
}

// NewMasterPool returns an pool of current Redis master instance.
func (s *Sentinel) NewMasterPool() (pool *redis.Pool, err error) {
	// TODO:
	return
}

// NewSlavePool returns a pool of current Redis slave instance.
func (s *Sentinel) NewSlavePool() (pool *redis.Pool, err error) {
	// TODO:
	return
}

// Discover watch server and update status
func (s *Sentinel) Discover() {

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
