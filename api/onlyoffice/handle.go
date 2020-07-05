package onlyoffice

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"yumi/pkg/log"
	"yumi/pkg/net/gin/binding"

	"yumi/api_model"
	"yumi/internal/onlyoffice"
	"yumi/pkg/conf"
	"yumi/pkg/ecode"
	"yumi/pkg/file_utility"
	"yumi/pkg/net/gin"
	"yumi/pkg/valuer"
)

func Index(c *gin.Context) {

}

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

	if mulfh.Size > conf.Get().OnlyOffice.Document.MaxFileSize.Size {
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

func Editor(c *gin.Context) {
	//TODO
}

func save() {

}

func forceSave() {

}

func track(c *gin.Context, body map[string]interface{}, fileName, userId string) {
	status := body["status"].(int64)
	if status == 1 {

	} else if status == 2 || status == 3 {

	} else if status == 6 || status == 7 {

	}

	_, _ = c.Writer.Write([]byte("{\"error\":0}"))
	return
}

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

func Convert(c *gin.Context) {
	reqm := api_model.ReqConvert{}
	if err := c.Bind(&reqm); err != nil {
		c.JSON(nil, ecode.ParamsErr(err))
		return
	}
	user := c.Get(valuer.KeyUser).User()

	fileUri := onlyoffice.Get().GetFileUri(reqm.FileName, user.UserId)
	fileExt := onlyoffice.Get().GetFileExtension(reqm.FileName, false)
	fileType := onlyoffice.Get().GetFileType(reqm.FileName)
	internalFileExt := onlyoffice.Get().GetInternalExtension(fileType)
	key := fileUri + file_utility.GetModTime(onlyoffice.Get().StoragePath(reqm.FileName, user.UserId)).Format("2006-01-02 15-04-05")
	key = onlyoffice.Get().GenerateRevisionId(key)

	//是已存换类型就直接返回
	if conf.Get().OnlyOffice.Document.ConvertedDocs.IndexOf(fileExt) != -1 {
		c.JSON(reqm.FileName, nil)
		return
	}

	resp, err := onlyoffice.Get().GetConvertUri(fileUri, fileExt, internalFileExt, key, true)
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
