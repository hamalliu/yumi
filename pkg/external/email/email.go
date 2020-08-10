package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/mail"
	"net/smtp"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/jhillyerd/enmime"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

const (
	mailBoxDrafts  = "Drafts"
	mailBoxSent    = "Sent"
	mailBoxInbox   = "INBOX"
	mailBoxDeleted = "Deleted"
)

const (
	mailHeaderSubject = "Subject"
	mailHeaderFrom    = "From"
	mailHeaderTo      = "To"
)

//UploadFile 上传文件
type UploadFile struct {
	FileName    string
	ContentType string
	Content     []byte
}

//UploadInlines 上传inlines结构
type UploadInlines struct {
	FileName    string
	ContentType string
	Content     []byte
	ContentID   string
}

//ReceiveEmail 接受邮件
type ReceiveEmail struct {
	Seen    bool   `json:"seen"`
	Date    string `json:"date"`
	UID     string `json:"uid"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`

	Text    string            `json:"text"`
	HTML    string            `json:"html"`
	Files   []string          `json:"files"`
	Inlines map[string][]byte `json:"inlines"`
}

//SendEmail 发送邮件
type SendEmail struct {
	Date    string
	From    string
	To      string
	Cc      string
	Bcc     string
	Subject string

	Text string
	HTML string

	UploadFiles   []UploadFile
	UploadInlines []UploadInlines

	from *mail.Address
	to   []*mail.Address
	cc   []*mail.Address
	bcc  []*mail.Address
}

//Config 邮箱服务器配置
type Config struct {
	ImapHost string `json:"imaphost"`
	ImapPort string `json:"imapport"`

	SMTPHost string `json:"smtphost"`
	SMTPPort string `json:"smtpport"`

	Domain string `json:"domain"`
}

//Model 邮箱服务器配置
type Model struct {
	conf Config
}

//New ...
func New(conf Config) (*Model, error) {
	//TODO:ping imap AND smtp
	return &Model{conf: conf}, nil
}

func buildManyMailAddress(ma []*mail.Address) string {
	var sa string
	for i := range ma {
		if sa == "" {
			sa = ma[i].String()
		} else {
			sa = fmt.Sprintf("%s,%s", sa, ma[i].String())
		}
	}
	return sa
}

