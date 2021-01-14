package media

import (
	"mime/multipart"
	"os"

	"yumi/usecase/media/entity"
)

// FileInfo contains upload file info
type FileInfo struct {
	Name string
	Size int64
	File multipart.File

	Creator   string
	Owner     string
	OwnerType entity.OwnerType
	Groups    []string
	Perm      os.FileMode
}

// CreateResponse is a param of Create function return
type CreateResponse struct {
}

// BatchCreateResponse is a param of BatchCreate function return
type BatchCreateResponse struct {
}

// GetResponse is a param of Get function return
type GetResponse struct {
}

// ListResponse is a param of List function return
type ListResponse struct {
}
