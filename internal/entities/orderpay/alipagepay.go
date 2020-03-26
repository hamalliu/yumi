package orderpay

import (
	"fmt"
	"time"

	"yumi/external/pay/ali_pagepay"
	"yumi/utils/internal_error"
)

type AliPagePay string

func GetAliPagePay() AliPagePay {
	return ""
}

func (alipay AliPagePay) Pay(mch ali_pagepay.Merchant, order *OrderPay) ([]byte, error) {
	//订单是否过期
	if order.PayExpire < time.Now().Format(TimeFormat) {
		return nil, fmt.Errorf("订单已过期不能发起支付")
	}

	//支付订单状态必须为已支付
	if order.Status != OrderPayStatusSubmitted {
		return nil, fmt.Errorf("不能发起重复支付订单")
	}

	//下单
	pagePay := ali_pagepay.PagePay{
		OutTradeNo:  order.OutTradeNo,
		ProductCode: order.Code,
		TotalAmount: toPrice(order.TotalFee),
		Subject:     order.Body,
		Body:        order.Detail,
		GoodsType:   "0",
		NotifyUrl:   order.NotifyUrl,
	}
	if ret, err := ali_pagepay.GetDefault().UnifiedOrder(mch, pagePay); err != nil {
		return nil, internal_error.With(err)
	} else {
		if err := order.SetPayWay(PayWayAliPagePay, mch.AppId, ret.SellerId); err != nil {
			return nil, internal_error.With(err)
		}

		return ret.PagePayHtml, nil
	}
}

func (alipay AliPagePay) QueryPayStatus(mch ali_pagepay.Merchant, order *OrderPay) (TradeStatus, error) {
	query := ali_pagepay.TradeQuery{
		TradeNo:    order.TransactionId,
		OutTradeNo: order.OutTradeNo,
	}

	if ret, err := ali_pagepay.GetDefault().TradeQuery(mch, query); err != nil {
		return "", internal_error.With(err)
	} else {
		if ret.TradeNo != order.TransactionId ||
			ret.OutTradeNo != order.OutTradeNo ||
			ret.TotalAmount != toPrice(order.TotalFee) {
			return "", fmt.Errorf("订单信息不一致")
		}
		if err := order.SetTransactionId(ret.TradeNo, ret.BuyerlogonId); err != nil {
			return "", internal_error.With(err)
		}

		switch ret.TradeStatus {
		case ali_pagepay.TradeStatusSuccess:
			return TradeStatusSuccess, nil
		case ali_pagepay.TradeStatusWaitBuyerPay:
			return TradeStatusNotPay, nil
		case ali_pagepay.TradeStatusCloseed:
			return TradeStatusClosed, nil
		case ali_pagepay.TradeStatusFinished:
			return TradeStatusFinished, nil
		default:
			err := fmt.Errorf("支付宝状态发生变动，请管理员及时更改")
			return "", internal_error.Critical(err)
		}
	}
}

func (alipay AliPagePay) TradeClose(mch ali_pagepay.Merchant, order *OrderPay) error {
	close := ali_pagepay.TradeClose{
		OutTradeNo: order.OutTradeNo,
		TradeNo:    order.TransactionId,
		OperatorId: "sys",
	}
	if order.Status != OrderPayStatusWaitPay {
		return fmt.Errorf("只能关闭待支付的账单")
	}

	if ret, err := ali_pagepay.GetDefault().TradeClose(mch, close); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TradeNo != order.TransactionId ||
			ret.OutTradeNo != order.OutTradeNo {
			return fmt.Errorf("订单信息不一致")
		}
	}

	return nil
}

func (alipay AliPagePay) PayNotify(mch ali_pagepay.Merchant, rawQuery string) (TradeStatus, error) {
	reqNotify, err := ali_pagepay.ParseQuery(rawQuery)
	if err != nil {
		return "", internal_error.With(err)
	}

	order := OrderPay{}
	if err := order.LoadByOutTradeNo(reqNotify.AppId, reqNotify.OutTradeNo); err != nil {
		return "", internal_error.With(err)
	}
	defer order.Release()

	if err := ali_pagepay.CheckPayNotify(mch, order.OutTradeNo, toPrice(order.TotalFee), order.MchId, reqNotify); err != nil {
		return "", err
	}

	if err := order.SetTransactionId(reqNotify.TradeNo, reqNotify.BuyerId); err != nil {
		return "", internal_error.With(err)
	}

	switch reqNotify.TradeStatus {
	case ali_pagepay.TradeStatusSuccess:
		return TradeStatusSuccess, nil
	case ali_pagepay.TradeStatusWaitBuyerPay:
		return TradeStatusNotPay, nil
	case ali_pagepay.TradeStatusCloseed:
		return TradeStatusClosed, nil
	case ali_pagepay.TradeStatusFinished:
		return TradeStatusFinished, nil
	default:
		err := fmt.Errorf("支付宝状态发生变动，请管理员及时更改")
		return "", internal_error.Critical(err)
	}
	//if reqNotify.TradeStatus == ali_pagepay.TradeStatusSuccess {
	//	if order.Status == OrderPayStatusPaid {
	//		return nil
	//	}
	//	if order.Status == OrderPayStatusWaitPay {
	//		if err := order.PaySuccess(reqNotify.GmtPayment); err != nil {
	//			return err
	//		}
	//		return nil
	//	}
	//	if order.Status == OrderPayStatusSubmitted ||
	//		order.Status == OrderPayStatusCancelled {
	//		//支付宝没有处理好，关闭订单和支付执行顺序。
	//		//TODO 这种情况可能引发重复支付，应该记录紧急日志，并通知管理员处理。
	//	}
	//}
	//
	//return fmt.Errorf("交易状态错误")
}

