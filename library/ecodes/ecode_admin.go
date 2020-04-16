package ecodes

/**
 *admin模块错误码（10000-11000）
 */

var (
	UserError     = NewStatus(10000, "用户名错误")
	PasswordError = NewStatus(10001, "密码错误")
)
