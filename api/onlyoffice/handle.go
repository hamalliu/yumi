package onlyoffice

import (
	"io"
	"mime/multipart"
	"os"

	"yumi/api_model"
	"yumi/internal/onlyoffice"
	"yumi/pkg/conf"
	"yumi/pkg/ecode"
	"yumi/pkg/net/gin"
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

	userId := c.UserId()
	userName := c.UserName()
	fileName := onlyoffice.Get().GetCorrectName(mulfh.Filename, userId)
	ext := onlyoffice.Get().GetFileExtension(fileName, false)

	if onlyoffice.Get().AllowUploadExtension(ext) {
		c.JSON(nil, ecode.NoAllowExtension)
		return
	}

	storagePath := onlyoffice.Get().StoragePath(fileName, userId)
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

	if err := onlyoffice.Get().SaveCreateInfo(fileName, userId, userName); err != nil {
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

}

func Editor(c *gin.Context) {

}

func Track(c *gin.Context) {

}

func Convert(c *gin.Context) {

}

func Download(c *gin.Context) {

}

func DeleteFile(c *gin.Context) {

}
