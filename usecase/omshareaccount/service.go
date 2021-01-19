package omshareaccount

import "yumi/usecase/omshareaccount/entity"

// Service ...
type Service struct {
}

// New a Service object
func New() (*Service, error) {
	return &Service{}, nil
}

// Create ...
func (s *Service) Create(req CreateShareRequest) (err error) {
	data := GetData()

	saa := req.ShareAccountAttribute()
	parentShareID := saa.Share.ParentShareID
	if parentShareID != "" {
		dataSas, err := data.GetShareAccount(parentShareID)
		if err != nil {
			return err
		}
		saa.Parent = dataSas.Attribute()
	}

	// 执法：检查当前对象数据关系是否合乎业务规定
	sa := entity.NewShareAccount(&saa)
	err = sa.LawEnforcement()
	if err != nil {
		return err
	}

	// 持久化
	err = data.ShareAccount(saa).Create()
	if err != nil {
		return err
	}
	return nil
}

// GetShare ...
func (s *Service) GetShare() (resp GetShareResponse, err error) {
	return
}

// CancelShare ...
func (s *Service) CancelShare(req CancelShareRequest) error {
	data := GetData()

	dataSas, err := data.GetShareAccount(req.ShareID)
	if err != nil {
		return err
	}
	sa := entity.NewShareAccount(dataSas.Attribute())
	err = sa.SetCancellationMsg()
	if err != nil {
		return err
	}

	// 持久化
	err = dataSas.Update()
	if err != nil {
		return err
	}

	return nil
}

// ReceiveAccount ...
func (s *Service) ReceiveAccount(req ReceiveAccountRequest) (ReceiveAccountResponse, error) {
	data := GetData()
	resp := ReceiveAccountResponse{}

	dataSas, err := data.GetShareAccount(req.ShareID)
	if err != nil {
		return resp, err
	}
	sa := entity.NewShareAccount(dataSas.Attribute())

	acct := getAcctName()
	err = sa.SetReceived(acct)
	if err != nil {
		return resp, err
	}
	caresp, err := s.createAccount(acct)
	if err != nil {
		return resp, err
	}
	resp.SetAcct(caresp)

	// 持久化
	err = dataSas.Update()
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *Service) createAccount(acct string) (createAccountResponse, error) {
	// TODO:
	return createAccountResponse{}, nil
}

func getAcctName() string {
	// TODO:
	return ""
}
