package onlyoffice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"yumi/internal/onlyoffice/config"
	"yumi/internal/onlyoffice/doc_manager"
	"yumi/internal/onlyoffice/doc_service"
	"yumi/pkg/conf"
	"yumi/pkg/file_utility"
)

type OnlyOffice struct {
	cfg conf.OnlyOffice
	doc_manager.DocManager
	doc_service.DocService
	config.Config
}

type History struct {
	doc_manager.ResponseHistory
	User    doc_manager.User `json:"user"`
	Created string           `json:"created"`

	Key     string `json:"key"`
	Version int    `json:"version"`
}

type HistoryData struct {
	ChangesUrl string      `json:"changesUrl"`
	Key        string      `json:"key"`
	Pervious   Pervious    `json:"pervious"`
	Url        string      `json:"url"`
	Version    interface{} `json:"version"`
}

type Pervious struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

var _oo OnlyOffice

func Init(cfg conf.OnlyOffice) {
	_oo = OnlyOffice{
		cfg:        cfg,
		DocManager: doc_manager.New(cfg.Document),
		DocService: doc_service.New(cfg.Token),
	}
}

func Get() OnlyOffice {
	return _oo
}

func (oo OnlyOffice) GetHistory(fileName, userId string) ([]History, []HistoryData, int, error) {
	hs := []History{}
	hds := []HistoryData{}

	countVersion := oo.CountHistoryVersion(fileName, userId) + 1
	uri := oo.GetFileUri(fileName, userId)

	if countVersion == 1 {
		ci, err := oo.GetCreateInfo(fileName, userId)
		if err != nil {
			return nil, nil, 0, err
		}

		key := oo.GenerateKey(fileName, userId)

		h := History{
			Key:     key,
			Version: 1,
			User: doc_manager.User{
				Id:   ci.UserId,
				Name: ci.UserName,
			},
			Created: ci.Created,
		}

		hd := HistoryData{
			Version: 1,
			Key:     key,
			Url:     uri,
		}

		hs = append(hs, h)
		hds = append(hds, hd)

		return hs, hds, countVersion, nil
	}

	for i := 1; i <= countVersion; i++ {
		rh, err := oo.GetHistoryChanges(fileName, userId, i)
		if err != nil {
			return nil, nil, 0, err
		}
		key, err := oo.GetHistoryKey(fileName, userId, i)
		if err != nil {
			return nil, nil, 0, err
		}

		h := History{
			ResponseHistory: rh,
			Key:             key,
			Version:         1,
			User:            rh.Changes[0].User,
			Created:         rh.Changes[0].Created,
		}

		hd := HistoryData{
			Version: i,
			Key:     key,
			Url:     uri,
		}
		if i > 1 && file_utility.ExistFile(oo.DiffPath(fileName, userId, i-1)) {
			hd.Pervious.Key = hds[i-2].Pervious.Key
			hd.Pervious.Url = hds[i-2].Pervious.Url
			hd.ChangesUrl = oo.GetLocalFileUri(fileName, userId, i-1) + "/diff.zip"
		}

		hs = append(hs, h)
		hds = append(hds, hd)
	}

	return hs, hds, countVersion, nil
}

func (oo OnlyOffice) GetFileUri(fileName, userId string) string {
	return oo.GetLocalFileUri(fileName, userId, 0)
}

func (oo OnlyOffice) GetLocalFileUri(fileName, userId string, version int) string {
	fileUri := fmt.Sprintf("%s/%s/%s/%s", oo.cfg.DocumentServerUrl, oo.cfg.Document.StoragePath, userId, fileName)
	if version != 0 {
		fileUri = fmt.Sprintf("%s-history/%d", fileUri, version)
	}

	return url.PathEscape(fileUri)
}

func (oo OnlyOffice) GenerateKey(fileName, userId string) string {
	key := oo.GetLocalFileUri(fileName, userId, 0)
	storagePath := oo.StoragePath(fileName, userId)
	f, _ := os.Stat(storagePath)
	key = key + oo.cfg.DocumentServerUrl + key + f.ModTime().Format("2006-01-02 15:04:05")

	return oo.GenerateRevisionId(key)
}

func (oo OnlyOffice) GenerateRevisionId(expectedKey string) string {
	if len(expectedKey) > 128 {
		return fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(expectedKey)))
	}

	re, err := regexp.Compile("[^0-9-.a-zA-Z_=]")
	if err != nil {
		panic(err)
	}

	matchs := re.FindAllString(expectedKey, -1)
	for i := range matchs {
		expectedKey = strings.ReplaceAll(expectedKey, matchs[i], "_")
	}

	return expectedKey
}

type RespConvert struct {
	EndConvert bool   `json:"endConvert"`
	FileUrl    string `json:"fileUrl"`
	Percent    int    `json:"percent"`
	Error      int    `json:"error"`
}

