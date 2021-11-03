package entity

import (
	"yumi/pkg/status"
)

var (
	m = status.NewMessages()

	// PasswordIncorrect 密码错误
	PasswordIncorrect = m.NewMessageID("", "", "")
	// UserFmtIncorrect 用户名格式不正确
	UserFmtIncorrect = m.NewMessageID("", "", "")
	// PasswordFmtIncorrect 密码格式不正确
	PasswordFmtIncorrect = m.NewMessageID("", "", "")
	// UserNotFound 用户名不存在
	UserNotFound = m.NewMessageID("", "", "")
	// UserAuthenticationExpired 用户身份认证已过期
	UserAuthenticationExpired = m.NewMessageID("The user identity authentication has expired", "用户身份认证已过期", "")
)

func init() {
	m.InitI18N()
}
