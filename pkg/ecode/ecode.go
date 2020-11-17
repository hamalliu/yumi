package ecode

import (
	"errors"
	"fmt"
	"strconv"

	"yumi/pkg/log"
)
var (
	//OK 正确
	OK = add(0, "正确") 

	//serverErr 服务器错误
	serverErr = add(500, "服务器错误") 
	//paramsErr 前端请求参数错误
	paramsErr = add(501, "请求参数错误") 
)

var (
	_codes = make(map[Code]string) // register codes.
)

func add(e int, msg string) Code {
	if _, ok := _codes[Code(e)]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[Code(e)] = msg
	return Int(e)
}

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

// A Code is an int error code spec.
type Code int

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int { return int(e) }

// Message return error message
func (e Code) Message() string {
	if msg, ok := _codes[e]; ok {
		return msg
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
