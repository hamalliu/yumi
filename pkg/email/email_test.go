package email

import (
	"testing"
	"time"
)

const (
	user = "247274526@qq.com"
	pwd  = "vtxgaszpubevbihg"
)

var config = Config{
	ImapHost: "imap.qq.com",
	ImapPort: "993",
	SMTPHost: "smtp.qq.com",
	SMTPPort: "25",
	Domain:   "qq.com",
}

var mm, _ = New(config)

var e = SendEmail{
	Date:    time.Now().Format(timeFormat),
	From:    "247274526@qq.com",
	To:      "247274526@qq.com",
	Subject: "test send email:subject",
	Text:    "test send email:text",
}

func TestModel_SmtpSendEmail(t *testing.T) {
	var err error
	if err = mm.SMTPSendEmail(user, pwd, e); err != nil {
		t.Error(err)
		return
	}
}
