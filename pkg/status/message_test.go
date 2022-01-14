package status

import (
	"testing"
)

func TestMessage(t *testing.T) {

	var (
		m = NewMessages()

		// FileIsTooLarge 文件太大
		FileIsTooLarge = m.NewMessageID("File is too large, size: {{.Size}} {{.Count}}", "文件太大, size: {{.Size}}", "文件太大, size: {{.Size}}")
	)

	if err := m.InitI18N(); err != nil {
		t.Error(err)
		return
	}
	FileIsTooLarge.SetArgs(Args{"Size": "100MB", "Count": 100})

	t.Log(FileIsTooLarge.T("en-US"))
	t.Log(FileIsTooLarge.T("zh-CN"))
	t.Log(FileIsTooLarge.T("zh-TW"))
}
