package media

// Config ...
type Config struct {

}

// Service ...
type Service struct {
	conf Config
}

// New a Service object
func New(conf Config) (*Service, error) {
	return nil, nil
}

// Create a media file
func (m *Service) Create(f FileInfo, operator string) (resp CreateResponse, err error) {
	return
}

// BatchCreate media files
func (m *Service) BatchCreate(fs []FileInfo, operator string) (resp BatchCreateResponse, err error) {
	return
}

// Get a media file
func (m *Service) Get(fileNo interface{}, operator string) (resp GetResponse, err error) {
	return
}

// List Get media file list
func (m *Service) List(page, line int, operator string) (err error) {
	return
}
