package wxpay

type Merchant struct {
	AppId      string `require:"true"` //公众账号id
	MchId      string `require:"true"` //商户号
	PrivateKey string `require:"true"` //私钥
	AppSecret  string //AppSecret（Secret）是APPID对应的接口密码，用于获取接口调用凭证access_token时使用
}
