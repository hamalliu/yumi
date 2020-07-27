package ecode

import (
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"

	"yumi/pkg/log"
)

var (
	_messages atomic.Value         // NOTE: stored map[string]map[int]string
	_codes    = map[int]struct{}{} // register codes.
)

// Register register ecode message map.
func Register(cm map[int]string) {
	_messages.Store(cm)
}

// New new a ecode.Codes by int value.
// NOTE: ecode must unique in global, the New will check repeat and then panic.
func New(e int) Code {
	if e <= 0 {
		panic("business ecode must greater than zero")
	}
	return add(e)
}

func add(e int) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return Int(e)
}

// Codes ecode error sinterface which has a code & message.
type Codes interface {
	// sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	//Detail get error detail,it may be nil.
	Details() []interface{}
	// Equal for compatible.
	// Deprecated: please use ecode.EqualError.
	Equal(error) bool
}

// A Code is an int error code spec.
type Code int

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int { return int(e) }

//ParamsErrMsg 参数错误消息内容
func (e Code) ParamsErrMsg() string {
	return "请求参数错误"
}

// Message return error message
func (e Code) Message() string {
	if e == paramsErr {
		return e.ParamsErrMsg()
	}

	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[e.Code()]; ok {
			return msg
		}
	}
	return e.Error()
}

//ParamsErr 参数错误
func (e Code) ParamsErr(err error) string {
	if e != paramsErr {
		return ""
	}
	return err.Error()
}

// Details return details.
func (e Code) Details() []interface{} { return nil }

// Int parse code int to error.
func Int(i int) Code { return Code(i) }

// String parse code string to error.
func String(e string) Code {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return serverErr
	}
	return Code(i)
}

//Must 转换err为code，如果失败且不是参数错误就panic
func Must(err error) Code {
	if err == nil {
		return OK
	}

	c, ok := err.(Code)
	if ok {
		return c
	}
	if errors.As(err, emptyParamsErr) {
		return paramsErr
	}
	panic(err)
}

//ServerErr 服务器错误
func ServerErr(err error) Code {
	if err == nil {
		return serverErr
	}
	if c, ok := err.(Code); ok {
		return c
	}
	log.Error2(err)
	return serverErr
}

//ErrorTypeParamsErr 用于前端未按要求传参的错误类型
type ErrorTypeParamsErr error

var emptyParamsErr = new(ErrorTypeParamsErr)

//ParamsErr 参数错误
func ParamsErr(err error) ErrorTypeParamsErr {
	return ErrorTypeParamsErr(err)
}
