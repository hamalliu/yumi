package service

import "yumi/usecase/media/entity"

// Config ...
type Config struct{}

// Service ...
type Service struct {
	data Data
}

// New a Service object
func New(data Data, options ...Option) (*Service, error) {
	do := Options{}
	for _, option := range options {
		option.f(&do)
	}

	return &Service{data: data}, nil
}

// Create a media file
func (s *Service) Create(f FileInfo) (resp CreateResponse, err error) {
	attr := f.attribute()
	mda := entity.NewMedia(&attr)
	err = mda.Create(f.File)
	if err != nil {
		return
	}

	// 持久化
	err = s.data.Create(attr)
	if err != nil {
		return
	}

	resp.setAttribute(attr)
	return
}

// BatchCreate media files
func (s *Service) BatchCreate(fs []FileInfo) (resp BatchCreateResponse, err error) {
	return
}

// Get a media file
func (s *Service) Get(fileNo interface{}) (resp GetResponse, err error) {
	return
}

// List Get media file list
func (s *Service) List(page, line int) (err error) {
	return
}
