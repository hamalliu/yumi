package onlyoffice

import (
	"fmt"
	"net/url"

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

func (oo OnlyOffice) GetHistory(fileName, userId string) ([]History, []HistoryData, error) {
	hs := []History{}
	hds := []HistoryData{}

	countVersion := oo.CountHistoryVersion(fileName, userId)
	uri := oo.GetFileUri(fileName, userId)

	if countVersion == 0 {
		ci, err := oo.GetCreateInfo(fileName, userId)
		if err != nil {
			return nil, nil, err
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

		return hs, hds, nil
	}

	for i := 1; i <= countVersion; i++ {
		rh, err := oo.GetHistoryChanges(fileName, userId, i)
		if err != nil {
			return nil, nil, err
		}
		key, err := oo.GetHistoryKey(fileName, userId, i)
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
		if i > 1 && file_utility.ExistFile(oo.DiffPath(fileName, userId, i-1)) {
			hd.Pervious.Key = hds[i-2].Pervious.Key
			hd.Pervious.Url = hds[i-2].Pervious.Url
			hd.ChangesUrl = oo.GetLocalFileUri(fileName, userId, i-1) + "/diff.zip"
		}

		hs = append(hs, h)
		hds = append(hds, hd)
	}

	return hs, hds, nil
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

func (oo OnlyOffice) GetConvertUri(documentUri, fromExtension, toExtension, documentRevisionId, string, async bool) {

}
