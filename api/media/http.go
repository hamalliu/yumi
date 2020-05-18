package media

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"yumi/internal/db"
	"yumi/pkg/conf"
	"yumi/pkg/ecode"
	"yumi/pkg/net/ymhttp"
)

func UploadMultipart(c *ymhttp.Context) {
	req := c.Request
	var (
		err error
	)

	if err = req.ParseMultipartForm(conf.Get().MaxFileSize); err != nil {
		c.JSON(nil, ecode.FileSizeTooBig)
		return
	}

	fds := req.MultipartForm.File["file[]"]
	l := len(fds)
	for i := 0; i < l; i++ {
		var (
			osf  *os.File
			mulf multipart.File
		)
		if mulf, err = fds[i].Open(); err != nil {
			c.JSON(nil, ecode.ServerErr(err))
			return
		}

		suffix := fds[i].Filename[strings.LastIndex(fds[i].Filename, ".")+1:]
		name := fmt.Sprintf("%d.%s", time.Now().UnixNano(), suffix)
		path := fmt.Sprintf("%s/%s", conf.Get().StoragePath, name)
		for {
			if _, err := os.Stat(path); os.IsExist(err) {
				path = fmt.Sprintf("%s/%d.%s", conf.Get().StoragePath, time.Now().UnixNano(), suffix)
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

		operatorid := req.Header.Get("xuid")
		operator := req.Header.Get("username")
		if _, err = db.Media().Add(suffix, name, fds[i].Filename, path, operator, operatorid); err != nil {
			c.JSON(nil, ecode.ServerErr(err))
			return
		}
	}
}

func Upload(c *ymhttp.Context) {
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

	suffix := mulfh.Filename[strings.LastIndex(mulfh.Filename, ".")+1:]
	name := fmt.Sprintf("%d.%s", time.Now().UnixNano(), suffix)
	path := fmt.Sprintf("%s/%s", conf.Get().StoragePath, name)
	for {
		if _, err := os.Stat(path); os.IsExist(err) {
			path = fmt.Sprintf("%s/%d.%s", conf.Get().StoragePath, time.Now().UnixNano(), suffix)
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

	operatorid := req.Header.Get("xuid")
	operator := req.Header.Get("username")
	if _, err = db.Media().Add(suffix, name, mulfh.Filename, path, operator, operatorid); err != nil {
		c.JSON(nil, ecode.ServerErr(err))
		return
	}
}