func (alipay AliPagePay) Refund() {
	//TODO
}

func (alipay AliPagePay) QueryRefundStatus() {
	//TODO
}

func (alipay AliPagePay) RefundNotify() {
	//TODO
}

func toPrice(amount int) string {
	return fmt.Sprintf("%d.%02d", amount/100, amount%100)
}

//可以添加对账服务，获取支付评价

//=======================================================================================================
func (alipay AliPagePay) PayCompleted(mch ali_pagepay.Merchant, order *OrderPay) error {

	query := ali_pagepay.TradeQuery{
		TradeNo:    order.TransactionId,
		OutTradeNo: order.OutTradeNo,
	}

	if ret, err := ali_pagepay.GetDefault().TradeQuery(mch, query); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TradeNo != order.TransactionId ||
			ret.OutTradeNo != order.OutTradeNo ||
			ret.TotalAmount != toPrice(order.TotalFee) {
			return fmt.Errorf("订单信息不一致")
		}

		if err := order.SetTransactionId(ret.TradeNo, ret.BuyerlogonId); err != nil {
			return err
		}
		if ret.TradeStatus == ali_pagepay.TradeStatusSuccess {
			if err := order.PaySuccess(ret.SendPayDate); err != nil {
				return err
			}
		}
		return nil
	}
}

func (alipay AliPagePay) PayProblem(mch ali_pagepay.Merchant, order *OrderPay) error {
	query := ali_pagepay.TradeQuery{
		TradeNo:    order.TransactionId,
		OutTradeNo: order.OutTradeNo,
	}

	if ret, err := ali_pagepay.GetDefault().TradeQuery(mch, query); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TradeNo != order.TransactionId ||
			ret.OutTradeNo != order.OutTradeNo ||
			ret.TotalAmount != toPrice(order.TotalFee) {
			return fmt.Errorf("订单信息不一致")
		}
		if err := order.SetTransactionId(ret.TradeNo, ret.BuyerlogonId); err != nil {
			return err
		}

		if ret.TradeStatus == ali_pagepay.TradeStatusSuccess {
			if err := order.PaySuccess(ret.SendPayDate); err != nil {
				return err
			}
		} else if ret.TradeStatus == ali_pagepay.TradeStatusWaitBuyerPay {
			//如果用户未支付则关闭该订单
			if err := alipay.CloseOrder(mch, order); err != nil {
				return err
			}
		}
		return nil
	}
}

func (alipay AliPagePay) CheckStartPay(mch ali_pagepay.Merchant, order *OrderPay) error {
	query := ali_pagepay.TradeQuery{
		TradeNo:    order.TransactionId,
		OutTradeNo: order.OutTradeNo,
	}

	if ret, err := ali_pagepay.GetDefault().TradeQuery(mch, query); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TradeNo != order.TransactionId ||
			ret.OutTradeNo != order.OutTradeNo ||
			ret.TotalAmount != toPrice(order.TotalFee) {
			return fmt.Errorf("订单信息不一致")
		}
		if err := order.SetTransactionId(ret.TradeNo, ret.BuyerlogonId); err != nil {
			return err
		}

		if ret.TradeStatus == ali_pagepay.TradeStatusSuccess {
			if err := order.PaySuccess(ret.SendPayDate); err != nil {
				return err
			}
			return fmt.Errorf("该订单已支付，不能重复发起")
		} else if ret.TradeStatus == ali_pagepay.TradeStatusWaitBuyerPay {
			//如果用户未支付则关闭该订单
			if err := alipay.CloseOrder(mch, order); err != nil {
				return err
			}
		}
		return nil
	}
}
