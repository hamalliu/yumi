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

	"yumi/api_model"
	"yumi/internal/onlyoffice"
	"yumi/pkg/conf"
	"yumi/pkg/ecode"
	"yumi/pkg/file_utility"
	"yumi/pkg/log"
	"yumi/pkg/net/gin"
	"yumi/pkg/net/gin/binding"
	"yumi/pkg/valuer"
)

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
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	if mulfh.Size > conf.Get().OnlyOffice.Document.MaxFileSize.Size() {
		c.JSON(nil, ecode.FileSizeTooBig)
		return
	}

	user := c.Get(valuer.KeyUser).User()
	fileName := onlyoffice.Get().GetCorrectName(mulfh.Filename, user.UserId)
	ext := onlyoffice.Get().GetFileExtension(fileName, false)

	if onlyoffice.Get().AllowUploadExtension(ext) {
		c.JSON(nil, ecode.NoAllowExtension)
		return
	}

	storagePath := onlyoffice.Get().StoragePath(fileName, user.UserId)
	if osf, err = os.OpenFile(storagePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0744); err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}
	if _, err := io.Copy(osf, mulf); err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}
	_ = mulf.Close()
	_ = osf.Close()

	if err := onlyoffice.Get().SaveCreateInfo(fileName, user.UserId, user.UserName); err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	c.JSONNoDataSuccess()
	return
}

