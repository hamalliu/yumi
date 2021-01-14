package media

import (
	"yumi/usecase/media/entity"
)

var _data Data

// Data ...
type Data interface {
	CreateMedia(ma entity.MediaAttribute) error
	GetMedia(fileNo string) (entity.MediaAttribute, error)
	ListMedia(page, line int) ([]entity.MediaAttribute, error)
}

// InitData ...
func InitData(data Data) {
	_data = data
}

// GetData ...
func GetData() Data {
	return _data
}
