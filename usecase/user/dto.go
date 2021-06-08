package user

import "yumi/usecase/user/entity"

// CreateRequest ...
type CreateRequest struct {
	UserID      string
	Password    string
	UserName    string
	PhoneNumber string
}

// userAttribute ...
func (ar *CreateRequest) userAttribute() entity.UserAttribute {
	return entity.UserAttribute{
		UserID:      ar.UserID,
		Password:    ar.Password,
		UserName:    ar.UserName,
		PhoneNumber: ar.PhoneNumber,
	}
}

// LoginByBcryptRequest ...
type LoginByBcryptRequest struct {
	Client   string
	UserID   string
	Password string
}

// LoginByBcryptResponse ...
type LoginByBcryptResponse struct {
	UserID            string
	SecureKey         string
}

// DisableRequest ...
type DisableRequest struct {
}

// AuthentcateRequest ...
type AuthentcateRequest struct {
}
