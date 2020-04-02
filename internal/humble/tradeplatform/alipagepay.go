package tradeplatform

import (
	"fmt"

	"yumi/external/pay/ali_pagepay"
	"yumi/internal/entities/trade"
	"yumi/internal/humble/db"
	"yumi/utils/internal_error"
)

const aliPagePay = trade.TradeWay("alipagepay")

type AliPagePay string

func GetAliPagePay() AliPagePay {
	return ""
}

func (alipay AliPagePay) Pay(e trade.OrderPay) (trade.TradePay, error) {
	ret := trade.TradePay{}

	////订单是否过期
	//if e.PayExpire.Unix() < time.Now().Unix() {
	//	return nil, fmt.Errorf("订单已过期不能发起支付")
	//}

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

	//获取收款商户信息
	mch, err := db.GetAliPayMerchantBySellerKey(e.SellerKey)
	if err != nil {
		return ret, err
	}
	aliMch := ali_pagepay.Merchant{
		AppId:      mch.AppId,
		PrivateKey: mch.PrivateKey,
		PublicKey:  mch.PublicKey,
	}

	if resp, err := ali_pagepay.GetDefault().UnifiedOrder(aliMch, pagePay); err != nil {
		return ret, internal_error.With(err)
	} else {
		ret.AppId = aliMch.AppId
		ret.MchId = resp.SellerId
		ret.Data = resp.PagePayHtml
		return ret, nil
	}
}

func (alipay AliPagePay) QueryPayStatus(e trade.OrderPay) (trade.TradePayQuery, error) {
	tradeInfo := trade.TradePayQuery{}

	tradeQuery := ali_pagepay.TradeQuery{
		TradeNo:    e.TransactionId,
		OutTradeNo: e.OutTradeNo,
	}

	//获取收款商户信息
	mch, err := db.GetAliPayMerchantBySellerKey(e.SellerKey)
	if err != nil {
		return tradeInfo, err
	}
	aliMch := ali_pagepay.Merchant{
		AppId:      mch.AppId,
		PrivateKey: mch.PrivateKey,
		PublicKey:  mch.PublicKey,
	}

	if ret, err := ali_pagepay.GetDefault().TradeQuery(aliMch, tradeQuery); err != nil {
		return tradeInfo, internal_error.With(err)
	} else {
		if ret.OutTradeNo != e.OutTradeNo {
			return tradeInfo, fmt.Errorf("订单号不一致")
		}
		if ret.TotalAmount != toPrice(e.TotalFee) {
			return tradeInfo, internal_error.Critical(fmt.Errorf("订单金额不一致"))
		}
		tradeInfo.TransactionId = ret.TradeNo
		tradeInfo.BuyerLogonId = ret.BuyerlogonId

		switch ret.TradeStatus {
		case ali_pagepay.TradeStatusSuccess:
			tradeInfo.TradeStatus = trade.Success
		case ali_pagepay.TradeStatusWaitBuyerPay:
			tradeInfo.TradeStatus = trade.NotPay
		case ali_pagepay.TradeStatusCloseed:
			tradeInfo.TradeStatus = trade.Closed
		case ali_pagepay.TradeStatusFinished:
			tradeInfo.TradeStatus = trade.Finished
		default:
			err := fmt.Errorf("支付宝状态发生变动，请管理员及时更改")
			return tradeInfo, internal_error.Critical(err)
		}
		return tradeInfo, nil
	}
}

func (alipay AliPagePay) TradeClose(e trade.OrderPay) error {
	tradeClose := ali_pagepay.TradeClose{
		OutTradeNo: e.OutTradeNo,
		TradeNo:    e.TransactionId,
		OperatorId: "sys",
	}

	//获取收款商户信息
	mch, err := db.GetAliPayMerchantBySellerKey(e.SellerKey)
	if err != nil {
		return err
	}
	aliMch := ali_pagepay.Merchant{
		AppId:      mch.AppId,
		PrivateKey: mch.PrivateKey,
		PublicKey:  mch.PublicKey,
	}

	if ret, err := ali_pagepay.GetDefault().TradeClose(aliMch, tradeClose); err != nil {
		return internal_error.With(err)
	} else {
		if ret.TradeNo != e.TransactionId ||
			ret.OutTradeNo != e.OutTradeNo {
			return fmt.Errorf("订单信息不一致")
		}
	}

	return nil
}

func (alipay AliPagePay) Refund(op trade.OrderPay, or trade.OrderRefund) error {
	//TODO
	//获取收款商户信息
	mch, err := db.GetAliPayMerchantBySellerKey(op.SellerKey)
	if err != nil {
		return err
	}
	aliMch := ali_pagepay.Merchant{
		AppId:      mch.AppId,
		PrivateKey: mch.PrivateKey,
		PublicKey:  mch.PublicKey,
	}

	rfd := ali_pagepay.Refund{}

	if ret, err := ali_pagepay.GetDefault().Refund(aliMch, rfd); err != nil {
		return internal_error.With(err)
	}

	return nil
}

func (alipay AliPagePay) QueryRefundStatus(op trade.OrderPay, or trade.OrderRefund) {
	//TODO
}

func toPrice(amount int) string {
	return fmt.Sprintf("%d.%02d", amount/100, amount%100)
}

//可以添加对账服务，获取支付评价
//func (alipay AliPagePay) PayNotify(mch ali_pagepay.Merchant, rawQuery string) (trade.Status, error) {
//	reqNotify, err := ali_pagepay.ParseQuery(rawQuery)
//	if err != nil {
//		return "", internal_error.With(err)
//	}
//
//	e := trade.Entity{}
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
//	switch reqNotify.Status {
//	case ali_pagepay.Success:
//		return trade.Success, nil
//	case ali_pagepay.TradeStatusWaitBuyerPay:
//		return trade.NotPay, nil
//	case ali_pagepay.TradeStatusCloseed:
//		return trade.Closed, nil
//	case ali_pagepay.Finished:
//		return trade.Finished, nil
//	default:
//		err := fmt.Errorf("支付宝状态发生变动，请管理员及时更改")
//		return "", internal_error.Critical(err)
//	}
//	//if reqNotify.Status == ali_pagepay.Success {
//	//	if order.OrderStatus == Paid {
//	//		return nil
//	//	}
//	//	if order.OrderStatus == WaitPay {
//	//		if err := order.PaySuccess(reqNotify.GmtPayment); err != nil {
//	//			return err
//	//		}
//	//		return nil
//	//	}
//	//	if order.OrderStatus == Submitted ||
//	//		order.OrderStatus == Cancelled {
//	//		//支付宝没有处理好，关闭订单和支付执行顺序。
//	//		//TODO 这种情况可能引发重复支付，应该记录紧急日志，并通知管理员处理。
//	//	}
//	//}
//	//
//	//return fmt.Errorf("交易状态错误")
//}

//func (alipay AliPagePay) RefundNotify() {
//	//TODO
//}
