package wxpay

import "yumi/usecase/trade/entity"

var _data Data

// FilterMerchantIDs ...
type FilterMerchantIDs struct {
	MchID     string
	SellerKey string
}

// Data ...
type Data interface {
	GetWxPayMerchant(ids FilterMerchantIDs) (entity.WxPayMerchant, error)
}
