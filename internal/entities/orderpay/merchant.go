package orderpay

type AliApp struct {
	AppId      string `db:"app_id"`      //开发应用id
	SellerId   string `db:"seller_id"`   //收款支付宝账号对应的支付宝唯一用户号
	PrivateKey string `db:"private_key"` //私钥
	PublicKey  string `db:"public_key"`  //公钥
}

type WxApp struct {
	AppId      string `db:"app_id"`      //公众账号id
	MchId      string `db:"mch_id"`      //商户号
	PrivateKey string `db:"private_key"` //私钥
}

func GetAliApp(appId string) (AliApp, error) {
	return AliApp{}, nil
}

func GetWxApp(appId string) (WxApp, error) {
	return WxApp{}, nil
}
