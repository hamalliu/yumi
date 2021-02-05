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

func (m *FileInfo) attribute() (attr entity.MediaAttribute) {
	attr.Name = m.Name
	attr.Size = m.Size
	attr.Creator = m.Creator
	attr.Owner = m.Owner
	attr.OwnerType = m.OwnerType
	attr.Groups = m.Groups
	attr.Perm = m.Perm

	return
}

// CreateResponse is a param of Create function return
type CreateResponse struct {
	MediaUUID string
}
func (m *CreateResponse) setAttribute(attr entity.MediaAttribute) {
	m.MediaUUID = attr.MediaUUID
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
