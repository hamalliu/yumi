package entity

import (
	"io"
	"mime/multipart"
	"os"
)

// OwnerType involves individuals and organizations
type OwnerType int

const (
	// OwnerTypeIndividual indicates that the owner type is an individual
	OwnerTypeIndividual OwnerType = iota
	// OwnerTypeOrganization indicates that the owner type is an Organization
	OwnerTypeOrganization
)

// MediaAttribute is attribute of media struct
type MediaAttribute struct {
	MediaUUID   string
	Name        string
	Size        int64
	StoragePath StoragePath

	Creator   string
	Owner     string
	OwnerType OwnerType
	Groups    []string
	Perm      os.FileMode
}

// Media is value object
type Media struct {
	attr *MediaAttribute
}

// NewMedia a media
func NewMedia(attr *MediaAttribute) Media {
	return Media{attr: attr}
}

// Create ...
func (m *Media) Create(file multipart.File) error {
	return nil
}

// GetFile return a io.reader
func (m *Media) GetFile() io.ReadCloser {
	return nil
}

// SaveFileFormHTTP ...
func (m *Media) SaveFileFormHTTP(file multipart.File) error {
	return nil
}

