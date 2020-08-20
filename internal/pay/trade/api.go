package trade

import (
	"fmt"
	"net/http"
	"time"

	"yumi/pkg/ecode"
	"yumi/pkg/log"
)

const timeFormat = "2006-01-02 15:04:05.999"

//SubmitOrderPay 提交支付订单
func SubmitOrderPay(buyerAccountGUID, sellerKey string, totalFee int, body, detail string,
	timeoutExpress time.Time) (string, error) {
	e, err := NewEntityByPayCode("")
	if err != nil {
		return "", err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	now := time.Now()
	code := getCode(OrderPayCode)
	return code, e.dataOp.Submit(buyerAccountGUID, sellerKey, "", "", totalFee, body, detail, timeoutExpress, now, code, Submitted)
}

func sendPay(tradeWay Way, e *Entity, payExpire time.Time, clientIP, notifyURL string) (string, error) {
	outTradeNo := getOutTradeNo()
	e.op.OutTradeNo = outTradeNo
	if err := e.dataOp.SetOutTradeNo(outTradeNo, notifyURL); err != nil {
		return "", err
	}

	trade := getTrade(tradeWay)
	if trade == nil {
		return "", ecode.NotSupportTradeWay
	}

	e.op.PayExpire = payExpire
	e.op.SpbillCreateIP = clientIP
	e.op.NotifyURL = notifyURL

	tp, err := trade.Pay(e.op)
	if err != nil {
		return "", err
	}
	if err := e.dataOp.SetWaitPay(tradeWay, tp.AppID, tp.MchID, clientIP, payExpire, WaitPay); err != nil {
		return "", err
	}

	return tp.Data, nil
}

//Pay 发起支付
func Pay(code string, tradeWay Way, clientIP, notifyURL string, payExpire time.Time) (string, error) {
	e, err := NewEntityByPayCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	//订单超时不能发起支付
	if e.op.TimeoutExpress.Format(timeFormat) < time.Now().Format(timeFormat) {
		return "", ecode.OrderPayTimeout
	}

	switch e.op.Status {
	case Submitted:
		return sendPay(tradeWay, e, payExpire, clientIP, notifyURL)
	case WaitPay:
		//查询当前支付状态
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			return "", ecode.NotSupportTradeWay
		}
		//查询之前支付状态
		if tpq, err := trade.QueryPayStatus(e.op); err != nil {
			return "", err
		} else if tpq.TradeStatus == NotPay {
			//关闭之前支付
			if err := trade.TradeClose(e.op); err != nil {
				return "", err
			}
			//发起支付
			return sendPay(tradeWay, e, payExpire, clientIP, notifyURL)
		} else if tpq.TradeStatus == Closed {
			//发起支付
			return sendPay(tradeWay, e, payExpire, clientIP, notifyURL)
		} else {
			return "", ecode.InvalidSendPay
		}
	default:
		return "", ecode.InvalidSendPay
	}
}

//CancelOrderPay 取消支付订单
func CancelOrderPay(code string) (err error) {
	e, err := NewEntityByPayCode(code)
	if err != nil {
		return err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	switch e.op.Status {
	case Submitted:
		return e.dataOp.SetCancelled(time.Now(), Cancelled)
	case WaitPay:
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			return ecode.NotSupportTradeWay
		}
		if err := trade.TradeClose(e.op); err != nil {
			return err
		}

		return e.dataOp.SetCancelled(time.Now(), Cancelled)
	default:
		return ecode.InvalidCancelOrderPay
	}
}

//PaySuccess 查询支付成功（只查询待支付订单）
func PaySuccess(code string) (res Status, err error) {
	e, err := NewEntityByPayCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	switch e.op.Status {
	case WaitPay:
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			return "", ecode.NotSupportTradeWay
		}
		tpq, err := trade.QueryPayStatus(e.op)
		if err != nil {
			return "", err
		}
		if tpq.TradeStatus == Success {
			if err := e.dataOp.SetSuccess(time.Now(), tpq.TransactionID, tpq.BuyerLogonID, Paid); err != nil {
				return "", err
			}
		}
		return tpq.TradeStatus, nil
	case Paid:
		return Success, nil
	case Cancelled:
		return Closed, nil
	default:
		return "", ecode.InvalidQueryPay
	}
}

