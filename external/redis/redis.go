package redis

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Model struct {
	redis.Conn
	conf Config
}

func New(conf Config) (*Model, error) {
	var (
		err   error
		cache = new(Model)
	)

	cache.conf = conf
	if err = cache.dialRedis(); err != nil {
		return nil, err
	}

	return cache, nil
}

func (m *Model) dialRedis() error {
	var (
		err error
	)

	m.Conn, err = redis.Dial("tcp",
		m.conf.Host+":"+m.conf.Port,
		[]redis.DialOption{redis.DialReadTimeout(10 * time.Second), redis.DialWriteTimeout(10 * time.Second)}...)
	if err != nil {
		return errors.New("redis连接失败：" + err.Error())
	}

	return nil
}

func (m *Model) PutStructOrMapFlat(key string, obj interface{}) error {
	if _, err := m.Do("HMSET", redis.Args{}.Add(key).AddFlat(obj)...); err != nil {
		return err
	}

	return nil
}

func (m *Model) GetStructOrMapFlat(key string, objptr interface{}) error {
	v, err := redis.Values(m.Do("HGETALL", key))
	if err != nil {
		return err
	}

	if err := redis.ScanStruct(v, objptr); err != nil {
		return err
	}

	return err
}

func (m *Model) MputStructOrMapFlat(key string, obj interface{}) error {
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
		if _, err := m.Do("HMSET", redis.Args{}.Add(key+fmt.Sprintf("%d", i)).AddFlat(ov.Index(i).Addr().Interface())...); err != nil {
			return err
		}

		if _, err := m.Do("LPUSH", key, i); err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) MgetStructOrMapFlat(key string, objptr interface{}) error {
	args := []interface{}{key, "BY", key}

	ov := reflect.ValueOf(objptr).Elem().Elem()
	ot := reflect.ValueOf(objptr)
	ot.Field(0)
	l := ov.NumField()
	for i := 0; i < l; i++ {
		args = append(args, "GET", key+"*->"+ov.Type().Field(i).Tag.Get("redis"))
	}
	v, err := redis.Values(m.Do("SORT", args...))
	if err != nil {
		return err
	}

	if err := redis.ScanSlice(v, objptr); err != nil {
		return err
	}

	return err
}

func (m *Model) PutString(key string, v string, ex int64) error {
	var args []interface{}
	if ex == 0 {
		args = append(args, key, v)
	} else {
		args = append(args, key, v, "EX", ex)
	}
	if _, err := m.Do("SET", args...); err != nil {
		return err
	}

	return nil
}

func (m *Model) GetString(key string) (string, error) {
	return redis.String(m.Do("GET", key))

}

func (m *Model) MputString(args []interface{}) error {
	_, err := m.Do("MSET", args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) MgetString(key []interface{}) ([]string, error) {
	return redis.Strings(m.Do("MGET", key...))

}

func (m *Model) GetFloat64(key string) (float64, error) {
	return redis.Float64(m.Do("GET", key))

}

func (m *Model) Del(key []interface{}) error {
	if _, err := m.Do("DEL", key...); err != nil {
		return err
	}

	return nil
}

func (m *Model) DelMap(key string, mkeys []interface{}) error {
	var args []interface{}
	args = append(args, key)
	args = append(args, mkeys...)
	if _, err := m.Do("HDEL", args...); err != nil {
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
