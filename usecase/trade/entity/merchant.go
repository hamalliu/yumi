package entity

//WxPayMerchant ...
type WxPayMerchant struct {
	//卖家key
	SellerKey string `db:"seller_key" json:"seller_key"`
	//开发平台唯一id
	AppID string `db:"app_id" json:"app_id"`
	//微信商户id
	MchID string `db:"mch_id" json:"mch_id"`
	//私钥
	PrivateKey string `db:"private_key" json:"private_key"`
	//AppSecret（Secret）是APPID对应的接口密码，用于获取接口调用凭证access_token时使用
	Secret string `db:"secret" json:"secret"`
}

//AliPayMerchant ...
type AliPayMerchant struct {
	//卖家key
	SellerKey string `db:"seller_key" json:"seller_key"`
	//商户id
	AppID string `db:"app_id" json:"app_id"`
	//私钥
	PrivateKey string `db:"private_key" json:"private_key"`
	//公钥
	PublicKey string `db:"public_key" json:"public_key"`
}
