package threepay

import (
	"fmt"
	"time"

	"yumi/external/pay/ali_pagepay"
	"yumi/internal/entities/orderpay"
	"yumi/utils/internal_error"
)

const aliPagePay = orderpay.TradeWay("alipagepay")

type AliPagePay string

func GetAliPagePay() AliPagePay {
	return ""
}

func (alipay AliPagePay) Pay(mch ali_pagepay.Merchant, e *orderpay.Entity) ([]byte, error) {
	//订单是否过期
	if e.PayExpire.Unix() < time.Now().Unix() {
		return nil, fmt.Errorf("订单已过期不能发起支付")
	}

	//支付订单状态必须为已支付
	if e.Status != orderpay.Submitted {
		return nil, fmt.Errorf("不能发起重复支付订单")
	}

	//下单
	pagePay := ali_pagepay.PagePay{
		OutTradeNo:  e.OutTradeNo,
		ProductCode: e.Code,
		TotalAmount: toPrice(e.TotalFee),
		Subject:     e.Body,
		Body:        e.Detail,
		GoodsType:   "0",
		NotifyUrl:   e.NotifyUrl,
	}
	if ret, err := ali_pagepay.GetDefault().UnifiedOrder(mch, pagePay); err != nil {
		return nil, internal_error.With(err)
	} else {
		if err := e.SetPayWay(aliPagePay, mch.AppId, ret.SellerId); err != nil {
			return nil, internal_error.With(err)
		}

		return ret.PagePayHtml, nil
	}
}

func (alipay AliPagePay) QueryPayStatus(mch ali_pagepay.Merchant, e *orderpay.Entity) (orderpay.TradeStatus, error) {
	query := ali_pagepay.TradeQuery{
		TradeNo:    e.TransactionId,
		OutTradeNo: e.OutTradeNo,
	}

	if ret, err := ali_pagepay.GetDefault().TradeQuery(mch, query); err != nil {
		return "", internal_error.With(err)
	} else {
		if ret.TradeNo != e.TransactionId ||
			ret.OutTradeNo != e.OutTradeNo ||
			ret.TotalAmount != toPrice(e.TotalFee) {
			return "", fmt.Errorf("订单信息不一致")
		}
		if err := e.SetTransactionId(ret.TradeNo, ret.BuyerlogonId); err != nil {
			return "", internal_error.With(err)
		}

		switch ret.TradeStatus {
		case ali_pagepay.TradeStatusSuccess:
			return orderpay.TradeStatusSuccess, nil
		case ali_pagepay.TradeStatusWaitBuyerPay:
			return orderpay.TradeStatusNotPay, nil
		case ali_pagepay.TradeStatusCloseed:
			return orderpay.TradeStatusClosed, nil
		case ali_pagepay.TradeStatusFinished:
			return orderpay.TradeStatusFinished, nil
		default:
			err := fmt.Errorf("支付宝状态发生变动，请管理员及时更改")
			return "", internal_error.Critical(err)
		}
	}
}

func (alipay AliPagePay) TradeClose(mch ali_pagepay.Merchant, e *orderpay.Entity) error {
	close := ali_pagepay.TradeClose{
		OutTradeNo: e.OutTradeNo,
		TradeNo:    e.TransactionId,
		OperatorId: "sys",
	}
	if e.Status != orderpay.WaitPay {
		return fmt.Errorf("只能关闭待支付的账单")
	}

	if ret, err := ali_pagepay.GetDefault().TradeClose(mch, close); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TradeNo != e.TransactionId ||
			ret.OutTradeNo != e.OutTradeNo {
			return fmt.Errorf("订单信息不一致")
		}
	}

	return nil
}

