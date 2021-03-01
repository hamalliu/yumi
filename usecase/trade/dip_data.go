package trade

import (
	"yumi/usecase/trade/entity"
)

var _data Data

// Data ...
type Data interface {
	CreateOrderPay(entity.OrderPayAttribute) error
	UpdateOrderPay(entity.OrderPayAttribute) error
	GetOrderPay(orderID string) (entity.OrderPayAttribute, error)

	CreateOrderRefund(entity.OrderRefundAttribute) error
	UpdateOrderRefund(entity.OrderRefundAttribute) error
	GetOrderRefund(orderID string) (entity.OrderRefundAttribute, error)

	CreateWxPayMerchant(entity.WxPayMerchant) error
	CreateAliPayMerchant(entity.AliPayMerchant) error
}

// InitData ...
func InitData(data Data) {
	_data = data
}

// getData ...
func getData() Data {
	return _data
}