//PayStatus 查询支付状态
func PayStatus(code string) (res Status, err error) {
	e, err := NewEntityByPayCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	trade := getTrade(e.op.TradeWay)
	if trade == nil {
		return "", ecode.NotSupportTradeWay
	}
	tpq, err := trade.QueryPayStatus(e.op)
	if err != nil {
		return "", err
	}

	return tpq.TradeStatus, nil
}

//CloseTrade 关闭交易（只关闭待支付订单）
func CloseTrade(code string) (err error) {
	e, err := NewEntityByPayCode(code)
	if err != nil {
		return err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	switch e.op.Status {
	case WaitPay:
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			return ecode.NotSupportTradeWay
		}
		if err := e.dataOp.SetSubmitted(Submitted); err != nil {
			return err
		}
		return trade.TradeClose(e.op)
	default:
		return ecode.InvalidCloseTrade
	}
}

//PayNotify 支付通知(待支付时处理通知)
func PayNotify(way Way, resp http.ResponseWriter, req *http.Request) (string, Status) {
	trade := getTrade(way)
	//解析通知参数
	ret, err := trade.PayNotifyReq(req)
	if err != nil {
		trade.PayNotifyResp(err, resp)
		return "", ""
	}
	e, err := NewEntityByPayCode(ret.OrderPayCode)
	if err != nil {
		log.Error(err)
		trade.PayNotifyResp(err, resp)
		return "", ""
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	switch e.op.Status {
	case WaitPay:
		//检查通知参数
		if err := trade.PayNotifyCheck(e.op, ret.ReqData); err != nil {
			trade.PayNotifyResp(err, resp)
			return "", ""
		}

		//查询支付状态
		tpq, err := trade.QueryPayStatus(e.op)
		if err != nil {
			trade.PayNotifyResp(err, resp)
			return "", ""
		}
		if tpq.TradeStatus == Success {
			if err := e.dataOp.SetSuccess(time.Now(), tpq.TransactionID, tpq.BuyerLogonID, Paid); err != nil {
				err = fmt.Errorf("服务器内部错误")
				trade.PayNotifyResp(err, resp)
				return "", ""
			}
			trade.PayNotifyResp(nil, resp)
			return ret.OrderPayCode, Success
		}
		log.Error(fmt.Errorf("非法请求"))
		return "", tpq.TradeStatus
	case Paid:
		trade.PayNotifyResp(nil, resp)
		return ret.OrderPayCode, Success
	default:
		//返回成功，停止通知
		trade.PayNotifyResp(nil, resp)
		return "", ""
	}
}

//Refund 提交退款订单（只退款已支付订单）
func Refund(orderPayCode, notifyURL string, refundAccountGUID string, refundFee int, refundDesc string, submitTime, timeoutExpress time.Time) (string, error) {
	e, err := NewEntityByPayCode(orderPayCode)
	if err != nil {
		return "", err
	}

	switch e.op.Status {
	//支付订单必须为已支付
	case Paid:
		//是否存在已提交或退款中的退款订单
		if exist, err := e.dataOr.ExistRefundingOrSubmitted(e.op.Code); err != nil {
			return "", err
		} else if exist {
			return "", ecode.CannotRepeatSendRefund
		}

		//退款金额是否超额
		count, countFee := 0, 0
		if count, countFee, err = e.dataOr.GetRefundFee(e.op.Code); err != nil {
			return "", err
		} else if countFee+refundFee > e.op.TotalFee {
			err = fmt.Errorf("退款总额超过支付金额")
			return "", ecode.ServerErr(err)
		}

		code := getCode(OrderRefundCode)
		outRefundNo := code
		if err := e.dataOr.Submit(code, orderPayCode, count, notifyURL, refundAccountGUID, e.op.TradeWay, outRefundNo, refundFee,
			refundDesc, submitTime, timeoutExpress, Submitted); err != nil {
			return "", err
		}

		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			_ = e.dataOr.SetError(time.Now(), "tradeway（交易方式）不存在", Error)
			err = fmt.Errorf("refundcode: %s，tradeway（交易方式）不存在", e.or.Code)
			return "", ecode.ServerErr(err)
		}

		_ = e.ReleaseOrderPay()
		e, err := NewEntityByRefundCode(code)
		if err != nil {
			return "", err
		}
		defer func() { _ = e.ReleaseOrderRefund() }()
		if err := trade.Refund(e.op, e.or); err != nil {
			return "", err
		}

		if err := e.dataOr.SetRefunding(Refunding); err != nil {
			return "", err
		}
		return code, nil
	default:
		return "", ecode.InvalidSendRefund
	}
}

