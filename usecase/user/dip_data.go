package user

import (
	"yumi/usecase/user/entity"
	"yumi/pkg/sessions"
)

var _data Data

// Data ...
type Data interface {
	Create(entity.UserAttribute) error
	Update(entity.UserAttribute) error
	GetUser(userID string) (entity.UserAttribute, error)

	GetSessionsStore() sessions.Store
}

// InitData ...
func InitData(data Data) {
	_data = data
}

// GetData ...
func GetData() Data {
	return _data
}
