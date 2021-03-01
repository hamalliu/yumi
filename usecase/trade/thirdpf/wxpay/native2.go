package wxpay

import (
	"yumi/pkg/externalapi/txapi/wxpay"
	"yumi/usecase/trade/entity"
)

//Native2 ...
type Native2 struct {
	Internal
}

//NewNative2 ...
func NewNative2() Native2 {
	return Native2{}
}

//Pay 发起支付
func (wxn1 Native2) Pay(op entity.OrderPayAttribute) (entity.ReturnPay, error) {
	ret := entity.ReturnPay{}
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
		PayExpire:      op.PayExpire.Time(),
		SpbillCreateIP: op.SpbillCreateIP,
	}

	retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeNative, wxMch, wxorder)
	if err != nil {
		return ret, err
	}

	ret.AppID = wxMch.AppID
	ret.MchID = wxMch.MchID
	ret.Data = retuo.CodeURL
	return ret, nil
}