//RefundSuccess 退款查询（只查询退款中的订单）
func RefundSuccess(code string) (res Status, err error) {
	e, err := NewEntityByRefundCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.ReleaseOrderRefund() }()

	switch e.or.Status {
	case Refunding:
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			_ = e.dataOr.SetError(time.Now(), "tradeway（交易方式）不存在", Error)
			err = fmt.Errorf("refundcode: %s，tradeway（交易方式）不存在", e.or.Code)
			return "", ecode.ServerErr(err)
		}

		ret, err := trade.QueryRefundStatus(e.op, e.or)
		if err != nil {
			return "", err
		}
		if ret.TradeStatus == Success {
			if err := e.dataOr.SetRefunded(ret.RefundID, time.Now(), Refunded); err != nil {
				log.Error(err)
				return "", err
			}
		}
		return ret.TradeStatus, nil
	case Refunded:
		return "", nil
	default:
		return "", ecode.InvalidQueryRefund
	}
}

//RefundNotify 退款通知
func RefundNotify(way Way, resp http.ResponseWriter, req *http.Request) Status {
	trade := getTrade(way)
	if trade == nil {
		err := fmt.Errorf("tradeway（交易方式）不存在")
		log.Error(err)
		panic(err)
	}

	ret, err := trade.RefundNotifyReq(req)
	if err != nil {
		trade.RefundNotifyResp(err, resp)
		return ""
	}
	e, err := NewEntityByRefundCode(ret.OrderRefundCode)
	if err != nil {
		trade.PayNotifyResp(err, resp)
		return ""
	}
	defer func() { _ = e.ReleaseOrderRefund() }()

	if err := trade.RefundNotifyCheck(e.op, e.or, ret.ReqData); err != nil {
		trade.PayNotifyResp(err, resp)
		return ""
	}

	status, err := RefundSuccess(ret.OrderRefundCode)
	if err != nil {
		trade.PayNotifyResp(err, resp)
		return ""
	}
	if status == Success {
		trade.RefundNotifyResp(nil, resp)
		return status
	}

	return status
}

//SetTimeout 将超时订单设置为已取消
func SetTimeout(code string) error {
	e, err := NewEntityByPayCode(code)
	if err != nil {
		return err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	if e.op.TimeoutExpress.Format(timeFormat) < time.Now().Format(timeFormat) {
		return e.dataOp.SetCancelled(time.Now(), Cancelled)
	}

	return nil
}

////申请退款（只退款已提交退款订单）
//func Refund(code string) error {
//	e, err := NewEntityByRefundCode(code)
//	if err != nil {
//		return InternalError
//	}
//	defer func() { _ = e.ReleaseOrderRefund() }()
//
//	if e.op.TimeoutExpress.Format(timeFormat) < time.Now().Format(timeFormat) {
//		return nil, fmt.Errorf("订单已过期，不能发起支付")
//	}
//
//	if e.or.Status == Submitted {
//		trade := getTrade(e.op.TradeWay)
//		if trade == nil {
//			if err := e.dataOr.SetError(time.Now(), "tradeway（交易方式）不存在", Error); err != nil {
//				logs.Error(err)
//				logs.Error(fmt.Sprintf("refundcode: %s，tradeway（交易方式）不存在", e.or.Code))
//				return InternalError
//			}
//			return fmt.Errorf("订单错误")
//		}
//
//		return trade.Refund(e.op, e.or)
//	}
//
//	return fmt.Errorf("不能发起退款申请")
//}
