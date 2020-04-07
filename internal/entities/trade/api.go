package trade

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

//提交支付订单
func SubmitOrderPay(accountGuid, sellerKey string, notifyUrl string, totalFee int, body, detail string,
	timeoutExpress time.Time) error {
	e := &Entity{dataOp: NewDataOrderPay()}
	now := time.Now()
	outTradeNo := getOutTradeNo()
	code := getCode(0)
	return e.dataOp.Submit(accountGuid, sellerKey, outTradeNo, notifyUrl, totalFee, body, detail, timeoutExpress, now, code, Submitted)
}

//发起支付（只发起已提交订单）
func Pay(code string, tradeWay TradeWay) (interface{}, error) {
	e := &Entity{dataOp: NewDataOrderPay()}
	if err := e.loadOrderPay(code); err != nil {
		return nil, err
	}
	defer func() { _ = e.releaseOrderPay() }()

	//超时不能支付
	if e.op.PayExpire.Unix() < time.Now().Unix() {
		return nil, fmt.Errorf("订单已过期，不能发起支付")
	}

	switch e.op.Status {
	case Submitted:
		trade := getTrade(tradeWay)
		if trade == nil {
			return nil, internal_error.With(fmt.Errorf("该支付方式不支持"))
		}
		if tp, err := trade.Pay(e.op); err != nil {
			return nil, err
		} else {
			if err := e.dataOp.SetPayWay(tradeWay, tp.AppId, tp.MchId, WaitPay); err != nil {
				return nil, internal_error.With(err)
			}

			return tp.Data, nil
		}
	default:
		return nil, fmt.Errorf("不能发起支付")
	}
}

//取消支付订单
func CancelleOrderPay(code string) (err error) {
	e := &Entity{dataOp: NewDataOrderPay()}
	if err = e.loadOrderPay(code); err != nil {
		return
	}
	defer func() { _ = e.releaseOrderPay() }()

	switch e.op.Status {
	case Submitted:
		return e.dataOp.SetCancelled(time.Now(), Cancelled)
	case WaitPay:
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			return internal_error.With(fmt.Errorf("该支付方式不支持"))
		}
		if err := trade.TradeClose(e.op); err != nil {
			return err
		}

		return e.dataOp.SetCancelled(time.Now(), Cancelled)
	default:
		return fmt.Errorf("不能取消订单")
	}
}

//查询支付成功（只查询待支付订单）
func PaySuccess(code string) (res Status, err error) {
	e := &Entity{dataOp: NewDataOrderPay()}
	if err = e.loadOrderPay(code); err != nil {
		return
	}
	defer func() { _ = e.releaseOrderPay() }()

	switch e.op.Status {
	case WaitPay:
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			return "", internal_error.With(fmt.Errorf("该支付方式不支持"))
		}
		if tpq, err := trade.QueryPayStatus(e.op); err != nil {
			return "", err
		} else {
			if err := e.dataOp.SetTransactionId(tpq.TransactionId, tpq.BuyerLogonId); err != nil {
				return "", err
			}
			if tpq.TradeStatus == Success {
				if err := e.dataOp.PaySuccess(time.Now(), Paid); err != nil {
					return "", err
				}
			}
			return tpq.TradeStatus, nil
		}
	case Paid:
		return Success, nil
	case Cancelled:
		return Closed, nil
	default:
		return "", fmt.Errorf("无效查询")
	}
}

//关闭交易（只关闭待支付订单）
func CloseTrade(code string) (err error) {
	e := &Entity{dataOp: NewDataOrderPay()}
	if err = e.loadOrderPay(code); err != nil {
		return
	}
	defer func() { _ = e.releaseOrderPay() }()

	switch e.op.Status {
	case WaitPay:
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			return internal_error.With(fmt.Errorf("该支付方式不支持"))
		}
		if err := e.dataOp.SetCancelled(time.Now(), Cancelled); err != nil {
			return err
		}
		return trade.TradeClose(e.op)
	default:
		return fmt.Errorf("不能关闭交易")
	}
}

//TODO 提交退款订单（只退款已支付订单）
func SubmitOrderRefund() error {
	return nil
}

//TODO 取消退款订单（只取消已提交退款订单）
func CancelleOrderRefund() error {
	return nil
}

//TODO 退款（只退款已提交退款订单）
func Refund(code string) error {
	return nil
}

//TODO 退款查询（只查询退款中的订单）
func RefundSuccess() error {
	return nil
}
