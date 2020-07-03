package onlyoffice

import (
	"io"
	"mime/multipart"
	"os"
	"path"

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

}

func Track(c *gin.Context) {

}

func Convert(c *gin.Context) {

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

}
