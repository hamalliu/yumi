package omshareaccount

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
		saa.Parent, err = data.Get(parentShareID)
		if err != nil {
			return err
		}
	}

	// 持久化
	err = data.Create(saa)
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
func (s *Service) CancelShare() error {
	return nil
}

// ReceiveAccount ...
func (s *Service) ReceiveAccount() error {
	return nil
}
