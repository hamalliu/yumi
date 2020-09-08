package tradeplatform

import (
	"yumi/internal/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/external/pay/wxpay"
)

//WxPayNATIVE1 ...
const WxPayNATIVE1 = trade.Way("wxpay_native1")

//WxNative1 ...
type WxNative1 struct {
	InternalWxPay
}

//GetWxNative1 ...
func GetWxNative1() WxNative1 {
	return WxNative1{}
}

//Pay 发起支付
func (wxn1 WxNative1) Pay(op trade.OrderPay) (trade.ReturnPay, error) {
	ret := trade.ReturnPay{}

	//获取收款商户信息
	wxMch, err := wxn1.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	bizURL, err := wxpay.GetDefault().BizPayShortURL(wxMch, op.OutTradeNo)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}

	ret.AppID = wxMch.AppID
	ret.MchID = wxMch.MchID
	ret.Data = bizURL
	return ret, nil
}

//func (wxpay WxNative1) PrepayNotify(mch wxpay.Merchant, req wxpay.ReqPrepayNotify) wxpay.RespPrepayNotify {
//	resp := wxpay.RespPrepayNotify{
//		ReturnCode: "SUCCESS",
//		ReturnMsg:  "",
//		AppId:      mch.AppId,
//		MchId:      mch.MchId,
//		NonceStr:   pay.GetNonceStr(),
//	}
//
//	order := trade.OrderPay{}
//	if err := order.Load(req.ProductId); err != nil {
//		resp.ResultCode = "FAIL"
//		resp.ErrCodeDes = "服务器错误"
//		resp.Sign = wxpay.Buildsign(resp, mch.PrivateKey)
//		return resp
//	}
//	defer order.Release()
//
//	if err := wxpay.CheckPrePayNotify(mch, req); err != nil {
//		logs.Error(err)
//		resp.ResultCode = "FAIL"
//		resp.ErrCodeDes = err.Error()
//		resp.Sign = wxpay.Buildsign(resp, mch.PrivateKey)
//		return resp
//	}
//
//	//下单
//	wxorder := wxpay.UnifiedOrder{
//		Body:       order.Body,
//		Detail:     order.Detail,
//		Attach:     "",
//		OutTradeNo: order.OutTradeNo,
//		TotalFee:   order.TotalFee,
//		NotifyUrl:  order.NotifyUrl,
//		ProductId:  order.Code,
//	}
//	if prepayId, _, err := wxpay.GetDefault().UnifiedOrder(mch, wxorder); err != nil {
//		resp.ResultCode = "FAIL"
//		resp.ErrCodeDes = err.Error()
//		resp.Sign = wxpay.Buildsign(resp, mch.PrivateKey)
//		return resp
//	} else {
//		resp.PrepayId = prepayId
//		resp.ResultCode = "SUCCESS"
//		resp.Sign = wxpay.Buildsign(resp, mch.PrivateKey)
//		return resp
//	}
//}

//func (wxpay WxNative1) RefundNotify() {
//	//TODO
//}

//...
//可以添加对账服务，获取支付评价
