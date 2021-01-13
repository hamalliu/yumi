package entity

import (
	"io"
	"mime/multipart"
	"os"

	"yumi/usecase/media"
)

// Media is value object
type Media struct {
	Name        string
	Size        int64
	StoragePath StoragePath

	Creator   string
	Owner     string
	OwnerType media.OwnerType
	Groups    []string
	Perm      os.FileMode
}

// StoragePath refarence media file storage path
type StoragePath struct {
	Host string
	Post int
	Path string
}

// GetFile return a io.reader
func (m *Media) GetFile() io.ReadCloser {
	return nil
}

// SaveFileFormHTTP ...
func (m *Media) SaveFileFormHTTP(file multipart.File) error {
	return nil
}
