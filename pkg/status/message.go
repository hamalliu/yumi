package status

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n"
)

type MessageID struct {
	id   string
	args []interface{}

	m *Messages
}

func (msg MessageID) T(language string) string {
	tmsg, _ := msg.m.T(language, msg.id, msg.args...)
	return tmsg
}

type I18nMessageID interface {
	T(language string) string
}

type Message struct {
	id   string
	enUS string
	zhCN string
	zhTW string
}

type Messages struct {
	msgs map[string]Message
}

func NewMessages() *Messages {
	m := Messages{}
	m.msgs = make(map[string]Message)
	return &Messages{}
}

func (m *Messages) NewMessageID(enUS, zhCN, zhTW string) MessageID {
	msg := Message{}
	msg.id = fmt.Sprintf("%p", &msg)
	m.msgs[msg.id] = msg

	return MessageID{id: msg.id, m: m}
}

func (m *Messages) T(language string, id string, args ...interface{}) (string, error) {
	T, err := i18n.Tfunc(language)
	if err != nil {
		return "", err
	}

	return T(id, args), nil
}

func (m *Messages) InitI18N() error {
	// TODO:

	return nil
}