func (m *Model) imapLogin(user, pwd string) (*client.Client, error) {
	var (
		addr = m.conf.ImapHost + ":" + m.conf.ImapPort
		c    *client.Client
		err  error
	)

	c, err = client.DialTLS(addr, &tls.Config{ServerName: m.conf.ImapHost, InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	if err := c.Login(user+"@"+m.conf.Domain, pwd); err != nil {
		return nil, err
	}

	return c, nil
}

func (m *Model) imapAddEmailFlagByUids(user, pwd string, uids string, flag string) error {
	var (
		seqset = &imap.SeqSet{}
		item   = imap.FormatFlagsOp(imap.AddFlags, true)
		flags  = []interface{}{flag}
		c      *client.Client
		err    error
	)

	c, err = m.imapLogin(user, pwd)
	if err != nil {
		return err
	}
	defer func() {
		if err = c.Logout(); err != nil {
			panic(err)
		}
	}()

	if err = seqset.Add(uids); err != nil {
		return nil
	}
	if err = c.UidStore(seqset, item, flags, nil); err != nil {
		return err
	}

	return nil
}

//SetSeenFlagByUids 通过uids设置已读标志
func (m *Model) SetSeenFlagByUids(user, pwd string, uids string) error {
	var (
		err error
	)

	if err = m.imapAddEmailFlagByUids(user, pwd, uids, imap.SeenFlag); err != nil {
		return err
	}

	return nil
}

//GetEmailListByMailBox 通过邮箱名称获取邮件列表
func (m *Model) GetEmailListByMailBox(user, pwd string, mbName string, offset, line uint) ([]ReceiveEmail, uint, uint, error) {
	var (
		err    error
		emails []ReceiveEmail
		count  uint
		c      *client.Client
		mbox   *imap.MailboxStatus
		env    *enmime.Envelope
		to     []*mail.Address
	)

	if mbName != mailBoxDrafts &&
		mbName != mailBoxSent &&
		mbName != mailBoxInbox {
		return nil, count, offset, fmt.Errorf("未知的邮箱名%s", mbName)
	}

	if c, err = m.imapLogin(user, pwd); err != nil {
		return nil, count, offset, err
	}
	defer func() {
		if err = c.Logout(); err != nil {
			panic(err)
		}
	}()

	if mbox, err = c.Select(mbName, true); err != nil {
		return nil, count, offset, err
	}

	if mbox.Messages == 0 || mbox.Messages <= uint32(offset) {
		return nil, count, offset, nil
	}
	count = uint(mbox.Messages)
	seqset := new(imap.SeqSet)
	offset = uint(mbox.Messages) - offset
	if line > offset {
		line = offset
	}
	seqset.AddRange(uint32(offset-line), uint32(offset))
	items := []imap.FetchItem{
		imap.FetchUid,
		imap.FetchFlags,
		imap.FetchRFC822Header,
		imap.FetchInternalDate,
	}
	done := make(chan error, 1)
	msgs := make(chan *imap.Message, line)
	go func() {
		done <- c.Fetch(seqset, items, msgs)
	}()
	if err = <-done; err != nil {
		return nil, count, offset, err
	}

	for msg := range msgs {
		var email ReceiveEmail
		rfc822Header := imap.FetchItem(imap.FetchRFC822Header)
		for bk, bv := range msg.Body {
			if bk.FetchItem() == rfc822Header {
				if env, err = enmime.ReadEnvelope(bv); err != nil {
					return emails, count, offset, err
				}

				email.Subject = env.GetHeader(mailHeaderSubject)
				email.From = env.GetHeader(mailHeaderFrom)
				if to, err = env.AddressList(mailHeaderTo); err != nil {
					return emails, count, offset, err
				}
				email.To = buildManyMailAddress(to)
				continue
			}
		}
		for _, flag := range msg.Flags {
			if flag == imap.SeenFlag {
				email.Seen = true
			}
		}
		email.Date = msg.InternalDate.Format(timeFormat)
		email.UID = fmt.Sprintf("%d", msg.Uid)
	}

	return emails, count, offset, nil
}

//GetEmailByUID 通过uid获取邮件
func (m *Model) GetEmailByUID(user, pwd string, uid string) (ReceiveEmail, error) {
	var (
		err    error
		email  ReceiveEmail
		c      *client.Client
		seqset = new(imap.SeqSet)
		env    *enmime.Envelope
		to     []*mail.Address
	)

	if c, err = m.imapLogin(user, pwd); err != nil {
		return email, err
	}
	defer func() {
		if err = c.Logout(); err != nil {
			panic(err)
		}
	}()

	if err = seqset.Add(uid); err != nil {
		return email, err
	}
	items := []imap.FetchItem{
		imap.FetchRFC822,
		imap.FetchFlags,
		imap.FetchInternalDate,
	}
	msgs := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.UidFetch(seqset, items, msgs)
	}()
	if err = <-done; err != nil {
		return email, err
	}
	for msg := range msgs {
		rfc822 := imap.FetchItem(imap.FetchRFC822)
		for bk, bv := range msg.Body {
			if bk.FetchItem() == rfc822 {
				if env, err = enmime.ReadEnvelope(bv); err != nil {
					return email, err
				}
				email.Subject = env.GetHeader(mailHeaderSubject)
				email.From = env.GetHeader(mailHeaderFrom)
				if to, err = env.AddressList(mailHeaderTo); err != nil {
					return email, err
				}
				email.To = buildManyMailAddress(to)
				email.Text = env.Text
				email.HTML = env.HTML
				for i := range env.Attachments {
					email.Files = append(email.Files, env.Attachments[i].FileName)
				}
				email.Inlines = make(map[string][]byte)
				for i := range env.Inlines {
					email.Inlines[env.Inlines[i].FileName] = env.Inlines[i].Content
				}
				continue
			}

			email.Date = msg.InternalDate.Format(timeFormat)
			email.UID = uid
			for i := range msg.Flags {
				if imap.SeenFlag == msg.Flags[i] {
					email.Seen = true
				}
			}
		}
	}

	return email, nil
}

//GetEmailAttachmentByFileName ...
func (m *Model) GetEmailAttachmentByFileName(user, pwd string, uid string, fileName string) (io.Reader, error) {
	var (
		err    error
		c      *client.Client
		seqset = new(imap.SeqSet)
		env    *enmime.Envelope
	)

	if c, err = m.imapLogin(user, pwd); err != nil {
		return nil, err
	}
	defer func() {
		if err = c.Logout(); err != nil {
			panic(err)
		}
	}()

	if err = seqset.Add(uid); err != nil {
		return nil, err
	}
	items := []imap.FetchItem{
		imap.FetchRFC822,
	}
	msgs := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.UidFetch(seqset, items, msgs)
	}()
	if err = <-done; err != nil {
		return nil, err
	}
	for msg := range msgs {
		rfc822 := imap.FetchItem(imap.FetchRFC822)
		for bk, bv := range msg.Body {
			if bk.FetchItem() == rfc822 {
				if env, err = enmime.ReadEnvelope(bv); err != nil {
					return nil, err
				}

				for i := range env.Attachments {
					if env.Attachments[i].FileName == fileName {
						return bytes.NewReader(env.Attachments[i].Content), nil
					}
					break
				}
				break
			}
		}
	}

	return nil, nil
}

