package platform

import "yumi/usecase/trade"

const (
	//TradeWayAliPayPage ...
	TradeWayAliPayPage = trade.Way("alipay_page")

	//TradeWayWxPayAPP ...
	TradeWayWxPayAPP = trade.Way("wxpay_app")
	//TradeWayWxPayJSAPI ...
	TradeWayWxPayJSAPI = trade.Way("wxpay_jsapi")
	//TradeWayWxPayMWEB ...
	TradeWayWxPayMWEB = trade.Way("wxpay_mweb")
	//TradeWayWxPayNATIVE1 ...
	TradeWayWxPayNATIVE1 = trade.Way("wxpay_native1")
	//TradeWayWxPayNATIVE2 ...
	TradeWayWxPayNATIVE2 = trade.Way("wxpay_native2")
)

// Init trade platform
func Init() {
	trade.RegisterTrade(TradeWayAliPayPage, NewAliPayPage())

	trade.RegisterTrade(TradeWayWxPayAPP, NewWxPayApp())
	trade.RegisterTrade(TradeWayWxPayJSAPI, NewWxPayJsapi())
	trade.RegisterTrade(TradeWayWxPayMWEB, NewWxPayMweb(WxPayMwebConfig{}))
	trade.RegisterTrade(TradeWayWxPayNATIVE1, NewWxPayNative1())
	trade.RegisterTrade(TradeWayWxPayNATIVE2, NewWxPayNative2())
}
