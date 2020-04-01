package tradeplatform

import (
	"fmt"

	"yumi/external/pay/wx_nativepay"
	"yumi/internal/entities/orderpay"
	"yumi/internal/humble/db"
	"yumi/utils/internal_error"
)

type WxNative1 string

func GetWxNative1() WxNative1 {
	return ""
}

func (wxpay WxNative1) Pay(e *orderpay.Entity) (orderpay.TradePay, error) {
	ret := orderpay.TradePay{}
	//获取收款商户信息
	mch, err := db.GetWxPayMerchantBySellerKey(e.SellerKey)
	if err != nil {
		return ret, err
	}
	wxMch := wx_nativepay.Merchant{
		AppId:      mch.AppId,
		MchId:      mch.MchId,
		PrivateKey: mch.PrivateKey,
	}

	if bizUrl, err := wx_nativepay.GetDefault().BizPayShortUrl(wxMch, e.Code); err != nil {
		return ret, internal_error.With(err)
	} else {
		ret.AppId = wxMch.AppId
		ret.MchId = wxMch.AppId
		ret.Data = bizUrl
		return ret, nil
	}
}

func (wxpay WxNative1) QueryPayStatus(e *orderpay.Entity) (orderpay.TradePayQuery, error) {
	ret := orderpay.TradePayQuery{}
	//获取收款商户信息
	mch, err := db.GetWxPayMerchantBySellerKey(e.SellerKey)
	if err != nil {
		return ret, err
	}
	wxMch := wx_nativepay.Merchant{
		AppId:      mch.AppId,
		MchId:      mch.MchId,
		PrivateKey: mch.PrivateKey,
	}

	if resp, err := wx_nativepay.GetDefault().OrderQuery(wxMch, e.TransactionId, e.OutTradeNo); err != nil {
		return ret, internal_error.With(err)
	} else {
		if e.OutTradeNo != resp.OutTradeNo {
			if resp.OutTradeNo != e.OutTradeNo {
				return ret, fmt.Errorf("订单号不一致")
			}
			if resp.TotalFee != e.TotalFee {
				return ret, internal_error.Critical(fmt.Errorf("订单金额不一致"))
			}
		}
		ret.TransactionId = resp.TransactionId
		switch resp.TradeState {
		case wx_nativepay.TradeStateSuccess:
			ret.TradeStatus = orderpay.TradeStatusSuccess
		case wx_nativepay.TradeStateNotpay, wx_nativepay.TradeStateUserPaying, wx_nativepay.TradeStatePayError,
			wx_nativepay.TradeStateRefund, wx_nativepay.TradeStateRevoked:
			ret.TradeStatus = orderpay.TradeStatusNotPay
		case wx_nativepay.TradeStateClosed:
			ret.TradeStatus = orderpay.TradeStatusClosed
		default:
			err := fmt.Errorf("微信支付状态发生变动，请管理员及时更改")
			return ret, internal_error.Critical(err)
		}
	}

	return ret, nil
}

func (wxpay WxNative1) TradeClose(e *orderpay.Entity) error {
	//获取收款商户信息
	mch, err := db.GetWxPayMerchantBySellerKey(e.SellerKey)
	if err != nil {
		return err
	}
	wxMch := wx_nativepay.Merchant{
		AppId:      mch.AppId,
		MchId:      mch.MchId,
		PrivateKey: mch.PrivateKey,
	}

	if err := wx_nativepay.GetDefault().CloseOrder(wxMch, e.OutTradeNo); err != nil {
		return internal_error.With(err)
	}

	return nil
}

func (wxpay WxNative1) Refund(e *orderpay.Entity) error {
	//TODO
	return nil
}

func (wxpay WxNative1) QueryRefundStatus(e *orderpay.Entity) {
	//TODO
}

//...
//可以添加对账服务，获取支付评价

//func (wxpay WxNative1) PrepayNotify(mch wx_nativepay.Merchant, req wx_nativepay.ReqPrepayNotify) wx_nativepay.RespPrepayNotify {
//	resp := wx_nativepay.RespPrepayNotify{
//		ReturnCode: "SUCCESS",
//		ReturnMsg:  "",
//		AppId:      mch.AppId,
//		MchId:      mch.MchId,
//		NonceStr:   pay.GetNonceStr(),
//	}
//
//	order := orderpay.OrderPay{}
//	if err := order.Load(req.ProductId); err != nil {
//		resp.ResultCode = "FAIL"
//		resp.ErrCodeDes = err.Error()
//		resp.Sign = wx_nativepay.Buildsign(resp, mch.PrivateKey)
//		return resp
//	}
//	defer order.Release()
//
//	if err := wx_nativepay.CheckPrePayNotify(mch, req); err != nil {
//		resp.ResultCode = "FAIL"
//		resp.ErrCodeDes = err.Error()
//		resp.Sign = wx_nativepay.Buildsign(resp, mch.PrivateKey)
//		return resp
//	}
//
//	//下单
//	wxorder := wx_nativepay.UnifiedOrder{
//		Body:       order.Body,
//		Detail:     order.Detail,
//		Attach:     "",
//		OutTradeNo: order.OutTradeNo,
//		TotalFee:   order.TotalFee,
//		NotifyUrl:  order.NotifyUrl,
//		ProductId:  order.Code,
//	}
//	if prepayId, _, err := wx_nativepay.GetDefault().UnifiedOrder(mch, wxorder); err != nil {
//		resp.ResultCode = "FAIL"
//		resp.ErrCodeDes = err.Error()
//		resp.Sign = wx_nativepay.Buildsign(resp, mch.PrivateKey)
//		return resp
//	} else {
//		resp.PrepayId = prepayId
//		resp.ResultCode = "SUCCESS"
//		resp.Sign = wx_nativepay.Buildsign(resp, mch.PrivateKey)
//		return resp
//	}
//}
//
//func (wxpay WxNative1) RefundNotify() {
//	//TODO
//}
