package ali_pagepay

type Merchant struct {
	AppId      string `require:"true"` //商户id
	PrivateKey string `require:"true"` //私钥
	PublicKey  string `require:"true"` //公钥
}

func NewMerchant(appId, privateKey string) Merchant {
	return Merchant{
		AppId:      appId,
		PrivateKey: privateKey,
	}
}

//生成商户订单号
func (mch Merchant) BuildOutTradeNo() string {
	//TODO
}

//生成商户退款订单号
func (mch Merchant) BuildOutRefundNo() string {
	//TODO
}
