package orderpay

import (
	"fmt"
	"time"

	"yumi/external/pay/ali_pagepay"
	"yumi/external/pay/wx_nativepay"
	"yumi/internal/humble/threepay"
	"yumi/utils/internal_error"
)

/**
 * 业务对象接口
 * 供用例（use case）对象调用，对外开放
 */

type PayResult struct {
	AliPayHtml  []byte
	WxPayBizUrl string
}

//发起支付（只发起已提交订单）
func Pay(code string) (res PayResult, err error) {
	e := &Entity{Data: NewData()}
	if err = e.load(code); err != nil {
		return
	}
	defer func() { _ = e.release() }()

	//订单是否过期
	if e.PayExpire.Unix() < time.Now().Unix() {
		return res, fmt.Errorf("订单已过期，不能发起支付")
	}

	switch e.Status {
	case Submitted:
		switch e.PayWay {
		case PayWayAliPagePay:
			if mch, err := getAliApp(e.AppId); err != nil {
				return res, err
			} else {
				if htmlByte, err := threepay.GetAliPagePay().Pay(mch, e); err != nil {
					return res, err
				} else {
					res.AliPayHtml = htmlByte
					return res, nil
				}

			}
		case PayWayWxNative1Pay:
			if mch, err := getWxApp(e.AppId); err != nil {
				return res, err
			} else {
				if bizUrl, err := threepay.GetWxNative1().Pay(mch, e); err != nil {
					return res, err
				} else {
					res.WxPayBizUrl = bizUrl
					return res, nil
				}
			}
		default:
			err := fmt.Errorf("不支持的支付方式")
			return res, internal_error.Critical(err)
		}

	case WaitPay:
		return res, fmt.Errorf("该订单待支付，不能重复发起支付")

	case Paid, Refunding, Refunded, Cancelled:
		return res, fmt.Errorf("不能发起支付")

	default:
		err := fmt.Errorf("该订单状态错误")
		return res, internal_error.Critical(err)
	}
}

//查询支付状态（只查询待支付订单）
func QueryPayStatus(code string) (res TradeStatus, err error) {
	e := &Entity{Data: NewData()}
	if err = e.load(code); err != nil {
		return
	}
	defer func() { _ = e.release() }()

	switch e.Status {
	case Submitted, Refunding, Refunded:
		return "", fmt.Errorf("无效查询")

	case WaitPay:
		switch e.PayWay {
		case PayWayAliPagePay:
			if mch, err := getAliApp(e.AppId); err != nil {
				return "", err
			} else {
				return threepay.GetAliPagePay().QueryPayStatus(mch, e)
			}
		case PayWayWxNative1Pay:
			if mch, err := getWxApp(e.AppId); err != nil {
				return "", err
			} else {
				return threepay.GetWxNative1().QueryPayStatus(mch, e)
			}
		default:
			err := fmt.Errorf("不支持的支付方式")
			return "", internal_error.Critical(err)
		}

	case Paid:
		return TradeStatusSuccess, nil

	case Cancelled:
		return TradeStatusClosed, nil

	default:
		err := fmt.Errorf("该订单状态错误")
		return "", internal_error.Critical(err)
	}
}

//关闭交易（只关闭待支付订单）
func CloseTrade(code string) (err error) {
	e := &Entity{Data: NewData()}
	if err = e.load(code); err != nil {
		return
	}
	defer func() { _ = e.release() }()

	switch e.Status {
	case Submitted:
		return fmt.Errorf("该订单未发起支付")

	case WaitPay:
		//TODO
		return nil
	case Paid, Cancelled, Refunding, Refunded:
		return nil

	default:
		err := fmt.Errorf("该订单状态错误")
		return internal_error.Critical(err)
	}
}

//TODO 退款（只退款已支付订单）
func Refund() error {
	return nil
}

//TODO 退款查询（只查询退款中的订单）
func QueryRefundStatus() error {
	return nil
}

//======================================================================================================================
func getWxApp(code string) (wx_nativepay.Merchant, error) {
	mch := wx_nativepay.Merchant{}
	if ret, err := GetWxApp(code); err != nil {
		return mch, err
	} else {
		mch.AppId = ret.AppId
		mch.MchId = ret.MchId
		mch.PrivateKey = ret.PrivateKey
		return mch, nil
	}
}

func getAliApp(code string) (ali_pagepay.Merchant, error) {
	mch := ali_pagepay.Merchant{}
	if ret, err := GetAliApp(code); err != nil {
		return mch, err
	} else {
		mch.AppId = ret.AppId
		mch.PublicKey = ret.PublicKey
		mch.PrivateKey = ret.PrivateKey
		return mch, nil
	}
}
