package media

// Config ...
type Config struct {

}

// Media ...
type Media struct {
	conf Config
}

// New a Media object
func New(conf Config) (*Media, error) {
	return nil, nil
}

// Create a media file
func (m *Media) Create(f FileInfo, operator string) (resp CreateResponse, err error) {
	return
}

// BatchCreate media files
func (m *Media) BatchCreate(fs []FileInfo, operator string) (resp BatchCreateResponse, err error) {
	return
}

// Delete a media file
func (m *Media) Delete(fileID interface{}, operator string) (err error) {
	return
}

// BatchDelete media files
func (m *Media) BatchDelete(fileIDs []interface{}, operator string) (err error) {
	return
}

// Replase a media file
func (m *Media) Replase(srcFileID interface{}, dest FileInfo, operator string) (err error) {
	return
}

// List Get media file list
func (m *Media) List(page, line int, operator string) (err error) {
	return
}


