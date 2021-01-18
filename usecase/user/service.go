package account

import (
	"yumi/usecase/user/entity"
)

// Service ...
type Service struct {
}

// New a Service object
func New() (*Service, error) {
	return &Service{}, nil
}

// Create ...
func (s *Service) Create(req CreateAccountRequest) (err error) {
	data := GetData()

	ua := req.UserAttribute()
	u := entity.NewUser(&ua)
	err = u.LawEnforcement()
	if err != nil {
		return err
	}
	err = u.BcryptPassword()
	if err != nil {
		return err
	}

	// 持久化
	err = data.Create(ua)
	if err != nil {
		return err
	}
	return nil
}

// Disable ...
func (s *Service) Disable(req DisableAccountRequest) (err error) {
	return nil
}

// LoginByBcrypt ...
// 密码必须经过MD5加密
func (s *Service) LoginByBcrypt(req LoginByBcryptRequest) (LoginByBcryptResponse, error) {
	data := GetData()
	resp := LoginByBcryptResponse{}

	ut, err := data.Get(req.UserID)
	if err != nil {
		return resp, err
	}

	u := entity.NewUser(&ut)
	err = u.VerifyPassword(req.Password)
	if err != nil {
		return resp, err
	}

	//返回token

	return resp, nil
}
