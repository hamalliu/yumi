package tradeplatform

import (
	"fmt"
	"time"

	"yumi/external/pay"
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

func (wxpay WxNative1) PayNotify(mch wx_nativepay.Merchant, req wx_nativepay.ReqPayNotify) error {
	order := orderpay.OrderPay{}
	if err := order.LoadByOutTradeNo(req.AppId, req.OutTradeNo); err != nil {
		return err
	}
	defer order.Release()

	if err := wx_nativepay.CheckPayNotify(mch, order.TotalFee, order.OutTradeNo, req); err != nil {
		return internal_error.With(err)
	}

	if err := order.SetTransactionId(req.TransactionId, ""); err != nil {
		return err
	}

	if order.Status == orderpay.Paid {
		return nil
	}
	if order.Status == orderpay.WaitPay {
		if payTime, err := toTimeFormat(req.TimeEnd); err != nil {
			return err
		} else {
			if err := order.PaySuccess(payTime); err != nil {
				return err
			}
		}
		return nil
	}
	if order.Status == orderpay.Submitted ||
		order.Status == orderpay.Cancelled {
		//支付宝没有处理好，关闭订单和支付执行顺序。
		//TODO 这种情况可能引发重复支付，应该记录紧急日志，并通知管理员处理。
	}

	return nil
}

func toTimeFormat(timeStr string) (time.Time, error) {
	t, err := time.Parse("20060102150405", timeStr)
	if err != nil {
		return t, err
	}
	return t, nil
}

//...
//可以添加对账服务，获取支付评价
func (wxpay WxNative1) PayCompleted(mch wx_nativepay.Merchant, order orderpay.OrderPay) error {
	if ret, err := wx_nativepay.GetDefault().OrderQuery(mch, order.TransactionId, order.OutTradeNo); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TransactionId != order.TransactionId ||
			ret.OutTradeNo != order.OutTradeNo ||
			ret.TotalFee != order.TotalFee {
			return fmt.Errorf("订单信息不一致")
		}

		if err := order.SetTransactionId(ret.TransactionId, ""); err != nil {
			return err
		}
		if ret.TradeState == wx_nativepay.TradeStateSuccess {
			if payTime, err := toTimeFormat(ret.TimeEnd); err != nil {
				return err
			} else {
				if err := order.PaySuccess(payTime); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func (wxpay WxNative1) PayProblem(mch wx_nativepay.Merchant, order orderpay.OrderPay) error {
	if ret, err := wx_nativepay.GetDefault().OrderQuery(mch, order.TransactionId, order.OutTradeNo); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TransactionId != order.TransactionId ||
			ret.OutTradeNo != order.OutTradeNo ||
			ret.TotalFee != order.TotalFee {
			return fmt.Errorf("订单信息不一致")
		}

		if err := order.SetTransactionId(ret.TransactionId, ""); err != nil {
			return err
		}

		//if ret.TradeState == wx_nativepay.TradeStateSuccess {
		//	if payTime, err := toTimeFormat(ret.TimeEnd); err != nil {
		//		return err
		//	} else {
		//		if err := order.PaySuccess(payTime); err != nil {
		//			return err
		//		}
		//	}
		//} else if ret.TradeState == wx_nativepay.TradeStateNotpay {
		//	//如果用户未支付则关闭该订单
		//	if err := wxpay.CloseOrder(mch, order); err != nil {
		//		return err
		//	}
		//}
		return nil
	}
}

func (wxpay WxNative1) CheckStartPay(mch wx_nativepay.Merchant, order orderpay.OrderPay) error {
	if ret, err := wx_nativepay.GetDefault().OrderQuery(mch, order.TransactionId, order.OutTradeNo); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TransactionId != order.TransactionId ||
			ret.OutTradeNo != order.OutTradeNo ||
			ret.TotalFee != order.TotalFee {
			return fmt.Errorf("订单信息不一致")
		}

		if err := order.SetTransactionId(ret.TransactionId, ""); err != nil {
			return err
		}

		if ret.TradeState == wx_nativepay.TradeStateSuccess {
			if payTime, err := toTimeFormat(ret.TimeEnd); err != nil {
				return err
			} else {
				if err := order.PaySuccess(payTime); err != nil {
					return err
				}
			}
		} else if ret.TradeState == wx_nativepay.TradeStateNotpay {
			//如果用户未支付则关闭该订单
			if err := wxpay.CloseOrder(mch, order); err != nil {
				return err
			}
		}
		return nil
	}
}

func (wxpay WxNative1) PrepayNotify(mch wx_nativepay.Merchant, req wx_nativepay.ReqPrepayNotify) wx_nativepay.RespPrepayNotify {
	resp := wx_nativepay.RespPrepayNotify{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "",
		AppId:      mch.AppId,
		MchId:      mch.MchId,
		NonceStr:   pay.GetNonceStr(),
	}

	order := orderpay.OrderPay{}
	if err := order.Load(req.ProductId); err != nil {
		resp.ResultCode = "FAIL"
		resp.ErrCodeDes = err.Error()
		resp.Sign = wx_nativepay.Buildsign(resp, mch.PrivateKey)
		return resp
	}
	defer order.Release()

	if err := wx_nativepay.CheckPrePayNotify(mch, req); err != nil {
		resp.ResultCode = "FAIL"
		resp.ErrCodeDes = err.Error()
		resp.Sign = wx_nativepay.Buildsign(resp, mch.PrivateKey)
		return resp
	}

	//下单
	wxorder := wx_nativepay.UnifiedOrder{
		Body:       order.Body,
		Detail:     order.Detail,
		Attach:     "",
		OutTradeNo: order.OutTradeNo,
		TotalFee:   order.TotalFee,
		NotifyUrl:  order.NotifyUrl,
		ProductId:  order.Code,
	}
	if prepayId, _, err := wx_nativepay.GetDefault().UnifiedOrder(mch, wxorder); err != nil {
		resp.ResultCode = "FAIL"
		resp.ErrCodeDes = err.Error()
		resp.Sign = wx_nativepay.Buildsign(resp, mch.PrivateKey)
		return resp
	} else {
		resp.PrepayId = prepayId
		resp.ResultCode = "SUCCESS"
		resp.Sign = wx_nativepay.Buildsign(resp, mch.PrivateKey)
		return resp
	}
}

func (wxpay WxNative1) RefundNotify() {
	//TODO
}
