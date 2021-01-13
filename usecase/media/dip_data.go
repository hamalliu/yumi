package media

var _data Data

// Data ...
type Data interface {
	CreateMedia(media Media) error
	GetMedia(fileNo string) (Media, error)
}

// InitData ...
func InitData(data Data) {
	_data = data
}

// GetData ...
func GetData() Data {
	return _data
}
