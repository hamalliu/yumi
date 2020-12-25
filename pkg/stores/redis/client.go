package redis

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gomodule/redigo/redis"
)

// NewClient create a redis client
func NewClient(pool *redis.Pool) *Client {
	return &Client{
		p: pool,
	}
}

// NewRWSplittingClient create a redis client with R/W splitting
func NewRWSplittingClient(readPool, writePool *redis.Pool) *Client {
	return &Client{
		read:  readPool,
		write: writePool,
		rwSplitting: true,
	}
}

// Client ...
type Client struct {
	rwSplitting bool
	p           *redis.Pool
	read        *redis.Pool
	write       *redis.Pool
}

// Get return a redis connect
func (cli *Client) Get(commond string) *Conn {
	p := cli.p
	if cli.rwSplitting {
		if cli.writeCommond(commond) {
			p = cli.write
		} else {
			p = cli.read
		}
	}
	return &Conn{commond: commond, rc: p.Get()}
}

// Conn warp a redis connect
type Conn struct {
	commond string
	rc      redis.Conn
}

// Do warp a Do fucntion of redis connect
func (c *Conn) Do(args ...interface{}) (reply interface{}, err error) {
	defer c.rc.Close()
	return c.rc.Do(c.commond, args)
}

//PutStructOrMapFlat ...
func (cli *Client) PutStructOrMapFlat(key string, obj interface{}) error {
	if _, err := cli.Get("HMSET").Do(redis.Args{}.Add(key).AddFlat(obj)...); err != nil {
		return err
	}

	return nil
}

//GetStructOrMapFlat ...
func (cli *Client) GetStructOrMapFlat(key string, objptr interface{}) error {
	v, err := redis.Values(cli.Get("HGETALL").Do(key))
	if err != nil {
		return err
	}

	if err := redis.ScanStruct(v, objptr); err != nil {
		return err
	}

	return err
}

//MputStructOrMapFlat ...
func (cli *Client) MputStructOrMapFlat(key string, obj interface{}) error {
	ov := reflect.ValueOf(obj)
	switch ov.Kind() {
	case reflect.Array, reflect.Slice:

	case reflect.Ptr:
		ov = ov.Elem()

	default:
		err := errors.New("obj类型错误")
		return err
	}

	l := ov.Len()
	for i := 0; i < l; i++ {
		if _, err := cli.Get("HMSET").Do(redis.Args{}.Add(key + fmt.Sprintf("%d", i)).AddFlat(ov.Index(i).Addr().Interface())...); err != nil {
			return err
		}

		if _, err := cli.Get("LPUSH").Do(key, i); err != nil {
			return err
		}
	}

	return nil
}

//MgetStructOrMapFlat ...
func (cli *Client) MgetStructOrMapFlat(key string, objptr interface{}) error {
	args := []interface{}{key, "BY", key}

	ov := reflect.ValueOf(objptr).Elem().Elem()
	ot := reflect.ValueOf(objptr)
	ot.Field(0)
	l := ov.NumField()
	for i := 0; i < l; i++ {
		args = append(args, "GET", key+"*->"+ov.Type().Field(i).Tag.Get("redis"))
	}
	v, err := redis.Values(cli.Get("SORT").Do(args...))
	if err != nil {
		return err
	}

	if err := redis.ScanSlice(v, objptr); err != nil {
		return err
	}

	return err
}

//PutString ...
func (cli *Client) PutString(key string, v string, ex int64) error {
	var args []interface{}
	if ex == 0 {
		args = append(args, key, v)
	} else {
		args = append(args, key, v, "EX", ex)
	}
	if _, err := cli.Get("SET").Do(args...); err != nil {
		return err
	}

	return nil
}

//GetString ...
func (cli *Client) GetString(key string) (string, error) {
	return redis.String(cli.Get("GET").Do(key))
}

//MputString ...
func (cli *Client) MputString(args []interface{}) error {
	_, err := cli.Get("MSET").Do(args...)
	if err != nil {
		return err
	}

	return nil
}

//MgetString ...
func (cli *Client) MgetString(key []interface{}) ([]string, error) {
	return redis.Strings(cli.Get("MGET").Do(key...))

}

//GetFloat64 ...
func (cli *Client) GetFloat64(key string) (float64, error) {
	return redis.Float64(cli.Get("GET").Do(key))
}

//Del ...
func (cli *Client) Del(key []interface{}) error {
	if _, err := cli.Get("DEL").Do(key...); err != nil {
		return err
	}

	return nil
}

//DelMap ...
func (cli *Client) DelMap(key string, mkeys []interface{}) error {
	var args []interface{}
	args = append(args, key)
	args = append(args, mkeys...)
	if _, err := cli.Get("HDEL").Do(args...); err != nil {
		return err
	}

	return nil
}

//func putStructCascade(c redis.Conn, key string, objPtr interface{}) error {
//	ot := reflect.TypeOf(objPtr)
//	ov := reflect.ValueOf(objPtr)
//	l := ot.Elem().NumField()
//	for i := 0; i < l; i++ {
//		switch ov.Field(i).Kind() {
//		case reflect.Struct:
//			putStructCascade()
//		case reflect.Slice, reflect.Array:
//
//		case reflect.Ptr:
//
//		case reflect.Int,reflect.String:
//
//		default:
//
//		}
//	}
//
//	if _, err := c.Do("HMSET", redis.Args{}.Add(key).AddFlat(objPtr)...); err != nil {
//		log.Error(err)
//		return err
//	}
//	return nil
//}
//
//func putMapCascade(c redis.Conn, key string, objMap interface{}) error {
//
//	return nil
//}
//
//func (m *Model) PutStructOrMapCascade(key string, obj interface{}) error {
//	ov := reflect.ValueOf(obj)
//	switch ov.Kind() {
//	case reflect.Ptr:
//		if ov.Elem().Kind() != reflect.Struct {
//			err := errors.New("obj指针类型错误")
//			log.Error(err)
//			return err
//		}
//		putStructCascade(m.c, key, obj)
//
//	case reflect.Struct:
//		putStructCascade(m.c, key, ov.Addr().Interface())
//
//	case reflect.Map:
//		putMapCascade(m.c, key, obj)
//
//	default:
//		err := errors.New("obj类型错误")
//		log.Error(err)
//		return err
//
//	}
//
//	return nil
//}
//
//func (m *Model) GetStructOrMapCascade(key string, objptr interface{}) error {
//
//	return nil
//}

func (cli *Client) writeCommond(commond string) bool {
	// TODO:
	return false
}
