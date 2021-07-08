package service

import (
	"yumi/pkg/status"
	"yumi/usecase/user/entity"
)

// Service ...
type Service struct {
	data Data
}

// New a Service object
func New(data Data) (*Service, error) {
	return &Service{data: data}, nil
}

// Create ...
func (s *Service) Create(req CreateRequest) (err error) {
	ua := req.userAttribute()
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
	err = s.data.Create(ua)
	if err != nil {
		return err
	}
	return nil
}

// Disable ...
func (s *Service) Disable(req DisableRequest) (err error) {
	// TODO
	return nil
}

// LoginByBcrypt ...
// 密码必须经过MD5加密
func (s *Service) LoginByBcrypt(req LoginByBcryptRequest) (LoginByBcryptResponse, error) {
	resp := LoginByBcryptResponse{}

	exist, attr, err := s.data.Exist(entity.UserAttributeIDs{UserID: req.UserID})
	if err != nil {
		return resp, err
	}
	if !exist {
		return resp, status.NotFound().WithMessage(entity.UserNotFound)
	}

	u := entity.NewUser(&attr)
	err = u.VerifyPassword(req.Password)
	if err != nil {
		return resp, err
	}

	//构建session
	sessID, err := u.Session(s.data.GetSessionsStore(), req.UserID, req.Password, req.Client)
	if err != nil {
		return resp, err
	}

	resp.UserID = req.UserID
	resp.SecureKey = sessID

	return resp, nil
}

// Authenticate ...
func (s *Service) Authenticate(req AuthentcateRequest) error {
	// TODO
	return nil
}
