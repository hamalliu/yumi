package account

import "yumi/usecase/user/entity"

var _data Data

// Data ...
type Data interface {
	Create(saa entity.UserAttribute) error
	Update(sa entity.UserAttribute) error
	Get(userID string) (entity.UserAttribute, error)
}

// InitData ...
func InitData(data Data) {
	_data = data
}

// GetData ...
func GetData() Data {
	return _data
}