func (oo OnlyOffice) GetConvertUri(documentUri, fromExtension, toExtension, documentRevisionId string, async bool) (RespConvert, error) {
	res := RespConvert{}

	if fromExtension == "" || toExtension == "" {
		return res, fmt.Errorf("the fromExtention or the toExtention is empty")
	}

	if documentRevisionId == "" {
		documentRevisionId = oo.GenerateRevisionId(documentUri)
	}

	params := struct {
		Async      bool   `json:"async"`
		Url        string `json:"url"`
		OutputType string `json:"outputtype"`
		FileType   string `json:"filetype"`
		Title      string `json:"title"`
		Key        string `json:"key"`
	}{
		Async:      async,
		Url:        documentUri,
		OutputType: strings.ReplaceAll(toExtension, ".", ""),
		FileType:   strings.ReplaceAll(fromExtension, ".", ""),
		Title:      oo.GetFileName(documentUri, false),
		Key:        documentRevisionId,
	}
	body, _ := json.Marshal(&params)

	uri := oo.cfg.SiteUrl + oo.cfg.ConverterUrl
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return res, err
	}
	req.Header.Set("Accept", "application/json")
	if oo.cfg.Token.Enable && oo.cfg.Token.UseForRequest {
		token, err := oo.FillJwtByUrl(uri, params, "", nil)
		if err != nil {
			return res, err
		}
		req.Header.Set(oo.cfg.Token.AuthorizationHeader, oo.cfg.Token.AuthorizationHeaderPrefix+string(token))
	}

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return res, err
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

type ReqCommandService struct {
	C        string   `json:"c"`
	Key      string   `json:"key"`
	Meta     Meta     `json:"meta"`
	Token    string   `json:"token"`
	UserData string   `json:"userdata"`
	Users    []string `json:"users"`
}
type Meta struct {
	Title string `json:"title"`
}

func (oo OnlyOffice) CommandForceSave(documentRevisionId string) error {
	documentRevisionId = oo.GenerateRevisionId(documentRevisionId)
	params := ReqCommandService{
		C:   "forcesave",
		Key: documentRevisionId,
	}
	body, _ := json.Marshal(&params)

	uri := oo.cfg.SiteUrl + oo.cfg.CommandUrl
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if oo.cfg.Token.Enable && oo.cfg.Token.UseForRequest {
		token, err := oo.FillJwtByUrl(uri, params, "", nil)
		if err != nil {
			return err
		}
		req.Header.Set(oo.cfg.Token.AuthorizationHeader, oo.cfg.Token.AuthorizationHeaderPrefix+string(token))
	}

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	res := struct {
		Error int    `json:"error"`
		Key   string `json:"key"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}

	return fmt.Errorf("%s", oo.CommandServiceErrorMessage(res.Error))
}

type File struct {
	Name        string
	Ext         string
	Uri         string
	Key         string
	Version     int
	Created     string
	Author      string
	Permissions string //权限 Full Access, Read Only, Deny Access
	User        string //用户名
}
type Customer struct {
	Address string //地址
	Info    string //附加信息
	Logo    string //头像
	Mail    string //邮箱
	Name    string //名称
	Www     string //个人或公司网站
}

type Editor struct {
	File        File
	GobackUrl   string
	Customer    Customer
	Mode        string
	CallbackUrl string
	UserId      string
	UserName    string

	Type         string
	DocumentType string
	Token        string
}

func (oo OnlyOffice) GetConfigStr(cfg Editor) string {
	c := oo.Config
	c.Type = cfg.Type
	c.Token = cfg.Token
	c.DocumentType = cfg.DocumentType

	c.Document.Title = cfg.File.Name
	c.Document.Url = cfg.File.Uri
	c.Document.FileType = cfg.File.Ext
	c.Document.Key = cfg.File.Key
	c.Document.Info.Author = cfg.File.Author
	c.Document.Info.Created = cfg.File.Created
	c.Document.Info.Owner = cfg.File.Author
	c.Document.Info.SharingSettings.User = cfg.File.Permissions
	c.Document.Info.SharingSettings.Permissions = cfg.File.User

	c.EditorConfig.Mode = cfg.Mode
	c.EditorConfig.CallbackUrl = cfg.CallbackUrl
	c.EditorConfig.User.Id = cfg.UserId
	c.EditorConfig.User.Name = cfg.UserName
	c.EditorConfig.Customization.Goback.Url = cfg.GobackUrl
	c.EditorConfig.Customization.Customer.Name = cfg.Customer.Name
	c.EditorConfig.Customization.Customer.Www = cfg.Customer.Www
	c.EditorConfig.Customization.Customer.Mail = cfg.Customer.Mail
	c.EditorConfig.Customization.Customer.Logo = cfg.Customer.Logo
	c.EditorConfig.Customization.Customer.Address = cfg.Customer.Address
	c.EditorConfig.Customization.Customer.Info = cfg.Customer.Info
	c.EditorConfig.Embedded.EmbedUrl = cfg.File.Uri
	c.EditorConfig.Embedded.SaveUrl = cfg.File.Uri
	c.EditorConfig.Embedded.ShareUrl = cfg.File.Uri

	return c.ToString()
}
