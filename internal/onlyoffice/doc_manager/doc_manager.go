package doc_manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"yumi/pkg/conf"
	"yumi/pkg/file_utility"
)

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

type CreateInfo struct {
	Created  string
	UserId   string
	UserName string
}

//======================================================================================================================
type DocManager struct {
	FileUtility
}

const timeFormat = "2006-01-02 15:04:05"

func New(cfg conf.Document) DocManager {
	return DocManager{FileUtility: FileUtility{cfg: cfg}}
}

func (dm DocManager) CreateDemo(fileName, sampleName, fileExt, userId, userName string) (string, error) {
	fileName = dm.GetCorrectName(fileName+fileExt, userId)
	filePath := dm.StoragePath(fileName, userId)

	sampleFile := ""
	if sampleName == "" {
		sampleName = "new"
	}
	sampleFile = path.Join(dm.cfg.SamplesPath, sampleName+fileExt)

	if err := file_utility.CopyFile(sampleFile, filePath); err != nil {
		return "", err
	}

	if err := dm.SaveCreateInfo(fileName, userId, userName); err != nil {
		return "", err
	}

	return fileName, nil
}

func (dm DocManager) SaveCreateInfo(fileName, userId, userName string) error {
	filePath := dm.StoragePath(fileName, userId)

	s, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	modTimeStr := s.ModTime().Format(timeFormat)

	directory := dm.HistoryPath(fileName, userId, true)
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

	historyPath := dm.HistoryPath(fileName, userId, false)

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

func (dm DocManager) GenerateKey(fileName, userId string) string {
	//TODO
	return ""
}

//return: history目录下changes.txt内容
func (dm DocManager) GetHistoryChanges(fileName, userId string, version int) (ResponseHistory, error) {
	rh := ResponseHistory{}

	changesPath := dm.ChangesPath(fileName, userId, version)
	f, err := os.Open(changesPath)
	if err != nil {
		return rh, err
	}

	if err := json.NewDecoder(f).Decode(&rh); err != nil {
		return rh, err
	}

	return rh, nil
}

//return: history 目录下key.txt内容
func (dm DocManager) GetHistoryKey(fileName, userId string, version int) (string, error) {
	keyPath := dm.KeyPath(fileName, userId, version)
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
	directory := dm.HistoryPath(fileName, userId, false)
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
