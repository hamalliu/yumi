package ecode

var (
	//OK 正确
	OK = add(0) 

	//serverErr 服务器错误
	serverErr = add(500) 
	//paramsErr 前端请求参数错误
	paramsErr = add(501) 
)
