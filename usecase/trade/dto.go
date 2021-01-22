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
	code := getCode(OrderPayCode)

	attr.Code = code
	attr.SellerKey = m.SellerKey
	attr.TotalFee = m.TotalFee
	attr.Body = m.Body
	attr.Detail =  m.Detail
	attr.SubmitTime = types.NowTimestamp()
	attr.TimeoutExpress = m.TimeoutExpress
	attr.Status = entity.Submitted

	return 
}

// CreateOrderPayResponse ...
type CreateOrderPayResponse struct {
	Code string
}

// PayRequest ...
type PayRequest struct {
	Code      string
	TradeWay  string
	ClientIP  string
	NotifyURL string
	PayExpire types.Timestamp
}

// Attribute ...
func (m *PayRequest) Attribute(attr *entity.OrderPayAttribute) {
	attr.TradeWay = m.TradeWay
	attr.SpbillCreateIP = m.ClientIP
	attr.NotifyURL = m.NotifyURL
	attr.PayExpire = m.PayExpire
	attr.Status = entity.WaitPay
	attr.OutTradeNo = getOutTradeNo()
}

// PayResponse ...
type PayResponse struct {
	Code       string
	OutTradeNo string
}
