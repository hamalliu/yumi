package wxpay

import (
	"yumi/pkg/externalapi/txapi/wxpay"
	"yumi/usecase/trade/entity"
)

//Native1 ...
type Native1 struct {
	Internal
}

//NewNative1 ...
func NewNative1() Native1 {
	return Native1{}
}

//Pay 发起支付
func (wxn1 Native1) Pay(op entity.OrderPayAttribute) (entity.ReturnPay, error) {
	ret := entity.ReturnPay{}

	//获取收款商户信息
	wxMch, err := wxn1.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	bizURL, err := wxpay.GetDefault().BizPayShortURL(wxMch, op.OutTradeNo)
	if err != nil {
		return ret, err
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
//	order := entity.OrderPay{}
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
