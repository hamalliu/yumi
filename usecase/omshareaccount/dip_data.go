package omshareaccount

import "yumi/usecase/omshareaccount/entity"

var _data Data

// Data ...
type Data interface {
	ShareAccount(entity.ShareAccountAttribute) DataShareAccount
	GetShareAccount(shareID string) (DataShareAccount, error)
}

// DataShareAccount ...
type DataShareAccount interface {
	Attribute() *entity.ShareAccountAttribute
	
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
