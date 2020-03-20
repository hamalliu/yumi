package pay

import (
	"fmt"
	"time"

	"yumi/external/pay"
	"yumi/external/pay/wx_nativepay"
	"yumi/internal/entity"
	"yumi/utils/internal_error"
)

type WxNative1 string

func GetWxNative1() WxNative1 {
	return ""
}

func (wxpay WxNative1) Pay(mch wx_nativepay.Merchant, code string) (string, error) {
	order := entity.OrderPay{}
	order.Load(code)

	//订单是否过期
	if order.PayExpire < time.Now().Format(TimeFormat) {
		return "", fmt.Errorf("订单已过期不能发起支付")
	}

	//支付订单状态必须为已支付
	if order.Status != entity.OrderPayStatusSubmitted {
		return "", fmt.Errorf("不能发起重复支付订单")
	}

	if bizUrl, err := wx_nativepay.GetDefault().BizPayShortUrl(mch, code); err != nil {
		return "", internal_error.With(err)
	} else {
		if err := order.SetPayWay(PayWayWxNative1Pay, mch.AppId, mch.MchId); err != nil {
			return "", err
		}
		return bizUrl, nil
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

	order := entity.OrderPay{}
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

func (wxpay WxNative1) PayCompleted(mch wx_nativepay.Merchant, order entity.OrderPay) error {
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

func (wxpay WxNative1) PayProblem(mch wx_nativepay.Merchant, order entity.OrderPay) error {
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

func (wxpay WxNative1) CheckStartPay(mch wx_nativepay.Merchant, order entity.OrderPay) error {
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

func (wxpay WxNative1) CloseOrder(mch wx_nativepay.Merchant, order entity.OrderPay) error {
	if order.Status != entity.OrderPayStatusWaitPay {
		return fmt.Errorf("只能关闭待支付的账单")
	}

	if err := wx_nativepay.GetDefault().CloseOrder(mch, order.OutTradeNo); err != nil {
		return internal_error.With(err)
	}

	if err := order.PayClose(); err != nil {
		return err
	}

	return nil
}

func (wxpay WxNative1) PayNotify(mch wx_nativepay.Merchant, req wx_nativepay.ReqPayNotify) error {
	order := entity.OrderPay{}
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

	if order.Status == entity.OrderPayStatusPaid {
		return nil
	}
	if order.Status == entity.OrderPayStatusWaitPay {
		if payTime, err := toTimeFormat(req.TimeEnd); err != nil {
			return err
		} else {
			if err := order.PaySuccess(payTime); err != nil {
				return err
			}
		}
		return nil
	}
	if order.Status == entity.OrderPayStatusSubmitted ||
		order.Status == entity.OrderPayStatusCancelled {
		//支付宝没有处理好，关闭订单和支付执行顺序。
		//TODO 这种情况可能引发重复支付，应该记录紧急日志，并通知管理员处理。
	}

	return nil
}

func (wxpay WxNative1) Refund() {
	//TODO
}

func (wxpay WxNative1) QueryRefundStatus() {
	//TODO
}

func (wxpay WxNative1) RefundNotify() {
	//TODO
}

func toTimeFormat(timeStr string) (string, error) {
	t, err := time.Parse("20060102150405", timeStr)
	if err != nil {
		return "", err
	}
	return t.Format(TimeFormat), nil
}

//...
//可以添加对账服务，获取支付评价
