package user

import (
	"yumi/usecase/user/entity"
	"yumi/pkg/sessions"
)

var _data Data

// Data ...
type Data interface {
	User(entity.UserAttribute) DataUser
	GetUser(userID string) (DataUser, error)

	GetSessionsStore() sessions.Store
}

// DataUser ...
type DataUser interface{
	Attribute() *entity.UserAttribute

	Create() error
	Update() error
}

// InitData ...
func InitData(data Data) {
	_data = data
}

// GetData ...
func GetData() Data {
	return _data
}
