package config

import (
	"encoding/json"
	"os"
	"strings"

	"yumi/pkg/conf"
)

type Document struct {
	Title       string      `json:"title"`    //为查看或编辑的文档定义所需的文件名，该文件名也将在下载文档时用作文件名。
	Url         string      `json:"url"`      //定义存储原始查看或编辑的文档的绝对URL。
	FileType    string      `json:"fileType"` //为查看或编辑的源文件定义文件类型docx
	Key         string      `json:"key"`      //定义用于服务识别文档的唯一文档标识符。万一已知密钥被发送，文档将从缓存中获取。每次编辑和保存文档时，都必须重新生成密钥。可以将文档url用作键，但不能包含特殊字符，并且长度限制为20个符号。
	Info        Info        `json:"info"`     //
	Permissions Permissions `json:"permissions"`
}

type Info struct {
	Author          string          `json:"author"`          //作者
	Created         string          `json:"created"`         //创建时间
	Uploaded        string          `json:"uploaded"`        //上传时间
	Floder          string          `json:"floder"`          //上传文件夹
	Owner           string          `json:"owner"`           //拥有者
	SharingSettings SharingSettings `json:"sharingSettings"` //共享设置
}

type SharingSettings struct {
	Permissions string `json:"permissions"` //权限 Full Access, Read Only, Deny Access
	User        string `json:"user"`        //用户名
}

type Permissions struct {
	Comment  bool `json:"comment"`  //文档是否能被评论
	Download bool `json:"download"` //文档是否能被下载
	Edit     bool `json:"edit"`     //文档是否能被编辑
	Print    bool `json:"print"`    //文档是否能被打印
	Review   bool `json:"review"`   //文档是否能被审阅

	//ChangeHistory        bool `json:"changeHistory"`
	//FillForms            bool `json:"fillForms"`
	//ModifyFilter         bool `json:"modifyFilter"`
	//ModifyContentControl bool `json:"modifyContentControl"`
}

type EditorConfig struct {
	//ActionLink
	Mode        string `json:"mode"`        //edit, view
	Lang        string `json:"lang"`        //语言zh-CN，en-US
	CallbackUrl string `json:"callbackUrl"` //指定文档存储服务的绝对URL
	CreateUrl   string `json:"createUrl"`   //创建文档时的url

	Plugins       Plugins       `json:"plugins"`
	User          User          `json:"user"`
	Embedded      Embedded      `json:"embedded"` //仅在（type = embedded）嵌入模式下使用
	Customization Customization `json:"customization"`
}

type Plugins struct {
}

type User struct {
	Id   string `json:"id"`   //当前用户id
	Name string `json:"name"` // 当前用户名称
}

type Embedded struct {
	FullscreenUrl string `json:"fullscreenUrl"`
	SaveUrl       string `json:"saveUrl"`
	EmbedUrl      string `json:"embedUrl"`
	ShareUrl      string `json:"shareUrl"`
	ToolbarDocked string `json:"toolbarDocked"`
}

type Customization struct {
	Chat              bool `json:"chat"`              //是否显示聊天
	CommentAuthorOnly bool `json:"commentAuthorOnly"` //评论是否只读
	Comments          bool `json:"comments"`          //是否显示评论
	Feedback          bool `json:"feedback"`          //是否显示反馈按钮
	Forcesave         bool `json:"forcesave"`         //当文档保存时

	Customer Customer `json:"customer"`
	Goback   Goback   `json:"goback"`
	Logo     Logo     `json:"logo"`
}

type Feedback struct{}

//编辑者的信息，对后续打开文件的人都可见
type Customer struct {
	Address string `json:"address"` //地址
	Info    string `json:"info"`    //附加信息
	Logo    string `json:"logo"`    //头像
	Mail    string `json:"mail"`    //邮箱
	Name    string `json:"name"`    //名称
	Www     string `json:"www"`     //个人或公司网站
}

type Goback struct {
	Url string `json:"url"` //打开文件位置按钮的url
}

type Logo struct {
	Image         string `json:"image"`         //172x40
	ImageEmbedded string `json:"imageEmbedded"` //248x40
	Url           string `json:"url"`           //点击图片跳转的url
}

type Config struct {
	Width        string       `json:"width"`        //宽度
	Height       string       `json:"height"`       //高度
	Type         string       `json:"type"`         //用于访问文档的平台类型，desktop/mobile/embedded
	DocumentType string       `json:"documentType"` //要打开的文档类型.doc，.docm，.docx，.dot，.dotm，.dotx，.epub，.fodt，.htm，.html，.mht，.odt 、. ott，.pdf，.rtf，.txt，.djvu，.xps
	Token        string       `json:"token"`        //用于文档服务器配置的加密签名。
	Document     Document     `json:"document"`
	EditorConfig EditorConfig `json:"editorConfig"`
}

var _conf Config

func Init(conf conf.OnlyOffice) {
	f, err := os.Open(conf.Document.ConfigPath)
	if err != nil {
		panic(err)
	}
	if err := json.NewDecoder(f).Decode(&_conf); err != nil {
		panic(err)
	}

	return
}

func Get() Config {
	return _conf
}

func (c Config) ToString() string {
	bs, _ := json.Marshal(&c)
	str := string(bs)
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	return str
}
