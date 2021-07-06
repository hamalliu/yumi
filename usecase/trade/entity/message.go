package entity

import (
	"yumi/pkg/status"
)

var (
	m = status.NewMessages()

	// 订单超时
	OrderTimeout = m.NewMessageID("", "", "")
	// 订单已存在
	OrderAlreadyExists = m.NewMessageID("", "", "")
	// 订单已完成，不能取消
	OrderFinishedRefuseCancel = m.NewMessageID("", "", "")
)

func init() {
	m.InitI18N()
}
