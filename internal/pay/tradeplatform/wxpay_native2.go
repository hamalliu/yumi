package tradeplatform

import (
	"yumi/internal/pay/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/external/pay/wxpay"
)

//WxPayNATIVE2 ...
const WxPayNATIVE2 = trade.Way("wxpay_native2")

//WxNative2 ...
type WxNative2 struct {
	InternalWxPay
}

//GetWxNative2 ...
func GetWxNative2() WxNative2 {
	return WxNative2{}
}

//Pay ...
func (wxn1 WxNative2) Pay(op trade.OrderPay) (trade.ReturnPay, error) {
	ret := trade.ReturnPay{}
	//获取收款商户信息
	wxMch, err := wxn1.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	wxorder := wxpay.UnifiedOrder{
		Body:           op.Body,
		Detail:         op.Detail,
		Attach:         op.Code,
		OutTradeNo:     op.OutTradeNo,
		TotalFee:       op.TotalFee,
		NotifyUrl:      op.NotifyURL,
		ProductId:      op.OutTradeNo,
		PayExpire:      op.PayExpire,
		SpbillCreateIp: op.SpbillCreateIP,
	}

	retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeNative, wxMch, wxorder)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}
	
	ret.AppID = wxMch.AppId
	ret.MchID = wxMch.MchId
	ret.Data = retuo.CodeUrl
	return ret, nil
}
