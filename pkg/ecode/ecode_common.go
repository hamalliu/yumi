package ecode

var (
	OK = add(0) // 正确

	serverErr = add(-500) // 服务器错误
	paramsErr = add(-501) //前端请求参数错误
)
