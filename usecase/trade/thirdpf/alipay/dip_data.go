package alipay

import "yumi/usecase/trade/entity"

var _data Data

// Data ...
type Data interface {
	GetAliPayMerchant(sellerKey string) (entity.AliPayMerchant, error)
}
