package conf

import "yumi/usecase/trade/thirdpf/wxpay"

// Trade ...
type Trade struct {
	H5Info H5Info
}

//H5Info ...
type H5Info struct {
	Type    string //场景类型
	WapURL  string //WAP网站URL地址
	WapName string //WAP网站名
}

// MwebConfig ...
func (t Trade) MwebConfig() wxpay.MwebConfig {
	return wxpay.MwebConfig{
		H5Info: wxpay.H5Info{
			Type: t.H5Info.Type,
			WapURL: t.H5Info.WapURL,
			WapName: t.H5Info.WapName,
		},
	}
}
