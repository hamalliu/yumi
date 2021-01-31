package thirdpf

import "yumi/usecase/trade/entity"

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
	entity.RegisterTrade(TradeWayAliPayPage, NewAliPayPage())

	entity.RegisterTrade(TradeWayWxPayAPP, NewWxPayApp())
	entity.RegisterTrade(TradeWayWxPayJSAPI, NewWxPayJsapi())
	entity.RegisterTrade(TradeWayWxPayMWEB, NewWxPayMweb(WxPayMwebConfig{}))
	entity.RegisterTrade(TradeWayWxPayNATIVE1, NewWxPayNative1())
	entity.RegisterTrade(TradeWayWxPayNATIVE2, NewWxPayNative2())
}
