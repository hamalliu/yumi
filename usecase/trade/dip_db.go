package trade

import (
	"yumi/usecase/trade/entity"
)

var _data Data

// Data ...
type Data interface {
	CreateOrderPay(entity.OrderPayAttribute) error
	GetOrderPay(orderID string) (DataOrderPay, error)

	CreateOrderRefund(entity.OrderRefundAttribute) error
	GetOrderRefund(orderID string) (DataOrderRefund, error)
}

// DataOrderPay ...
type DataOrderPay interface {
	Attribute() *entity.OrderPayAttribute

	Update() error
}

// DataOrderRefund ...
type DataOrderRefund interface {
	Attribute() *entity.OrderRefundAttribute

	Update() error
}

// InitData ...
func InitData(data Data) {
	_data = data
}

// GetData ...
func GetData() Data {
	return _data
}
