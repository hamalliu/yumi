package entity

import (
	"yumi/pkg/status"
)

var (
	m = status.NewMessages()

	// 密码错误
	PasswordIncorrect = m.NewMessageID("", "", "")
	// 用户名格式不正确
	UserFmtIncorrect = m.NewMessageID("", "", "")
	// 密码格式不正确
	PasswordFmtIncorrect = m.NewMessageID("", "", "")
	// 用户名不存在
	UserNotFound = m.NewMessageID("", "", "")
)

func init() {
	m.InitI18N()
}
