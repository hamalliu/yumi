package media

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"yumi/conf"
	"yumi/pkg/ecode"
	"yumi/gin"
	"yumi/usecase/media/data"
)

//UploadMultipart 多文件上传
func UploadMultipart(c *gin.Context) {
	req := c.Request
	var (
		err error
	)

	if err = req.ParseMultipartForm(conf.Get().Media.MultipleFileUploadsMaxSize.Size()); err != nil {
		c.JSON(nil, ecode.FileSizeTooBig)
		return
	}

	// 检查文件大小
	fds := req.MultipartForm.File["file[]"]
	l := len(fds)
	for i := 0; i < l; i++ {
		if fds[i].Size > conf.Get().Media.SingleFileUploadsMaxSize.Size() {
			c.JSON(nil, ecode.FileSizeTooBig)
			return
		}
	}

	for i := 0; i < l; i++ {
		var (
			osf  *os.File
			mulf multipart.File
		)
		if mulf, err = fds[i].Open(); err != nil {
			c.JSON(nil, ecode.ServerErr(err))
			return
		}

		// 创建一个不重复的文件名，复制文件
		suffix := fds[i].Filename[strings.LastIndex(fds[i].Filename, ".")+1:]
		name := fmt.Sprintf("%d.%s", time.Now().UnixNano(), suffix)
		path := fmt.Sprintf("%s/%s", conf.Get().Media.StoragePath, name)
		for {
			if _, err := os.Stat(path); os.IsExist(err) {
				path = fmt.Sprintf("%s/%d.%s", conf.Get().Media.StoragePath, time.Now().UnixNano(), suffix)
			} else if os.IsNotExist(err) {
				break
			}
		}
		if osf, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0744); err != nil {
			c.JSON(nil, ecode.ServerErr(err))
			return
		}

		if _, err := io.Copy(osf, mulf); err != nil {
			c.JSON(nil, ecode.ServerErr(err))
			return
		}
		_ = mulf.Close()
		_ = osf.Close()

		// 添加上传记录
		operatorid := req.Header.Get("xuid")
		operator := req.Header.Get("username")
		if _, err = data.DB().Insert(suffix, name, fds[i].Filename, path, operator, operatorid); err != nil {
			c.JSON(nil, ecode.ServerErr(err))
			return
		}
	}
}

//Upload 单文件上传
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

	// 检查文件大小
	if mulfh.Size > conf.Get().Media.SingleFileUploadsMaxSize.Size() {
		c.JSON(nil, ecode.FileSizeTooBig)
		return
	}

	// 创建一个不重复的文件名，复制文件
	suffix := mulfh.Filename[strings.LastIndex(mulfh.Filename, ".")+1:]
	name := fmt.Sprintf("%d.%s", time.Now().UnixNano(), suffix)
	path := fmt.Sprintf("%s/%s", conf.Get().Media.StoragePath, name)
	for {
		if _, err := os.Stat(path); os.IsExist(err) {
			path = fmt.Sprintf("%s/%d.%s", conf.Get().Media.StoragePath, time.Now().UnixNano(), suffix)
		} else if os.IsNotExist(err) {
			break
		}
	}
	if osf, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0744); err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}

	if _, err := io.Copy(osf, mulf); err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}
	_ = mulf.Close()
	_ = osf.Close()

	// 添加上传记录
	operatorid := req.Header.Get("xuid")
	operator := req.Header.Get("username")
	if _, err = data.DB().Insert(suffix, name, mulfh.Filename, path, operator, operatorid); err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}
}
