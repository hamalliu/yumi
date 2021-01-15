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
	entities.ShareMessage
}

// ShareAccountAttribute ...
func (c *CreateShareRequest) ShareAccountAttribute() entity.ShareAccountAttribute {
	share := entities.ShareAttribute{
		ParentShareID: c.ParentShareID,
		ShareID:       entities.NewShareID(),
		Sender:        c.Sender,
		Message:       c.ShareMessage,
		SendTime:      types.NowTimestamp(),
	}

	return entity.ShareAccountAttribute{
		ShareID: share.ShareID,
		Total:   c.AccountNum,
		CanShareNumber: c.AccountNum,
		Share:   &share,
	}
}

// GetShareResponse ...
type GetShareResponse struct {
}
