package trade

import (
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/pkg/stores/mgoc"
	"yumi/usecase/trade/data"
	"yumi/usecase/trade/service"
	"yumi/usecase/trade/thirdpf"
	"yumi/usecase/trade/thirdpf/alipay"
	"yumi/usecase/trade/thirdpf/wxpay"
)

// Usecase 新建用例
func Usecase(mongoC *mgoc.Client, mysqlC *mysqlx.Client, wxMwebConf wxpay.MwebConfig) (*service.Service, error) {
	data := data.New(mysqlC)
	trades := thirdpf.New()
	// 支付宝
	trades.AddTrade(thirdpf.TradeWayAliPayPage, alipay.NewPage(data))
	// 微信
	trades.AddTrade(thirdpf.TradeWayWxPayAPP, wxpay.NewApp(data))
	trades.AddTrade(thirdpf.TradeWayWxPayJSAPI, wxpay.NewJsapi(data))
	trades.AddTrade(thirdpf.TradeWayWxPayMWEB, wxpay.NewMweb(wxMwebConf, data))
	trades.AddTrade(thirdpf.TradeWayWxPayNATIVE1, wxpay.NewNative1(data))
	trades.AddTrade(thirdpf.TradeWayWxPayNATIVE2, wxpay.NewNative2(data))

	return service.New(data, trades)
}
