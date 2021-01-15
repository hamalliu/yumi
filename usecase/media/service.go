package media

// Config ...
type Config struct {

}

// Service ...
type Service struct {
	conf Config
}

// New a Service object
func New(options ...ServiceOption) (*Service, error) {
	do := serviceOptions{
	}
	for _, option := range options {
		option.f(&do)
	}

	return nil, nil
}

// Create a media file
func (s *Service) Create(f FileInfo, operator string) (resp CreateResponse, err error) {
	return
}

// BatchCreate media files
func (s *Service) BatchCreate(fs []FileInfo, operator string) (resp BatchCreateResponse, err error) {
	return
}

// Get a media file
func (s *Service) Get(fileNo interface{}, operator string) (resp GetResponse, err error) {
	return
}

// List Get media file list
func (s *Service) List(page, line int, operator string) (err error) {
	return
}
