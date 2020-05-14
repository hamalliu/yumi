package tradeplatform

import (
	"yumi/internal/pay/entities/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/external/pay/wxpay"
)

const WxPay_NATIVE2 = trade.Way("wxpay_native2")

type WxNative2 struct {
	InternalWxPay
}

func GetWxNative2() WxNative2 {
	return WxNative2{}
}

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
		NotifyUrl:      op.NotifyUrl,
		ProductId:      op.OutTradeNo,
		PayExpire:      op.PayExpire,
		SpbillCreateIp: op.SpbillCreateIp,
	}
	if retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeNative, wxMch, wxorder); err != nil {
		return ret, ecode.ServerErr(err)
	} else {
		ret.AppId = wxMch.AppId
		ret.MchId = wxMch.MchId
		ret.Data = retuo.CodeUrl
		return ret, nil
	}
}
