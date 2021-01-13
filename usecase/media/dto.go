package media

import (
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

// FileInfo contains upload file info
type FileInfo struct {
	Name string
	Size int64
	File multipart.File

	Creator   string
	Owner     string
	OwnerType OwnerType
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
