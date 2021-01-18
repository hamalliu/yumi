package omshareaccount

import (
	"yumi/entities"
	"yumi/pkg/types"
	"yumi/usecase/omshareaccount/entity"
)

// CreateShareRequest ...
type CreateShareRequest struct {
	ParentShareID string
	Sender        string
	AccountNum    int

	// Message
	entities.ShareMessage
}

// ShareAccountAttribute ...
func (c *CreateShareRequest) ShareAccountAttribute() entity.ShareAccountAttribute {
	shareID := entities.NewShareID()
	sendTime := types.NowTimestamp()

	c.ShareMessage.SetLink(shareID)

	share := entities.ShareAttribute{
		ParentShareID: c.ParentShareID,
		ShareID:       shareID,
		Sender:        c.Sender,
		Message:       c.ShareMessage,
		SendTime:      sendTime,
	}

	return entity.ShareAccountAttribute{
		ShareID:        share.ShareID,
		Total:          c.AccountNum,
		CanShareNumber: c.AccountNum,
		Share:          &share,
	}
}

// GetShareResponse ...
type GetShareResponse struct {
}

// CancelShareRequest ...
type CancelShareRequest struct {
	ShareID string
}

// ReceiveAccountRequest ...
type ReceiveAccountRequest struct {
	ShareID string
}

type createAccountResponse struct {
	Account  string
	Password string
}

// ReceiveAccountResponse ... 
type ReceiveAccountResponse struct {
	Account  string
	Password string
}

// SetAcct ...
func (ar *ReceiveAccountResponse) SetAcct(car createAccountResponse) {
	ar.Account = car.Account
	ar.Password = car.Password
}
