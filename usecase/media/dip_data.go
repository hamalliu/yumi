package media

import (
	"yumi/usecase/media/entity"
)

var _data Data

// Data ...
type Data interface {
	Create(ma entity.MediaAttribute) error
	Get(fileNo string) (entity.MediaAttribute, error)
	List(page, line int) ([]entity.MediaAttribute, error)
}

// InitData ...
func InitData(data Data) {
	_data = data
}

// GetData ...
func GetData() Data {
	return _data
}
