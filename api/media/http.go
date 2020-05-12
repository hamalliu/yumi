package media

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"yumi/conf"
	"yumi/internal/db"
	"yumi/pkg/net/ymhttp"
	"yumi/response"
)

func UploadMultipart(ctx *ymhttp.Context) {
	req := ctx.Request
	resp := ctx.Writer
	var (
		err error
	)

	if err = req.ParseMultipartForm(conf.Get().MaxFileSize); err != nil {
		response.Json(resp, req, response.FileSizeTooBig(), nil)
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
			response.Json(resp, req, response.InternalError(err), nil)
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
			response.Json(resp, req, response.InternalError(err), nil)
			return
		}

		if _, err := io.Copy(osf, mulf); err != nil {
			response.Json(resp, req, response.InternalError(err), nil)
			return
		}
		_ = mulf.Close()
		_ = osf.Close()

		operatorid := req.Header.Get("xuid")
		operator := req.Header.Get("username")
		if _, err = db.Media().Add(suffix, name, fds[i].Filename, path, operator, operatorid); err != nil {
			response.Json(resp, req, response.InternalError(err), nil)
			return
		}
	}
}

func Upload(ctx *ymhttp.Context) {
	req := ctx.Request
	resp := ctx.Writer
	var (
		mulf  multipart.File
		mulfh *multipart.FileHeader
		osf   *os.File

		err error
	)

	if mulf, mulfh, err = req.FormFile("file"); err != nil {
		response.Json(resp, req, response.InternalError(err), nil)
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
		response.Json(resp, req, response.InternalError(err), nil)
		return
	}

	if _, err := io.Copy(osf, mulf); err != nil {
		response.Json(resp, req, response.InternalError(err), nil)
		return
	}
	_ = mulf.Close()
	_ = osf.Close()

	operatorid := req.Header.Get("xuid")
	operator := req.Header.Get("username")
	if _, err = db.Media().Add(suffix, name, mulfh.Filename, path, operator, operatorid); err != nil {
		response.Json(resp, req, response.InternalError(err), nil)
		return
	}
}
