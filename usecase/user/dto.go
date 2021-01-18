package account

import "yumi/usecase/user/entity"

// CreateAccountRequest ...
type CreateAccountRequest struct {
	UserID      string
	Password    string
	UserName    string
	PhoneNumber string
}

// UserAttribute ...
func (ar *CreateAccountRequest) UserAttribute() entity.UserAttribute {
	return entity.UserAttribute{
		UserID:      ar.UserID,
		Password:    ar.Password,
		UserName:    ar.UserName,
		PhoneNumber: ar.PhoneNumber,
	}
}

// LoginByBcryptRequest ...
type LoginByBcryptRequest struct {
	UserID      string
	Password    string
}

// LoginByBcryptResponse ...
type LoginByBcryptResponse struct {
	
}

// DisableAccountRequest ...
type DisableAccountRequest struct {
}
