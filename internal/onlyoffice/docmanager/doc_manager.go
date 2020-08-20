package docmanager

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

//ResponseHistory ...
type ResponseHistory struct {
	ServerVersion string   `json:"serverVersion"`
	Changes       []Change `json:"changes"`
}

//Change ...
type Change struct {
	Created string `json:"created"`
	User    User   `json:"user"`
}

//User ...
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//CreateInfo ...
type CreateInfo struct {
	Created  string
	UserID   string
	UserName string
}

//======================================================================================================================

//DocManager ...
type DocManager struct {
	FileUtility
}

const timeFormat = "2006-01-02 15:04:05"

//New ...
func New(cfg conf.Document) DocManager {
	return DocManager{FileUtility: FileUtility{cfg: cfg}}
}

//SaveCreateInfo ...
func (dm DocManager) SaveCreateInfo(fileName, userID, userName string) error {
	filePath := dm.StoragePath(fileName, userID)

	s, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	modTimeStr := s.ModTime().Format(timeFormat)

	directory := dm.HistoryPath(fileName, userID, true)
	if err := file_utility.CreateDir(directory); err != nil {
		return err
	}

	f, err := os.OpenFile(path.Join(directory, fileName+".txt"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s,%s,%s", modTimeStr, userID, userName))
	if err != nil {
		return err
	}

	return nil
}

//GetCreateInfo ...
func (dm DocManager) GetCreateInfo(fileName, userID string) (CreateInfo, error) {
	ci := CreateInfo{}

	historyPath := dm.HistoryPath(fileName, userID, false)

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
	ci.UserID = feilds[1]
	ci.UserName = feilds[2]

	return ci, nil
}

//GetHistoryChanges return: history目录下changes.txt内容
func (dm DocManager) GetHistoryChanges(fileName, userID string, version int) (ResponseHistory, error) {
	rh := ResponseHistory{}

	changesPath := dm.ChangesPath(fileName, userID, version)
	f, err := os.Open(changesPath)
	if err != nil {
		return rh, err
	}

	if err := json.NewDecoder(f).Decode(&rh); err != nil {
		return rh, err
	}

	return rh, nil
}

//GetHistoryKey return: history 目录下key.txt内容
func (dm DocManager) GetHistoryKey(fileName, userID string, version int) (string, error) {
	keyPath := dm.KeyPath(fileName, userID, version)
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

//CountHistoryVersion ...
func (dm DocManager) CountHistoryVersion(fileName, userID string) int {
	directory := dm.HistoryPath(fileName, userID, false)
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
