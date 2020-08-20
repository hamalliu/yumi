package alipay

//Merchant ...
type Merchant struct {
	AppID      string `require:"true"` //商户id
	PrivateKey string `require:"true"` //私钥
	PublicKey  string `require:"true"` //公钥
}
