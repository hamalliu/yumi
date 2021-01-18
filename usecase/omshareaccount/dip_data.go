package omshareaccount

import (
	"yumi/usecase/omshareaccount/entity"
)

var _data Data

// Data ...
type Data interface {
	Create(saa entity.ShareAccountAttribute) error
	Update(sa entity.ShareAccountAttribute) error
	Get(shareID string) (entity.ShareAccountAttribute, error)
}

// InitData ...
func InitData(data Data) {
	_data = data
}

// GetData ...
func GetData() Data {
	return _data
}
