package orderpay

import (
	"fmt"
	"time"

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

//提交订单
func SubmitOrder(sellerKey string, notifyUrl string, totalFee int, body, detail, accountGuid, code string, timeoutExpress time.Time) error {
	e := &Entity{Data: NewData()}
	return e.submitOrder(sellerKey, notifyUrl, totalFee, body, detail, accountGuid, code, timeoutExpress)
}

//取消订单
func CancelleOrder(code string) (err error) {
	e := &Entity{Data: NewData()}
	if err = e.load(code); err != nil {
		return
	}
	defer func() { _ = e.release() }()

	switch e.Status {
	case Submitted:
		return fmt.Errorf("该订单未发起支付")
	case WaitPay:
		trade := getTrade(e.TradeWay)
		if trade == nil {
			return internal_error.With(fmt.Errorf("该支付方式不支持"))
		}
		if err := e.setCancelled(); err != nil {
			return err
		}
		return trade.TradeClose(e)
	case Paid, Cancelled, Refunding, Refunded:
		return nil
	default:
		err := fmt.Errorf("该订单状态错误")
		return internal_error.Critical(err)
	}
}

//发起支付（只发起已提交订单）
func Pay(code string, tradeWay TradeWay) (interface{}, error) {
	e := &Entity{Data: NewData()}
	if err := e.load(code); err != nil {
		return nil, err
	}
	defer func() { _ = e.release() }()

	//超时不能支付
	if e.PayExpire.Unix() < time.Now().Unix() {
		return nil, fmt.Errorf("订单已过期，不能发起支付")
	}

	switch e.Status {
	case Submitted:
		trade := getTrade(tradeWay)
		if trade == nil {
			return nil, internal_error.With(fmt.Errorf("该支付方式不支持"))
		}
		if tp, err := trade.Pay(e); err != nil {
			return nil, err
		} else {
			if err := e.setPayWay(tradeWay, tp.AppId, tp.MchId); err != nil {
				return nil, internal_error.With(err)
			}

			return tp.Data, nil
		}
	case WaitPay:
		return nil, fmt.Errorf("该订单待支付，不能重复发起支付")
	case Paid, Refunding, Refunded, Cancelled:
		return nil, fmt.Errorf("不能发起支付")
	default:
		err := fmt.Errorf("该订单状态错误")
		return nil, internal_error.Critical(err)
	}
}

//查询支付状态（只查询待支付订单）
func PaySuccess(code string) (res TradeStatus, err error) {
	e := &Entity{Data: NewData()}
	if err = e.load(code); err != nil {
		return
	}
	defer func() { _ = e.release() }()

	switch e.Status {
	case Submitted, Refunding, Refunded:
		return "", fmt.Errorf("无效查询")
	case WaitPay:
		trade := getTrade(e.TradeWay)
		if trade == nil {
			return "", internal_error.With(fmt.Errorf("该支付方式不支持"))
		}
		if tpq, err := trade.QueryPayStatus(e); err != nil {
			return "", err
		} else {
			if err := e.setTransactionId(tpq.TransactionId, tpq.BuyerLogonId); err != nil {
				return "", err
			}
			if tpq.TradeStatus == TradeStatusSuccess {
				if err := e.paySuccess(); err != nil {
					return "", err
				}
			}
			return tpq.TradeStatus, nil
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
		trade := getTrade(e.TradeWay)
		if trade == nil {
			return internal_error.With(fmt.Errorf("该支付方式不支持"))
		}
		if err := e.setCancelled(); err != nil {
			return err
		}
		return trade.TradeClose(e)
	case Paid, Cancelled, Refunding, Refunded:
		return nil
	default:
		err := fmt.Errorf("该订单状态错误")
		return internal_error.Critical(err)
	}
}

//TODO 退款（只退款已支付订单）
func Refund(code string) error {
	return nil
}

//TODO 退款查询（只查询退款中的订单）
func RefundSuccess() error {
	return nil
}