//func (alipay AliPagePay) PayNotify(mch ali_pagepay.Merchant, rawQuery string) (orderpay.TradeStatus, error) {
//	reqNotify, err := ali_pagepay.ParseQuery(rawQuery)
//	if err != nil {
//		return "", internal_error.With(err)
//	}
//
//	e := orderpay.Entity{}
//	if err := e.loadByOutTradeNo(reqNotify.AppId, reqNotify.OutTradeNo); err != nil {
//		return "", internal_error.With(err)
//	}
//	defer func() { _ = e.release() }()
//
//	if err := ali_pagepay.CheckPayNotify(mch, e.OutTradeNo, toPrice(e.TotalFee), e.MchId, reqNotify); err != nil {
//		return "", err
//	}
//
//	if err := e.setTransactionId(reqNotify.TradeNo, reqNotify.BuyerId); err != nil {
//		return "", internal_error.With(err)
//	}
//
//	switch reqNotify.TradeStatus {
//	case ali_pagepay.TradeStatusSuccess:
//		return orderpay.TradeStatusSuccess, nil
//	case ali_pagepay.TradeStatusWaitBuyerPay:
//		return orderpay.TradeStatusNotPay, nil
//	case ali_pagepay.TradeStatusCloseed:
//		return orderpay.TradeStatusClosed, nil
//	case ali_pagepay.TradeStatusFinished:
//		return orderpay.TradeStatusFinished, nil
//	default:
//		err := fmt.Errorf("支付宝状态发生变动，请管理员及时更改")
//		return "", internal_error.Critical(err)
//	}
//	//if reqNotify.TradeStatus == ali_pagepay.TradeStatusSuccess {
//	//	if order.Status == Paid {
//	//		return nil
//	//	}
//	//	if order.Status == WaitPay {
//	//		if err := order.PaySuccess(reqNotify.GmtPayment); err != nil {
//	//			return err
//	//		}
//	//		return nil
//	//	}
//	//	if order.Status == Submitted ||
//	//		order.Status == Cancelled {
//	//		//支付宝没有处理好，关闭订单和支付执行顺序。
//	//		//TODO 这种情况可能引发重复支付，应该记录紧急日志，并通知管理员处理。
//	//	}
//	//}
//	//
//	//return fmt.Errorf("交易状态错误")
//}

func (alipay AliPagePay) Refund() {
	//TODO
}

func (alipay AliPagePay) QueryRefundStatus() {
	//TODO
}

//func (alipay AliPagePay) RefundNotify() {
//	//TODO
//}

func toPrice(amount int) string {
	return fmt.Sprintf("%d.%02d", amount/100, amount%100)
}

//可以添加对账服务，获取支付评价
//=======================================================================================================
func (alipay AliPagePay) PayCompleted(mch ali_pagepay.Merchant, e *orderpay.Entity) error {

	query := ali_pagepay.TradeQuery{
		TradeNo:    e.TransactionId,
		OutTradeNo: e.OutTradeNo,
	}

	if ret, err := ali_pagepay.GetDefault().TradeQuery(mch, query); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TradeNo != e.TransactionId ||
			ret.OutTradeNo != e.OutTradeNo ||
			ret.TotalAmount != toPrice(e.TotalFee) {
			return fmt.Errorf("订单信息不一致")
		}

		if err := e.setTransactionId(ret.TradeNo, ret.BuyerlogonId); err != nil {
			return err
		}
		if ret.TradeStatus == ali_pagepay.TradeStatusSuccess {
			if err := e.paySuccess(ret.SendPayDate); err != nil {
				return err
			}
		}
		return nil
	}
}

func (alipay AliPagePay) PayProblem(mch ali_pagepay.Merchant, e *orderpay.Entity) error {
	query := ali_pagepay.TradeQuery{
		TradeNo:    e.TransactionId,
		OutTradeNo: e.OutTradeNo,
	}

	if ret, err := ali_pagepay.GetDefault().TradeQuery(mch, query); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TradeNo != e.TransactionId ||
			ret.OutTradeNo != e.OutTradeNo ||
			ret.TotalAmount != toPrice(e.TotalFee) {
			return fmt.Errorf("订单信息不一致")
		}
		if err := e.setTransactionId(ret.TradeNo, ret.BuyerlogonId); err != nil {
			return err
		}

		if ret.TradeStatus == ali_pagepay.TradeStatusSuccess {
			if err := e.paySuccess(ret.SendPayDate); err != nil {
				return err
			}
		} else if ret.TradeStatus == ali_pagepay.TradeStatusWaitBuyerPay {
			//如果用户未支付则关闭该订单
			if err := alipay.CloseOrder(mch, e); err != nil {
				return err
			}
		}
		return nil
	}
}

func (alipay AliPagePay) CheckStartPay(mch ali_pagepay.Merchant, e *orderpay.Entity) error {
	query := ali_pagepay.TradeQuery{
		TradeNo:    e.TransactionId,
		OutTradeNo: e.OutTradeNo,
	}

	if ret, err := ali_pagepay.GetDefault().TradeQuery(mch, query); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TradeNo != e.TransactionId ||
			ret.OutTradeNo != e.OutTradeNo ||
			ret.TotalAmount != toPrice(e.TotalFee) {
			return fmt.Errorf("订单信息不一致")
		}
		if err := e.setTransactionId(ret.TradeNo, ret.BuyerlogonId); err != nil {
			return err
		}

		if ret.TradeStatus == ali_pagepay.TradeStatusSuccess {
			if err := e.paySuccess(ret.SendPayDate); err != nil {
				return err
			}
			return fmt.Errorf("该订单已支付，不能重复发起")
		} else if ret.TradeStatus == ali_pagepay.TradeStatusWaitBuyerPay {
			//如果用户未支付则关闭该订单
			if err := alipay.CloseOrder(mch, e); err != nil {
				return err
			}
		}
		return nil
	}
}
