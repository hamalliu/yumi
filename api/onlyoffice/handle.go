package onlyoffice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"

	"yumi/conf"
	"yumi/gin"
	"yumi/gin/valuer"
	"yumi/pkg/binding"
	"yumi/pkg/fileutility"
	"yumi/pkg/log"
	"yumi/pkg/status"
	"yumi/usecase/onlyoffice"
)

//ReqSample ...
type ReqSample struct {
	SampleName    string `json:"sample_name" binding:"required"`
	FileName      string `json:"file_name" binding:"required"`
	FileExtension string `json:"file_extension" binding:"required"`
}

//ReqDownload ...
type ReqDownload struct {
	FileName string `query:"file_name" binding:"required"`
}

//ReqDeleteFile ...
type ReqDeleteFile struct {
	FileName string `query:"file_name"`
}

//ReqConvert ...
type ReqConvert struct {
	FileName string `query:"file_name" binding:"required"`
}

//ReqTrack ...
type ReqTrack struct {
	FileName string `query:"file_name" binding:"required"`
	UserID   string `query:"user_id" binding:"required"`
}

//ReqEditor ...
type ReqEditor struct {
	Mode     string `query:"mode" binding:"oneof=view edit"`
	Type     string `query:"type"`
	FileName string `query:"file_name" binding:"required"`
}

//Upload 上传文件
func Upload(c *gin.Context) {
	req := c.Request
	var (
		mulf  multipart.File
		mulfh *multipart.FileHeader
		osf   *os.File

		err error
	)

	if mulf, mulfh, err = req.FormFile("file"); err != nil {
		c.WriteJSON(nil, err)
		return
	}

	if mulfh.Size > conf.Get().OnlyOffice.Document.MaxFileSize.Size() {
		c.WriteJSON(nil, status.OutOfRange())
		return
	}

	user := c.Get(valuer.KeyUser).User()
	fileName := onlyoffice.Get().GetCorrectName(mulfh.Filename, user.UserID)
	ext := onlyoffice.Get().GetFileExtension(fileName, false)

	if onlyoffice.Get().AllowUploadExtension(ext) {
		c.WriteJSON(nil, status.FailedPrecondition())
		return
	}

	storagePath := onlyoffice.Get().StoragePath(fileName, user.UserID)
	if osf, err = os.OpenFile(storagePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0744); err != nil {
		c.WriteJSON(nil, status.Internal().WithError(err))
		return
	}
	if _, err := io.Copy(osf, mulf); err != nil {
		c.WriteJSON(nil, status.Internal().WithError(err))
		return
	}
	_ = mulf.Close()
	_ = osf.Close()

	if err := onlyoffice.Get().SaveCreateInfo(fileName, user.UserID, user.UserName); err != nil {
		c.WriteJSON(nil, status.Internal().WithError(err))
		return
	}

	c.WriteJSON(nil, nil)
}

