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
	"yumi/internal/onlyoffice/docmanager"
	"yumi/internal/onlyoffice/docservice"
	"yumi/pkg/conf"
	"yumi/pkg/fileutility"
)

//OnlyOffice ...
type OnlyOffice struct {
	cfg conf.OnlyOffice
	docmanager.DocManager
	docservice.DocService
	config.Config
}

//History ...
type History struct {
	docmanager.ResponseHistory
	User    docmanager.User `json:"user"`
	Created string          `json:"created"`

	Key     string `json:"key"`
	Version int    `json:"version"`
}

//HistoryData ...
type HistoryData struct {
	ChangesURL string      `json:"changesUrl"`
	Key        string      `json:"key"`
	Pervious   Pervious    `json:"pervious"`
	URL        string      `json:"url"`
	Version    interface{} `json:"version"`
}

//Pervious ...
type Pervious struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

var _oo OnlyOffice

//Init ...
func Init(cfg conf.OnlyOffice) {
	_oo = OnlyOffice{
		cfg:        cfg,
		DocManager: docmanager.New(cfg.Document),
		DocService: docservice.New(cfg.Token),
	}
}

//Get ...
func Get() OnlyOffice {
	return _oo
}

//GetHistory ...
func (oo OnlyOffice) GetHistory(fileName, userID string) ([]History, []HistoryData, int, error) {
	hs := []History{}
	hds := []HistoryData{}

	countVersion := oo.CountHistoryVersion(fileName, userID) + 1
	uri := oo.GetFileURI(fileName, userID)

	if countVersion == 1 {
		ci, err := oo.GetCreateInfo(fileName, userID)
		if err != nil {
			return nil, nil, 0, err
		}

		key := oo.GenerateKey(fileName, userID)

		h := History{
			Key:     key,
			Version: 1,
			User: docmanager.User{
				ID:   ci.UserID,
				Name: ci.UserName,
			},
			Created: ci.Created,
		}

		hd := HistoryData{
			Version: 1,
			Key:     key,
			URL:     uri,
		}

		hs = append(hs, h)
		hds = append(hds, hd)

		return hs, hds, countVersion, nil
	}

	for i := 1; i <= countVersion; i++ {
		rh, err := oo.GetHistoryChanges(fileName, userID, i)
		if err != nil {
			return nil, nil, 0, err
		}
		key, err := oo.GetHistoryKey(fileName, userID, i)
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
			URL:     uri,
		}
		if i > 1 && fileutility.ExistFile(oo.DiffPath(fileName, userID, i-1)) {
			hd.Pervious.Key = hds[i-2].Pervious.Key
			hd.Pervious.URL = hds[i-2].Pervious.URL
			hd.ChangesURL = oo.GetLocalFileURI(fileName, userID, i-1) + "/diff.zip"
		}

		hs = append(hs, h)
		hds = append(hds, hd)
	}

	return hs, hds, countVersion, nil
}

//GetFileURI ...
func (oo OnlyOffice) GetFileURI(fileName, userID string) string {
	return oo.GetLocalFileURI(fileName, userID, 0)
}

//GetLocalFileURI ...
func (oo OnlyOffice) GetLocalFileURI(fileName, userID string, version int) string {
	fileURI := fmt.Sprintf("%s/%s/%s/%s", oo.cfg.DocumentServerURL, oo.cfg.Document.StoragePath, userID, fileName)
	if version != 0 {
		fileURI = fmt.Sprintf("%s-history/%d", fileURI, version)
	}

	return url.PathEscape(fileURI)
}

//GenerateKey ...
func (oo OnlyOffice) GenerateKey(fileName, userID string) string {
	key := oo.GetLocalFileURI(fileName, userID, 0)
	storagePath := oo.StoragePath(fileName, userID)
	f, _ := os.Stat(storagePath)
	key = key + oo.cfg.DocumentServerURL + key + f.ModTime().Format("2006-01-02 15:04:05")

	return oo.GenerateRevisionID(key)
}

