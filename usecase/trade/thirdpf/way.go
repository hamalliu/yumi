package thirdpf

import (
	"yumi/usecase/trade/entity"
	"yumi/usecase/trade/thirdpf/alipay"
	"yumi/usecase/trade/thirdpf/wxpay"
)

const (
	//TradeWayAliPayPage ...
	TradeWayAliPayPage = entity.Way("alipay_page")

	//TradeWayWxPayAPP ...
	TradeWayWxPayAPP = entity.Way("wxpay_app")
	//TradeWayWxPayJSAPI ...
	TradeWayWxPayJSAPI = entity.Way("wxpay_jsapi")
	//TradeWayWxPayMWEB ...
	TradeWayWxPayMWEB = entity.Way("wxpay_mweb")
	//TradeWayWxPayNATIVE1 ...
	TradeWayWxPayNATIVE1 = entity.Way("wxpay_native1")
	//TradeWayWxPayNATIVE2 ...
	TradeWayWxPayNATIVE2 = entity.Way("wxpay_native2")
)

// Init trade platform
func Init() {
	entity.RegisterTrade(TradeWayAliPayPage, alipay.NewPage())

	entity.RegisterTrade(TradeWayWxPayAPP, wxpay.NewApp())
	entity.RegisterTrade(TradeWayWxPayJSAPI, wxpay.NewJsapi())
	entity.RegisterTrade(TradeWayWxPayMWEB, wxpay.NewMweb(wxpay.MwebConfig{}))
	entity.RegisterTrade(TradeWayWxPayNATIVE1, wxpay.NewNative1())
	entity.RegisterTrade(TradeWayWxPayNATIVE2, wxpay.NewNative2())
}