//Sample 从样品中复制文件，用于新建文件
func Sample(c *gin.Context) {
	reqm := ReqSample{}
	if err := c.Bind(&reqm); err != nil {
		c.WriteJSON(nil, status.InvalidArgument().WithError(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()

	reqm.FileName = onlyoffice.Get().GetCorrectName(reqm.FileName+reqm.FileExtension, user.UserID)
	filePath := onlyoffice.Get().StoragePath(reqm.FileName, user.UserID)
	sampleFile := ""
	if reqm.SampleName == "" {
		reqm.SampleName = "new"
	}
	sampleFile = path.Join(conf.Get().OnlyOffice.Document.SamplesPath, reqm.SampleName+reqm.FileExtension)
	if err := fileutility.CopyFile(sampleFile, filePath); err != nil {
		c.WriteJSON(nil, status.Internal().WithError(err))
		return
	}

	if err := onlyoffice.Get().SaveCreateInfo(reqm.FileName, user.UserID, user.UserName); err != nil {
		c.WriteJSON(nil, err)
		return
	}

	c.WriteJSON(nil, nil)
}

type renderConfig struct {
	APIURL      string
	Config      string
	Version     int
	History     []onlyoffice.History
	HistoryData []onlyoffice.HistoryData
}

func getCallbackURL(c *gin.Context, fileName, userID string) string {
	return c.Request.URL.Scheme + "://" + c.Request.URL.Host + "/api/track?file_name=" + fileName + "&user_id=" + userID
}

func renderError(c *gin.Context, err error) {
	log.Error(err)
	t, _ := template.ParseFiles("./views/onlyoffice/editor.tpl")
	_ = t.Execute(c.Writer, struct{ message string }{message: "服务器内部错误"})
}

//Editor 按配置加载数据，返回office编辑模板
func Editor(c *gin.Context) {
	reqm := ReqEditor{}
	if err := c.Bind(&reqm); err != nil {
		c.WriteJSON(nil, status.InvalidArgument().WithError(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()
	confOo := conf.Get().OnlyOffice
	history, historyData, version, err := onlyoffice.Get().GetHistory(reqm.FileName, user.UserID)
	if err != nil {
		renderError(c, err)
	}

	editor := onlyoffice.Editor{
		File: onlyoffice.File{
			Name:    reqm.FileName,
			Ext:     onlyoffice.Get().GetFileExtension(reqm.FileName, true),
			URI:     onlyoffice.Get().GetFileURI(reqm.FileName, user.UserID),
			Key:     onlyoffice.Get().GenerateKey(reqm.FileName, user.UserID),
			Version: version,
			Created: time.Now().Format("2006-01-02 15:04:05"),
		},
		Customer: onlyoffice.Customer{
			Name:    user.UserName,
			Info:    "",
			Logo:    "",
			Mail:    "",
			Address: "",
			Www:     "",
		},
		GobackURL:    "",
		Mode:         reqm.Mode,
		CallbackURL:  getCallbackURL(c, reqm.FileName, user.UserID),
		UserID:       user.UserID,
		UserName:     user.UserName,
		Type:         reqm.Type,
		DocumentType: onlyoffice.Get().GetFileType(reqm.FileName),
	}

	if confOo.Token.Enable {
		vals := make(map[string]interface{})
		cs := onlyoffice.Get().GetConfigStr(editor)
		_ = json.Unmarshal([]byte("{"+cs+"}"), &vals)

		token, err := onlyoffice.Get().GetToken(vals)
		if err != nil {
			renderError(c, err)
		}
		editor.Token = string(token)
	}

	rc := renderConfig{
		APIURL:      confOo.SiteURL + confOo.APIURL,
		History:     history,
		HistoryData: historyData,
		Config:      onlyoffice.Get().GetConfigStr(editor),
	}

	t, err := template.ParseFiles("./views/onlyoffice/editor.tpl")
	if err != nil {
		renderError(c, err)
	}

	if err := t.Execute(c.Writer, rc); err != nil {
		renderError(c, err)
	}
}

func save(c *gin.Context, body map[string]interface{}, downloadURI, fileName, userID string) {
	curExt := onlyoffice.Get().GetFileExtension(fileName, false)
	downloadExt := onlyoffice.Get().GetFileExtension(downloadURI, false)

	if curExt != downloadExt {
		key := onlyoffice.Get().GenerateRevisionID(downloadURI)

		resp, err := onlyoffice.Get().GetConvertURI(downloadURI, downloadExt, curExt, key, true)
		if err != nil {
			log.Error(err)
			goto Error
		}
		if resp.Error != 0 {
			err := onlyoffice.Get().ConvertURIErrorMessage(resp.Error)
			log.Error(err)
			goto Error
		}
		save(c, body, resp.FileURL, fileName, userID)
	} else {
		storagePath := onlyoffice.Get().StoragePath(fileName, userID)

		if fileutility.ExistFile(storagePath) {
			historyPath := onlyoffice.Get().HistoryPath(fileName, userID, false)
			if historyPath == "" {
				historyPath = onlyoffice.Get().HistoryPath(fileName, userID, true)
				err := fileutility.CreateDir(historyPath)
				if err != nil {
					log.Error(err)
					goto Error
				}
			}

			countVersion := onlyoffice.Get().CountHistoryVersion(fileName, userID)
			version := countVersion + 1
			versionPath := onlyoffice.Get().VersionPath(fileName, userID, version)
			err := fileutility.CreateDir(versionPath)
			if err != nil {
				log.Error(err)
				goto Error
			}

			//diff
			downloadZip := body["changesurl"].(string)
			if downloadZip != "" {
				diffPath := onlyoffice.Get().DiffPath(fileName, userID, version)
				resp, err := http.Get(downloadZip)
				if err != nil {
					log.Error(err)
					goto Error
				}
				err = fileutility.WriteFile(resp.Body, diffPath)
				if err != nil {
					log.Error(err)
					goto Error
				}
				_ = resp.Body.Close()
			}

			//changes
			changesStr := ""
			if body["changeshistory"] != nil {
				changesStr = fmt.Sprintf("%v", body["changeshistory"])
			} else if body["history"] != nil {
				changesStr = fmt.Sprintf("%v", body["history"])
			}
			changesPath := onlyoffice.Get().ChangesPath(fileName, userID, version)
			err = fileutility.WriteFile(bytes.NewBuffer([]byte(changesStr)), changesPath)
			if err != nil {
				log.Error(err)
				goto Error
			}

			//key
			key := body["key"].(string)
			pathKey := onlyoffice.Get().KeyPath(fileName, userID, version)
			err = fileutility.WriteFile(bytes.NewBuffer([]byte(key)), pathKey)
			if err != nil {
				log.Error(err)
				goto Error
			}

			//prev
			prevPath := onlyoffice.Get().PrevFilePath(fileName, userID, version)
			err = fileutility.CopyFile(storagePath, prevPath)
			if err != nil {
				log.Error(err)
				goto Error
			}

			//storagePath
			resp, err := http.Get(downloadZip)
			if err != nil {
				log.Error(err)
				goto Error
			}
			err = fileutility.WriteFile(resp.Body, storagePath)
			if err != nil {
				log.Error(err)
				goto Error
			}

			//delete forcesavepath
			forceSavePath := onlyoffice.Get().ForcesavePath(fileName, userID, false)
			if forceSavePath != "" {
				err := os.RemoveAll(forceSavePath)
				if err != nil {
					log.Error(err)
					goto Error
				}
			}
		}
	}

	_, _ = c.Writer.Write([]byte("{\"error\":0}"))
	return

Error:
	_, _ = c.Writer.Write([]byte("{\"error\":1}"))
}

func forceSave(c *gin.Context, body map[string]interface{}, downloadURI, fileName, userID string) {
	curExt := onlyoffice.Get().GetFileExtension(fileName, false)
	downloadExt := onlyoffice.Get().GetFileExtension(downloadURI, false)

	if curExt != downloadExt {
		key := onlyoffice.Get().GenerateRevisionID(downloadURI)

		resp, err := onlyoffice.Get().GetConvertURI(downloadURI, downloadExt, curExt, key, true)
		if err != nil {
			log.Error(err)
			goto Error
		}
		if resp.Error != 0 {
			err := onlyoffice.Get().ConvertURIErrorMessage(resp.Error)
			log.Error(err)
			goto Error
		}
		forceSave(c, body, resp.FileURL, fileName, userID)
	} else {
		forceSavePath := onlyoffice.Get().ForcesavePath(fileName, userID, false)
		if forceSavePath == "" {
			forceSavePath = onlyoffice.Get().ForcesavePath(fileName, userID, true)
			err := fileutility.CreateDir(forceSavePath)
			if err != nil {
				log.Error(err)
				goto Error
			}
		}

		resp, err := http.Get(downloadURI)
		if err != nil {
			log.Error(err)
			goto Error
		}

		err = fileutility.WriteFile(resp.Body, forceSavePath)
		if err != nil {
			log.Error(err)
			goto Error
		}
	}
	_, _ = c.Writer.Write([]byte("{\"error\":0}"))
	return

Error:
	_, _ = c.Writer.Write([]byte("{\"error\":1}"))
}

func track(c *gin.Context, body map[string]interface{}, fileName, userID string) {
	status := body["status"].(float64)
	if status == 1 {
		if body["actions"] != nil {
			action := body["actions"].([]interface{})[0].(map[string]interface{})
			if action["type"].(float64) == 0 {
				user := action["userid"]

				exist := false
				users := body["users"].([]interface{})
				for i := range users {
					if users[i].(string) == user.(string) {
						exist = true
					}
				}

				if !exist {
					if err := onlyoffice.Get().CommandForceSave(body["key"].(string)); err != nil {
						log.Error(err)
						_, _ = c.Writer.Write([]byte("{\"error\":1}"))
						return
					}
				}
			}
		}
	} else if status == 2 || status == 3 {
		save(c, body, body["url"].(string), fileName, userID)
	} else if status == 6 || status == 7 {
		forceSave(c, body, body["url"].(string), fileName, userID)
	}

	_, _ = c.Writer.Write([]byte("{\"error\":0}"))
}

//Track onlyoffice的回调处理
func Track(c *gin.Context) {
	reqm := ReqTrack{}
	bodyMap := make(map[string]interface{})

	if err := c.BindWith(&reqm, binding.Query); err != nil {
		log.Error(err, c.Request.URL)
		goto Error
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&bodyMap); err != nil {
		log.Error(err, c.Request.Body)
		goto Error
	}

	if conf.Get().OnlyOffice.Token.Enable && conf.Get().OnlyOffice.Token.UseForRequest {
		var body map[string]interface{}
		var err error
		if token, ok := bodyMap["token"].(string); ok && token != "" {
			body, err = onlyoffice.Get().ReadToken(token)
			if err != nil {
				log.Error(err)
				goto Error
			}
		} else {
			body, err = onlyoffice.Get().CheckJwtHeader(c.Request)
			if err != nil || body == nil {
				log.Error(err)
				goto Error
			}
			if body["preload"] == nil {
				log.Error("preload 为空。")
				goto Error
			}
			if preload, ok := body["preload"].(map[string]interface{}); ok {
				body = preload
			}
		}

		if body == nil {
			log.Error("token 为空")
			goto Error
		}

		track(c, body, reqm.FileName, reqm.UserID)
	}

	track(c, bodyMap, reqm.FileName, reqm.UserID)

Error:
	_, _ = c.Writer.Write([]byte("{\"error\":1}"))
}

// Convert 转换格式为open office xml格式
func Convert(c *gin.Context) {
	reqm := ReqConvert{}
	if err := c.Bind(&reqm); err != nil {
		c.WriteJSON(nil, status.InvalidArgument().WithError(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()

	fileURI := onlyoffice.Get().GetFileURI(reqm.FileName, user.UserID)
	fileExt := onlyoffice.Get().GetFileExtension(reqm.FileName, false)
	fileType := onlyoffice.Get().GetFileType(reqm.FileName)
	internalFileExt := onlyoffice.Get().GetInternalExtension(fileType)
	key := fileURI + fileutility.GetModTime(onlyoffice.Get().StoragePath(reqm.FileName, user.UserID)).Format("2006-01-02 15-04-05")
	key = onlyoffice.Get().GenerateRevisionID(key)

	//是已存换类型就直接返回
	if conf.Get().OnlyOffice.Document.ConvertedDocs.IndexOf(fileExt) != -1 {
		c.WriteJSON(reqm.FileName, nil)
		return
	}

	resp, err := onlyoffice.Get().GetConvertURI(fileURI, fileExt, internalFileExt, key, true)
	if err != nil {
		c.WriteJSON(nil, err)
		return
	}
	if resp.Error != 0 {
		c.WriteJSON(nil, status.Internal().WithError(onlyoffice.Get().ConvertURIErrorMessage(resp.Error)))
		return
	}

	file, err := http.Get(resp.FileURL)
	if err != nil {
		c.WriteJSON(nil, err)
		return
	}

	correctName := onlyoffice.Get().GetCorrectName(onlyoffice.Get().GetFileName(reqm.FileName, true)+internalFileExt, user.UserID)
	correctStoragePath := onlyoffice.Get().StoragePath(correctName, user.UserID)
	destf, err := os.Create(correctStoragePath)
	if err != nil {
		c.WriteJSON(nil, status.Internal().WithError(err))
		return
	}

	_, err = io.Copy(destf, file.Body)
	if err != nil {
		c.WriteJSON(nil, status.Internal().WithError(err))
		return
	}
	_ = destf.Close()

	storagePath := onlyoffice.Get().StoragePath(reqm.FileName, user.UserID)
	err = os.Remove(storagePath)
	if err != nil {
		c.WriteJSON(nil, status.Internal().WithError(err))
		return
	}

	historyPath := onlyoffice.Get().HistoryPath(reqm.FileName, user.UserID, true)
	correctHistoryPath := onlyoffice.Get().HistoryPath(correctName, user.UserID, true)
	err = os.Rename(historyPath, correctHistoryPath)
	if err != nil {
		c.WriteJSON(nil, status.Internal().WithError(err))
		return
	}

	createInfoPath := path.Join(historyPath, reqm.FileName+".txt")
	correctCreateInfoPath := path.Join(historyPath, correctName+".txt")
	err = os.Rename(createInfoPath, correctCreateInfoPath)
	if err != nil {
		c.WriteJSON(nil, status.Internal().WithError(err))
		return
	}

	c.WriteJSON(correctName, nil)
}

//Download 下载文件
func Download(c *gin.Context) {
	reqm := ReqDownload{}
	if err := c.Bind(&reqm); err != nil {
		c.WriteJSON(nil, status.InvalidArgument().WithError(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()

	filePath := onlyoffice.Get().ForcesavePath(reqm.FileName, user.UserID, false)
	if filePath == "" {
		filePath = onlyoffice.Get().StoragePath(reqm.FileName, user.UserID)
	}

	c.FileAttachment(filePath, reqm.FileName)

}

//DeleteFile 删除自己目录下的文件
func DeleteFile(c *gin.Context) {
	reqm := ReqDownload{}
	if err := c.Bind(&reqm); err != nil {
		c.WriteJSON(nil, status.InvalidArgument().WithError(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()

	if reqm.FileName != "" {
		storagePath := onlyoffice.Get().StoragePath(reqm.FileName, user.UserID)
		if err := os.Remove(storagePath); err != nil {
			c.WriteJSON(nil, status.Internal().WithError(err))
			return
		}

		historyPath := onlyoffice.Get().HistoryPath(reqm.FileName, user.UserID, true)
		if err := onlyoffice.Get().CleanFolderRecursive(historyPath, true); err != nil {
			c.WriteJSON(nil, status.Internal().WithError(err))
			return
		}
	} else {
		//delete all
		if err := onlyoffice.Get().CleanFolderRecursive(onlyoffice.Get().StoragePath("", user.UserID), false); err != nil {
			c.WriteJSON(nil, status.Internal().WithError(err))
			return
		}
	}

	c.WriteJSON(nil, nil)
}
