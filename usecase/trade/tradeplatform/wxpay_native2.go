package tradeplatform

import (
	"yumi/usecase/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/external/trade/wxpay"
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

//Pay 发起支付
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
		NotifyURL:      op.NotifyURL,
		ProductID:      op.OutTradeNo,
		PayExpire:      op.PayExpire,
		SpbillCreateIP: op.SpbillCreateIP,
	}

	retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeNative, wxMch, wxorder)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}

	ret.AppID = wxMch.AppID
	ret.MchID = wxMch.MchID
	ret.Data = retuo.CodeURL
	return ret, nil
}
