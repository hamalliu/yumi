package wx_nativepay

type Merchant struct {
	AppId      string `require:"true"` //公众账号id
	MchId      string `require:"true"` //商户号
	PrivateKey string `require:"true"` //私钥
}

func NewMerchant(appId, mchId, privateKey string) Merchant {
	return Merchant{
		AppId:      appId,
		MchId:      mchId,
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
