package entity

import (
	"yumi/pkg/status"
)

var (
	m = status.NewMessages()

	// OrderTimeout 订单超时
	OrderTimeout = m.NewMessageID("", "", "")
	// OrderAlreadyExists 订单已存在
	OrderAlreadyExists = m.NewMessageID("", "", "")
	// OrderFinishedRefuseCancel 订单已完成，不能取消
	OrderFinishedRefuseCancel = m.NewMessageID("", "", "")
)

func init() {
	m.InitI18N()
}