//DeleteEmailByUID 删除邮件通过uid
func (m *Model) DeleteEmailByUID(user, pwd string, uid string) error {
	var (
		err    error
		c      *client.Client
		seqset = new(imap.SeqSet)
		item   = imap.FormatFlagsOp(imap.AddFlags, true)
		flags  = []interface{}{imap.DeletedFlag}
	)

	if c, err = m.imapLogin(user, pwd); err != nil {
		return err
	}
	defer func() {
		if err = c.Logout(); err != nil {
			panic(err)
		}
	}()

	if err = seqset.Add(uid); err != nil {
		return err
	}
	if err := c.UidStore(seqset, item, flags, nil); err != nil {
		return err
	}

	if err = c.Expunge(nil); err != nil {
		return err
	}

	return nil
}

//MoveInDeletedMailBox 把uid的邮件移入已删除邮箱
func (m *Model) MoveInDeletedMailBox(user, pwd string, uid string) error {
	var (
		err    error
		c      *client.Client
		seqset = new(imap.SeqSet)
	)

	if c, err = m.imapLogin(user, pwd); err != nil {
		return err
	}
	defer func() {
		if err = c.Logout(); err != nil {
			panic(err)
		}
	}()

	if err = seqset.Add(uid); err != nil {
		return err
	}
	items := []imap.FetchItem{
		imap.FetchRFC822,
		imap.FetchFlags,
		imap.FetchInternalDate,
	}
	msgs := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.UidFetch(seqset, items, msgs)
	}()
	if err = <-done; err != nil {
		return err
	}

	go func() {
		done <- m.DeleteEmailByUID(user, pwd, uid)
	}()
	if err = <-done; err != nil {
		return err
	}

	for msg := range msgs {
		rfc822 := imap.FetchItem(imap.FetchRFC822)
		for bk, bv := range msg.Body {
			if bk.FetchItem() == rfc822 {
				if err = c.Append(mailBoxDeleted, msg.Flags, msg.InternalDate, bv); err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

//AppendToMyMailbox 追加到我的邮箱
func (m *Model) AppendToMyMailbox(user, pwd string, mbName string, flags []string, e SendEmail) error {
	var (
		err error
		c   *client.Client
	)

	if mbName != mailBoxDrafts &&
		mbName != mailBoxSent &&
		mbName != mailBoxInbox {
		return fmt.Errorf("未知的邮箱名%s", mbName)
	}

	if c, err = m.imapLogin(user, pwd); err != nil {
		return err
	}
	defer func() {
		if err = c.Logout(); err != nil {
			panic(err)
		}
	}()

	b := &bytes.Buffer{}
	if err = e.buildMessage(b); err != nil {
		return err
	}
	if err = c.Append(mbName, flags, time.Now(), b); err != nil {
		return err
	}

	return nil
}

//SMTPSendEmail Smtp发送邮件
func (m *Model) SMTPSendEmail(user, pwd string, e SendEmail) error {
	var (
		err error

		done = make(chan error, 1)
	)

	go func() {
		done <- e.send(m.conf.SMTPHost+":"+m.conf.SMTPPort, smtp.PlainAuth("", user, pwd, m.conf.SMTPHost))
	}()

	if err = m.AppendToMyMailbox(user, pwd, mailBoxSent, nil, e); err != nil {
		return err
	}

	if err = <-done; err != nil {
		return err
	}

	return nil
}

//BuildMessage 构建消息内容
func (e *SendEmail) buildMessage(outPut io.Writer) error {
	var (
		err     error
		builder enmime.MailBuilder
		p       *enmime.Part
	)

	//set from
	if e.from, err = mail.ParseAddress(e.From); err != nil {
		return err
	}
	builder = builder.From(e.from.Name, e.from.Address)
	//set to
	if e.to, err = mail.ParseAddressList(e.To); err != nil {
		return err
	}
	for i := range e.to {
		builder = builder.To(e.to[i].Name, e.to[i].Address)
	}
	//set cc
	if e.Cc != "" {
		if e.cc, err = mail.ParseAddressList(e.Cc); err != nil {
			return err
		}
		for i := range e.cc {
			builder = builder.To(e.cc[i].Name, e.cc[i].Address)
		}
	}
	//set bcc
	if e.Bcc != "" {
		if e.bcc, err = mail.ParseAddressList(e.Bcc); err != nil {
			return err
		}
		for i := range e.bcc {
			builder = builder.To(e.bcc[i].Name, e.bcc[i].Address)
		}

	}
	//set subject
	if e.Subject != "" {
		builder = builder.Subject(e.Subject)
	}
	//set html
	if e.HTML != "" {
		builder = builder.HTML([]byte(e.HTML))
	}
	//set text
	if e.Text != "" {
		builder = builder.Text([]byte(e.Text))
	}
	//set upfile
	for i := range e.UploadFiles {
		builder = builder.AddAttachment(
			e.UploadFiles[i].Content,
			e.UploadFiles[i].ContentType,
			e.UploadFiles[i].FileName)
		//builder = builder.AddFileAttachment(e.Files[i])
	}
	//set upinlines
	for i := range e.UploadInlines {
		builder = builder.AddInline(
			e.UploadInlines[i].Content,
			e.UploadInlines[i].ContentType,
			e.UploadInlines[i].FileName,
			e.UploadInlines[i].ContentID)
	}
	//set date
	builder = builder.Date(time.Now())

	if p, err = builder.Build(); err != nil {
		return err
	}

	if err = p.Encode(outPut); err != nil {
		return err
	}

	return nil
}

//Send 发送
func (e *SendEmail) send(addr string, a smtp.Auth) error {
	var (
		err error

		buf = &bytes.Buffer{}
	)

	if err = e.buildMessage(buf); err != nil {
		return err
	}

	recips := make([]string, 0, len(e.to)+len(e.cc)+len(e.bcc))
	for _, a := range e.to {
		recips = append(recips, a.Address)
	}
	for _, a := range e.cc {
		recips = append(recips, a.Address)
	}
	for _, a := range e.bcc {
		recips = append(recips, a.Address)
	}
	return smtp.SendMail(addr, a, e.from.Address, recips, buf.Bytes())
}
