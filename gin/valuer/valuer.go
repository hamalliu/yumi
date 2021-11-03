package valuer

import (
	"fmt"
)

//Key ...
type Key string

const (
	//KeyUser ...
	KeyUser      Key = "user"
	//KeySecret ...
	KeySecret    Key = "secret"
	//KeyTimestamp ...
	KeyTimestamp Key = "timestamp"
	//KeyNonce ...
	KeyNonce     Key = "nonce"
	//KeySignature ...
	KeySignature Key = "signature"
)

//Valuer ...
type Valuer struct {
	Value interface{}
}

//User ...
type User struct {
	UserID   string
	UserName string
}

//User ...
func (v Valuer) User() User {
	val, ok := v.Value.(User)
	if ok {
		return val
	}

	return User{}
}

//String ...
func (v Valuer) String() string {
	val, ok := v.Value.(string)
	if ok {
		return val
	}
	panic(fmt.Errorf("%v is not string", v.Value))
}

//Bytes ...
func (v Valuer) Bytes() []byte {
	val, ok := v.Value.([]byte)
	if ok {
		return val
	}
	panic(fmt.Errorf("%v is not bytes", v.Value))
}

//Float64 ...
func (v Valuer) Float64() float64 {
	val, ok := v.Value.(float64)
	if ok {
		return val
	}
	panic(fmt.Errorf("%v is not float64", v.Value))
}

//Float32 ...
func (v Valuer) Float32() float32 {
	val, ok := v.Value.(float32)
	if ok {
		return val
	}
	panic(fmt.Errorf("%v is not float32", v.Value))
}

//Int ...
func (v Valuer) Int() int {
	val, ok := v.Value.(int)
	if ok {
		return val
	}
	panic(fmt.Errorf("%v is not int", v.Value))
}

//Int64 ...
func (v Valuer) Int64() int64 {
	val, ok := v.Value.(int64)
	if ok {
		return val
	}
	panic(fmt.Errorf("%v is not int64", v.Value))
}
