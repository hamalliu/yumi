package trade

import (
	"yumi/pkg/types"
	"yumi/usecase/trade/entity"
)

// CreateOrderPayRequest ...
type CreateOrderPayRequest struct {
	BuyerAccountGUID string
	SellerKey        string
	TotalFee         int
	Body             string
	Detail           string
	TimeoutExpress   types.Timestamp
}

// Attribute ...
func (m *CreateOrderPayRequest) Attribute(attr *entity.OrderPayAttribute) {
	attr.BuyerAccountGUID = m.BuyerAccountGUID
	attr.SellerKey = m.SellerKey
	attr.TotalFee = m.TotalFee
	attr.Body = m.Body
	attr.Detail = m.Detail
	attr.TimeoutExpress = m.TimeoutExpress
	attr.Status = entity.Submitted

	return
}

// CreateOrderPayResponse ...
type CreateOrderPayResponse struct {
	Code string
}

func (m *CreateOrderPayResponse) set(op entity.OrderPayAttribute) {
	m.Code = op.Code
}

// PayRequest ...
type PayRequest struct {
	Code      string
	TradeWay  string
	ClientIP  string
	NotifyURL string
}

// PayResponse ...
type PayResponse struct {
	Code       string
	OutTradeNo string
	Data       string
}

func (m *PayResponse) set(op entity.OrderPayAttribute) {
	m.Code = op.Code
	m.OutTradeNo = op.OutTradeNo
}