//GenerateRevisionID ...
func (oo OnlyOffice) GenerateRevisionID(expectedKey string) string {
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

//RespConvert ...
type RespConvert struct {
	EndConvert bool   `json:"endConvert"`
	FileURL    string `json:"fileUrl"`
	Percent    int    `json:"percent"`
	Error      int    `json:"error"`
}

//GetConvertURI ...
func (oo OnlyOffice) GetConvertURI(documentURI, fromExtension, toExtension, documentRevisionID string, async bool) (RespConvert, error) {
	res := RespConvert{}

	if fromExtension == "" || toExtension == "" {
		return res, fmt.Errorf("the fromExtention or the toExtention is empty")
	}

	if documentRevisionID == "" {
		documentRevisionID = oo.GenerateRevisionID(documentURI)
	}

	params := struct {
		Async      bool   `json:"async"`
		URL        string `json:"url"`
		OutputType string `json:"outputtype"`
		FileType   string `json:"filetype"`
		Title      string `json:"title"`
		Key        string `json:"key"`
	}{
		Async:      async,
		URL:        documentURI,
		OutputType: strings.ReplaceAll(toExtension, ".", ""),
		FileType:   strings.ReplaceAll(fromExtension, ".", ""),
		Title:      oo.GetFileName(documentURI, false),
		Key:        documentRevisionID,
	}
	body, _ := json.Marshal(&params)

	uri := oo.cfg.SiteURL + oo.cfg.ConverterURL
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return res, err
	}
	req.Header.Set("Accept", "application/json")
	if oo.cfg.Token.Enable && oo.cfg.Token.UseForRequest {
		token, err := oo.FillJwtByURL(uri, params, "", nil)
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

//ReqCommandService ...
type ReqCommandService struct {
	C        string   `json:"c"`
	Key      string   `json:"key"`
	Meta     Meta     `json:"meta"`
	Token    string   `json:"token"`
	UserData string   `json:"userdata"`
	Users    []string `json:"users"`
}

//Meta ...
type Meta struct {
	Title string `json:"title"`
}

//CommandForceSave ...
func (oo OnlyOffice) CommandForceSave(documentRevisionID string) error {
	documentRevisionID = oo.GenerateRevisionID(documentRevisionID)
	params := ReqCommandService{
		C:   "forcesave",
		Key: documentRevisionID,
	}
	body, _ := json.Marshal(&params)

	uri := oo.cfg.SiteURL + oo.cfg.CommandURL
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if oo.cfg.Token.Enable && oo.cfg.Token.UseForRequest {
		token, err := oo.FillJwtByURL(uri, params, "", nil)
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

//File ...
type File struct {
	Name        string
	Ext         string
	URI         string
	Key         string
	Version     int
	Created     string
	Author      string
	Permissions string //权限 Full Access, Read Only, Deny Access
	User        string //用户名
}

//Customer ...
type Customer struct {
	Address string //地址
	Info    string //附加信息
	Logo    string //头像
	Mail    string //邮箱
	Name    string //名称
	Www     string //个人或公司网站
}

//Editor ...
type Editor struct {
	File        File
	GobackURL   string
	Customer    Customer
	Mode        string
	CallbackURL string
	UserID      string
	UserName    string

	Type         string
	DocumentType string
	Token        string
}

//GetConfigStr ...
func (oo OnlyOffice) GetConfigStr(cfg Editor) string {
	c := oo.Config
	c.Type = cfg.Type
	c.Token = cfg.Token
	c.DocumentType = cfg.DocumentType

	c.Document.Title = cfg.File.Name
	c.Document.URL = cfg.File.URI
	c.Document.FileType = cfg.File.Ext
	c.Document.Key = cfg.File.Key
	c.Document.Info.Author = cfg.File.Author
	c.Document.Info.Created = cfg.File.Created
	c.Document.Info.Owner = cfg.File.Author
	c.Document.Info.SharingSettings.User = cfg.File.Permissions
	c.Document.Info.SharingSettings.Permissions = cfg.File.User

	c.EditorConfig.Mode = cfg.Mode
	c.EditorConfig.CallbackURL = cfg.CallbackURL
	c.EditorConfig.User.ID = cfg.UserID
	c.EditorConfig.User.Name = cfg.UserName
	c.EditorConfig.Customization.Goback.URL = cfg.GobackURL
	c.EditorConfig.Customization.Customer.Name = cfg.Customer.Name
	c.EditorConfig.Customization.Customer.Www = cfg.Customer.Www
	c.EditorConfig.Customization.Customer.Mail = cfg.Customer.Mail
	c.EditorConfig.Customization.Customer.Logo = cfg.Customer.Logo
	c.EditorConfig.Customization.Customer.Address = cfg.Customer.Address
	c.EditorConfig.Customization.Customer.Info = cfg.Customer.Info
	c.EditorConfig.Embedded.EmbedURL = cfg.File.URI
	c.EditorConfig.Embedded.SaveURL = cfg.File.URI
	c.EditorConfig.Embedded.ShareURL = cfg.File.URI

	return c.ToString()
}
