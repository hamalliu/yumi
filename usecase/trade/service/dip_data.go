package service

import (
	"yumi/usecase/trade/entity"
)

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
