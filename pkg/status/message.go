package status

import (
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
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

type Args map[string]interface{}

func (msg *MessageID) SetArgs(args Args) {
	msg.args = append(msg.args, args)
}

func (msg *MessageID) SetArg(key string, value interface{}) {
	msg.args = append(msg.args, map[string]interface{}{key: value})
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
	b    *bundle.Bundle
	msgs map[string]Message
}

func NewMessages() *Messages {
	m := Messages{}
	m.msgs = make(map[string]Message)
	return &m
}

func (m *Messages) NewMessageID(enUS, zhCN, zhTW string) MessageID {
	msg := Message{}
	// FIXME: id生成需优化
	msg.id = fmt.Sprintf("%p", &msg)
	msg.enUS = enUS
	msg.zhCN = zhCN
	msg.zhTW = zhTW
	m.msgs[msg.id] = msg

	return MessageID{id: msg.id, m: m}
}

func (m *Messages) T(language string, id string, args ...interface{}) (string, error) {
	T, err := m.b.Tfunc(language)
	if err != nil {
		return "", err
	}
	return T(id, args...), nil
}

type Translation struct {
	ID          string `json:"id"`
	Translation string `json:"translation"`
}

func (m *Messages) marshalEnUs() ([]byte, error) {
	var ts []Translation
	for k, v := range m.msgs {
		ts = append(ts, Translation{ID: k, Translation: v.enUS})
	}
	return json.Marshal(&ts)
}

func (m *Messages) marshalZhCN() ([]byte, error) {
	var ts []Translation
	for k, v := range m.msgs {
		ts = append(ts, Translation{ID: k, Translation: v.zhCN})
	}
	return json.Marshal(&ts)
}

func (m *Messages) marshalZhTW() ([]byte, error) {
	var ts []Translation
	for k, v := range m.msgs {
		ts = append(ts, Translation{ID: k, Translation: v.zhTW})
	}
	return json.Marshal(&ts)
}

func (m *Messages) InitI18N() error {
	m.b = bundle.New()

	enUs, err := m.marshalEnUs()
	if err != nil {
		return err
	}
	zhCn, err := m.marshalZhCN()
	if err != nil {
		return err
	}
	zhTw, err := m.marshalZhTW()
	if err != nil {
		return err
	}

	if err := m.b.ParseTranslationFileBytes("en-US.all.json", enUs); err != nil {
		return err
	}
	if err := m.b.ParseTranslationFileBytes("zh-CN.all.json", zhCn); err != nil {
		return err
	}
	if err := m.b.ParseTranslationFileBytes("zh-TW.all.json", zhTw); err != nil {
		return err
	}

	return nil
}
