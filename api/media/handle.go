package media

import (
	"net/http"
	"yumi/conf"
	"yumi/gin"
	"yumi/status"
	"yumi/usecase/media"
)

var mediaSrv *media.Service

func init() {
	var err error
	mediaSrv, err = media.New()
	if err != nil {
		panic(err)
	}
}

//UploadMultiple 多文件上传
func UploadMultiple(c *gin.Context) {
	req := c.Request
	var (
		err error
	)

	mediaConf := conf.Get().Media
	if err = req.ParseMultipartForm(mediaConf.MultipleFileUploadsMaxSize.Size()); err != nil {
		if err == http.ErrLineTooLong {
			c.WriteJSON(nil, status.FailedPrecondition().WithMessage(status.FileIsTooLarge))
		} else {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(err))
		}
		return
	}

	fs := []media.FileInfo{}
	fds := req.MultipartForm.File["file[]"]
	l := len(fds)
	for i := 0; i < l; i++ {
		mulf, err := fds[i].Open()
		if err != nil {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(err))
			return
		}

		f := media.FileInfo{}
		f.Name = fds[i].Filename
		f.Size = fds[i].Size
		f.File = mulf
		f.Creator = ""
		f.Owner = ""
		f.OwnerType = 0
		f.Groups = nil
		f.Perm = 0444

		fs = append(fs, f)
	}

	resp, err := mediaSrv.BatchCreate(fs)
	c.WriteJSON(resp, err)
}

//Upload 单文件上传
func Upload(c *gin.Context) {
	req := c.Request

	mulf, mulfh, err := req.FormFile("file")
	if err != nil {
		if err == http.ErrLineTooLong {
			c.WriteJSON(nil, status.FailedPrecondition().WithMessage(status.FileIsTooLarge))
		} else {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(err))
		}
	}

	f := media.FileInfo{}
	f.Name = mulfh.Filename
	f.Size = mulfh.Size
	f.File = mulf
	f.Creator = ""
	f.Owner = ""
	f.OwnerType = 0
	f.Groups = nil
	f.Perm = 0444

	resp, err := mediaSrv.Create(f)
	c.WriteJSON(resp, err)
}
