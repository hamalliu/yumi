package ecode

/**
 *admin模块错误码（10000-11000）
 */

var (
	//UserError 用户名错误
	UserError     = add(10000, "用户名错误") 
	//PasswordError 密码错误
	PasswordError = add(10001, "密码错误")
)
