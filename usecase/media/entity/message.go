package entity

import "yumi/pkg/status"

var (
	m = status.NewMessages()

	// FileIsTooLarge 文件太大
	FileIsTooLarge = m.NewMessageID("", "", "")
)

func init() {
	m.InitI18N()
}
