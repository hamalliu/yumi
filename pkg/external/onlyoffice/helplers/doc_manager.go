package helplers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"yumi/pkg/conf"
	"yumi/pkg/file_utility"
)

type History struct {
	ResponseHistory
	User    User   `json:"user"`
	Created string `json:"created"`

	Key     string `json:"key"`
	Version int    `json:"version"`
}

type ResponseHistory struct {
	ServerVersion string   `json:"serverVersion"`
	Changes       []Change `json:"changes"`
}

type Change struct {
	Created string `json:"created"`
	User    User   `json:"user"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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

type CreateInfo struct {
	Created  string
	UserId   string
	UserName string
}

//======================================================================================================================
type DocManager struct {
	req *http.Request
}

const timeFormat = "2006-01-02 15:04:05"

func New(req *http.Request) DocManager {
	return DocManager{req: req}
}

func (dm DocManager) CreateDemo(fileName, sampleName, fileExt, userId, userName string) (string, error) {
	fileName = fileUtility.GetCorrectName(fileName+fileExt, userId)
	filePath := fileUtility.StoragePath(fileName, userId)

	cfgOffice := conf.Get().Office
	sampleFile := ""
	if sampleName == "" {
		sampleName = "new"
	}
	sampleFile = path.Join(cfgOffice.SamplesPath, sampleName+fileExt)

	if err := file_utility.CopyFile(sampleFile, filePath); err != nil {
		return "", err
	}

	if err := dm.SaveCreateInfo(fileName, userId, userName); err != nil {
		return "", err
	}

	return fileName, nil
}

func (dm DocManager) SaveCreateInfo(fileName, userId, userName string) error {
	filePath := fileUtility.StoragePath(fileName, userId)

	s, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	modTimeStr := s.ModTime().Format(timeFormat)

	directory := fileUtility.HistoryPath(fileName, userId, true)
	if err := file_utility.CreateDir(directory); err != nil {
		return err
	}

	f, err := os.OpenFile(path.Join(directory, fileName+".txt"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s,%s,%s", modTimeStr, userId, userName))
	if err != nil {
		return err
	}

	return nil
}

func (dm DocManager) GetCreateInfo(fileName, userId string) (CreateInfo, error) {
	ci := CreateInfo{}

	historyPath := fileUtility.HistoryPath(fileName, userId, false)

	f, err := os.Open(path.Join(historyPath, fileName+".txt"))
	if err != nil {
		return ci, err
	}

	cb, err := ioutil.ReadAll(f)
	if err != nil {
		return ci, err
	}

	feilds := strings.Split(string(cb), ",")
	ci.Created = feilds[0]
	ci.UserId = feilds[1]
	ci.UserName = feilds[2]

	return ci, nil
}

func (dm DocManager) GetFileUri(fileName, userId string) string {
	return dm.GetLocalFileUri(fileName, userId, 0)
}

func (dm DocManager) GetLocalFileUri(fileName, userId string, version int) string {
	serverUrl := dm.GetServerUrl()
	storagePath := conf.Get().Office.StoragePath
	fileUri := fmt.Sprintf("%s/%s/%s/%s", serverUrl, storagePath, userId, fileName)
	if version != 0 {
		fileUri = fmt.Sprintf("%s-history/%d", fileUri, version)
	}

	return url.PathEscape(fileUri)
}

func (dm DocManager) GetServerUrl() string {
	exampleUrl := conf.Get().Office.DocumentServerUrl
	if exampleUrl == "" {
		return fmt.Sprintf("%s://%s", dm.req.URL.Scheme, dm.req.URL.Host)
	} else {
		return exampleUrl
	}
}

func (dm DocManager) GetCallback(fileName, userId string) string {
	return fmt.Sprintf("%s/track?filename=%s&userid=%s", dm.GetServerUrl(), fileName, userId)
}

func (dm DocManager) GenerateKey(fileName, userId string) string {
	//TODO
	return ""
}

func (dm DocManager) GetHistoryChanges(fileName, userId string, version int) (ResponseHistory, error) {
	rh := ResponseHistory{}

	changesPath := fileUtility.ChangesPath(fileName, userId, version)
	f, err := os.Open(changesPath)
	if err != nil {
		return rh, err
	}

	if err := json.NewDecoder(f).Decode(&rh); err != nil {
		return rh, err
	}

	return rh, nil
}

func (dm DocManager) GetHistoryKey(fileName, userId string, version int) (string, error) {
	keyPath := fileUtility.KeyPath(fileName, userId, version)
	f, err := os.Open(keyPath)
	if err != nil {
		return "", err
	}

	key, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

func (dm DocManager) CountHistoryVersion(fileName, userId string) int {
	directory := fileUtility.HistoryPath(fileName, userId, false)
	if directory == "" {
		return 0
	}

	index := 0
	for {
		index++
		if !file_utility.ExistDir(path.Join(directory, fmt.Sprintf("%d", index))) {
			break
		}
	}

	return index - 1
}

func (dm DocManager) GetHistory(fileName, userId string) ([]History, []HistoryData, error) {
	hs := []History{}
	hds := []HistoryData{}

	countVersion := dm.CountHistoryVersion(fileName, userId)
	uri := dm.GetFileUri(fileName, userId)

	if countVersion == 0 {
		ci, err := dm.GetCreateInfo(fileName, userId)
		if err != nil {
			return nil, nil, err
		}

		key := dm.GenerateKey(fileName, userId)

		h := History{
			Key:     key,
			Version: 1,
			User: User{
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

		return hs, hds, nil
	}

	for i := 1; i <= countVersion; i++ {
		rh, err := dm.GetHistoryChanges(fileName, userId, i)
		if err != nil {
			return nil, nil, err
		}
		key, err := dm.GetHistoryKey(fileName, userId, i)
		if err != nil {
			return nil, nil, err
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

		hs = append(hs, h)
		hds = append(hds, hd)
	}

	return hs, hds, nil
}