//Sample 从样品中复制文件，用于新建文件
func Sample(c *gin.Context) {
	reqm := api_model.ReqSample{}
	if err := c.Bind(&reqm); err != nil {
		c.JSON(nil, ecode.ParamsErr(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()

	reqm.FileName = onlyoffice.Get().GetCorrectName(reqm.FileName+reqm.FileExtension, user.UserId)
	filePath := onlyoffice.Get().StoragePath(reqm.FileName, user.UserId)
	sampleFile := ""
	if reqm.SampleName == "" {
		reqm.SampleName = "new"
	}
	sampleFile = path.Join(conf.Get().OnlyOffice.Document.SamplesPath, reqm.SampleName+reqm.FileExtension)
	if err := file_utility.CopyFile(sampleFile, filePath); err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	if err := onlyoffice.Get().SaveCreateInfo(reqm.FileName, user.UserId, user.UserName); err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	c.JSONNoDataSuccess()
	return
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
	log.Error2(err)
	t, _ := template.ParseFiles("./views/onlyoffice/editor.tpl")
	_ = t.Execute(c.Writer, struct{ message string }{message: "服务器内部错误"})
}

//Editor 按配置加载数据，返回office编辑模板
func Editor(c *gin.Context) {
	reqm := api_model.ReqEditor{}
	if err := c.Bind(&reqm); err != nil {
		c.JSON(nil, ecode.ParamsErr(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()
	confOo := conf.Get().OnlyOffice
	history, historyData, version, err := onlyoffice.Get().GetHistory(reqm.FileName, user.UserId)
	if err != nil {
		renderError(c, err)
	}

	editor := onlyoffice.Editor{
		File: onlyoffice.File{
			Name:    reqm.FileName,
			Ext:     onlyoffice.Get().GetFileExtension(reqm.FileName, true),
			Uri:     onlyoffice.Get().GetFileUri(reqm.FileName, user.UserId),
			Key:     onlyoffice.Get().GenerateKey(reqm.FileName, user.UserId),
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
		GobackUrl:    "",
		Mode:         reqm.Mode,
		CallbackUrl:  getCallbackURL(c, reqm.FileName, user.UserId),
		UserId:       user.UserId,
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
		key := onlyoffice.Get().GenerateRevisionId(downloadURI)

		resp, err := onlyoffice.Get().GetConvertUri(downloadURI, downloadExt, curExt, key, true)
		if err != nil {
			log.Error(err)
			goto Error
		}
		if resp.Error != 0 {
			err := onlyoffice.Get().ConvertUriErrorMessage(resp.Error)
			log.Error(err)
			goto Error
		}
		save(c, body, resp.FileUrl, fileName, userID)
	} else {
		storagePath := onlyoffice.Get().StoragePath(fileName, userID)

		if file_utility.ExistFile(storagePath) {
			historyPath := onlyoffice.Get().HistoryPath(fileName, userID, false)
			if historyPath == "" {
				historyPath = onlyoffice.Get().HistoryPath(fileName, userID, true)
				err := file_utility.CreateDir(historyPath)
				if err != nil {
					log.Error(err)
					goto Error
				}
			}

			countVersion := onlyoffice.Get().CountHistoryVersion(fileName, userID)
			version := countVersion + 1
			versionPath := onlyoffice.Get().VersionPath(fileName, userID, version)
			err := file_utility.CreateDir(versionPath)
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
				err = file_utility.WriteFile(resp.Body, diffPath)
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
			err = file_utility.WriteFile(bytes.NewBuffer([]byte(changesStr)), changesPath)
			if err != nil {
				log.Error(err)
				goto Error
			}

			//key
			key := body["key"].(string)
			pathKey := onlyoffice.Get().KeyPath(fileName, userID, version)
			err = file_utility.WriteFile(bytes.NewBuffer([]byte(key)), pathKey)
			if err != nil {
				log.Error(err)
				goto Error
			}

			//prev
			prevPath := onlyoffice.Get().PrevFilePath(fileName, userID, version)
			err = file_utility.CopyFile(storagePath, prevPath)
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
			err = file_utility.WriteFile(resp.Body, storagePath)
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
	c.Writer.Write([]byte("{\"error\":1}"))
	return
}

func forceSave(c *gin.Context, body map[string]interface{}, downloadURI, fileName, userID string) {
	curExt := onlyoffice.Get().GetFileExtension(fileName, false)
	downloadExt := onlyoffice.Get().GetFileExtension(downloadURI, false)

	if curExt != downloadExt {
		key := onlyoffice.Get().GenerateRevisionId(downloadURI)

		resp, err := onlyoffice.Get().GetConvertUri(downloadURI, downloadExt, curExt, key, true)
		if err != nil {
			log.Error(err)
			goto Error
		}
		if resp.Error != 0 {
			err := onlyoffice.Get().ConvertUriErrorMessage(resp.Error)
			log.Error(err)
			goto Error
		}
		forceSave(c, body, resp.FileUrl, fileName, userID)
	} else {
		forceSavePath := onlyoffice.Get().ForcesavePath(fileName, userID, false)
		if forceSavePath == "" {
			forceSavePath = onlyoffice.Get().ForcesavePath(fileName, userID, true)
			err := file_utility.CreateDir(forceSavePath)
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

		err = file_utility.WriteFile(resp.Body, forceSavePath)
		if err != nil {
			log.Error(err)
			goto Error
		}
	}
	_, _ = c.Writer.Write([]byte("{\"error\":0}"))
	return

Error:
	c.Writer.Write([]byte("{\"error\":1}"))
	return
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
	return
}

//Track onlyoffice的回调处理
func Track(c *gin.Context) {
	reqm := api_model.ReqTrack{}
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

		track(c, body, reqm.FileName, reqm.UserId)
	}

	track(c, bodyMap, reqm.FileName, reqm.UserId)

Error:
	c.Writer.Write([]byte("{\"error\":1}"))
	return
}

//Convert 转换格式为open office xml格式
func Convert(c *gin.Context) {
	reqm := api_model.ReqConvert{}
	if err := c.Bind(&reqm); err != nil {
		c.JSON(nil, ecode.ParamsErr(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()

	fileURI := onlyoffice.Get().GetFileUri(reqm.FileName, user.UserId)
	fileExt := onlyoffice.Get().GetFileExtension(reqm.FileName, false)
	fileType := onlyoffice.Get().GetFileType(reqm.FileName)
	internalFileExt := onlyoffice.Get().GetInternalExtension(fileType)
	key := fileURI + file_utility.GetModTime(onlyoffice.Get().StoragePath(reqm.FileName, user.UserId)).Format("2006-01-02 15-04-05")
	key = onlyoffice.Get().GenerateRevisionId(key)

	//是已存换类型就直接返回
	if conf.Get().OnlyOffice.Document.ConvertedDocs.IndexOf(fileExt) != -1 {
		c.JSON(reqm.FileName, nil)
		return
	}

	resp, err := onlyoffice.Get().GetConvertUri(fileURI, fileExt, internalFileExt, key, true)
	if err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}
	if resp.Error != 0 {
		c.JSON(nil, ecode.ServerErr(onlyoffice.Get().ConvertUriErrorMessage(resp.Error)))
		return
	}

	file, err := http.Get(resp.FileUrl)
	if err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	correctName := onlyoffice.Get().GetCorrectName(onlyoffice.Get().GetFileName(reqm.FileName, true)+internalFileExt, user.UserId)
	correctStoragePath := onlyoffice.Get().StoragePath(correctName, user.UserId)
	destf, err := os.Create(correctStoragePath)
	if err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	_, err = io.Copy(destf, file.Body)
	if err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}
	_ = destf.Close()

	storagePath := onlyoffice.Get().StoragePath(reqm.FileName, user.UserId)
	err = os.Remove(storagePath)
	if err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	historyPath := onlyoffice.Get().HistoryPath(reqm.FileName, user.UserId, true)
	correctHistoryPath := onlyoffice.Get().HistoryPath(correctName, user.UserId, true)
	err = os.Rename(historyPath, correctHistoryPath)
	if err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	createInfoPath := path.Join(historyPath, reqm.FileName+".txt")
	correctCreateInfoPath := path.Join(historyPath, correctName+".txt")
	err = os.Rename(createInfoPath, correctCreateInfoPath)
	if err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	c.JSON(correctName, nil)
	return
}

//Download 下载文件
func Download(c *gin.Context) {
	reqm := api_model.ReqDownload{}
	if err := c.Bind(&reqm); err != nil {
		c.JSON(nil, ecode.ParamsErr(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()

	filePath := onlyoffice.Get().ForcesavePath(reqm.FileName, user.UserId, false)
	if filePath == "" {
		filePath = onlyoffice.Get().StoragePath(reqm.FileName, user.UserId)
	}

	c.FileAttachment(filePath, reqm.FileName)

	return
}

//DeleteFile 删除自己目录下的文件
func DeleteFile(c *gin.Context) {
	reqm := api_model.ReqDownload{}
	if err := c.Bind(&reqm); err != nil {
		c.JSON(nil, ecode.ParamsErr(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()

	if reqm.FileName != "" {
		storagePath := onlyoffice.Get().StoragePath(reqm.FileName, user.UserId)
		if err := os.Remove(storagePath); err != nil {
			c.JSON(nil, ecode.ServerErr(err))
			return
		}

		historyPath := onlyoffice.Get().HistoryPath(reqm.FileName, user.UserId, true)
		if err := onlyoffice.Get().CleanFolderRecursive(historyPath, true); err != nil {
			c.JSON(nil, ecode.ServerErr(err))
			return
		}
	} else {
		//delete all
		if err := onlyoffice.Get().CleanFolderRecursive(onlyoffice.Get().StoragePath("", user.UserId), false); err != nil {
			c.JSON(nil, ecode.ServerErr(err))
			return
		}
	}

	c.JSONNoDataSuccess()
	return
}
