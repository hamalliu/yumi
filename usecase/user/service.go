package user

import (
	"yumi/pkg/status"
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
func (s *Service) Create(req CreateRequest) (err error) {
	data := GetData()

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
	err = data.Create(ua)
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
	data := GetData()
	resp := LoginByBcryptResponse{}

	exist, attr, err := data.Exist(entity.UserAttributeIDs{UserID: req.UserID})
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
	sessID, err := u.Session(data.GetSessionsStore(), req.UserID, req.Password, req.Client)
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
