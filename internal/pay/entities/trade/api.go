package trade

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"yumi/utils/log"
)

const timeFormat = "2006-01-02 15:04:05.999"

var (
	InternalError = errors.New("服务器内部错误")
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
func SubmitOrderPay(buyerAccountGuid, sellerKey string, totalFee int, body, detail string,
	timeoutExpress time.Time) (string, error) {
	e, err := NewEntityByPayCode("")
	if err != nil {
		return "", err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	now := time.Now()
	code := getCode(OrderPayCode)
	return code, e.dataOp.Submit(buyerAccountGuid, sellerKey, "", "", totalFee, body, detail, timeoutExpress, now, code, Submitted)
}

func sendPay(tradeWay Way, e *Entity, payExpire time.Time, clientIp, notifyUrl string) (string, error) {
	outTradeNo := getOutTradeNo()
	e.op.OutTradeNo = outTradeNo
	if err := e.dataOp.SetOutTradeNo(outTradeNo, notifyUrl); err != nil {
		return "", err
	}

	trade := getTrade(tradeWay)
	if trade == nil {
		return "", fmt.Errorf("该支付方式不支持")
	}

	e.op.PayExpire = payExpire
	e.op.SpbillCreateIp = clientIp
	e.op.NotifyUrl = notifyUrl
	if tp, err := trade.Pay(e.op); err != nil {
		return "", err
	} else {
		if err := e.dataOp.SetWaitPay(tradeWay, tp.AppId, tp.MchId, clientIp, payExpire, WaitPay); err != nil {
			return "", err
		}

		return tp.Data, nil
	}
}

//发起支付
func Pay(code string, tradeWay Way, clientIp, notifyUrl string, payExpire time.Time) (string, error) {
	e, err := NewEntityByPayCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	//订单超时不能发起支付
	if e.op.TimeoutExpress.Format(timeFormat) < time.Now().Format(timeFormat) {
		return "", fmt.Errorf("订单已过期，不能发起支付")
	}

	switch e.op.Status {
	case Submitted:
		return sendPay(tradeWay, e, payExpire, clientIp, notifyUrl)
	case WaitPay:
		//查询当前支付状态
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			return "", fmt.Errorf("该支付方式不支持")
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
			return sendPay(tradeWay, e, payExpire, clientIp, notifyUrl)
		} else if tpq.TradeStatus == Closed {
			//发起支付
			return sendPay(tradeWay, e, payExpire, clientIp, notifyUrl)
		} else {
			return "", fmt.Errorf("不能发起支付")
		}
	default:
		return "", fmt.Errorf("不能发起支付")
	}
}

//取消支付订单
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
			return fmt.Errorf("该支付方式不支持")
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
	e, err := NewEntityByPayCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	switch e.op.Status {
	case WaitPay:
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			return "", fmt.Errorf("该支付方式不支持")
		}
		if tpq, err := trade.QueryPayStatus(e.op); err != nil {
			return "", err
		} else {
			if tpq.TradeStatus == Success {
				if err := e.dataOp.SetSuccess(time.Now(), tpq.TransactionId, tpq.BuyerLogonId, Paid); err != nil {
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

//查询支付状态
func PayStatus(code string) (res Status, err error) {
	e, err := NewEntityByPayCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.ReleaseOrderPay() }()

	trade := getTrade(e.op.TradeWay)
	if trade == nil {
		return "", fmt.Errorf("该支付方式不支持")
	}
	if tpq, err := trade.QueryPayStatus(e.op); err != nil {
		return "", err
	} else {
		return tpq.TradeStatus, nil
	}
}

//关闭交易（只关闭待支付订单）
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
			return fmt.Errorf("该支付方式不支持")
		}
		if err := e.dataOp.SetSubmitted(Submitted); err != nil {
			return err
		}
		return trade.TradeClose(e.op)
	default:
		return fmt.Errorf("不能关闭交易")
	}
}

//支付通知(待支付时处理通知)
func PayNotify(way Way, resp http.ResponseWriter, req *http.Request) (string, Status) {
	trade := getTrade(way)
	//解析通知参数
	if ret, err := trade.PayNotifyReq(req); err != nil {
		trade.PayNotifyResp(err, resp)
		return "", ""
	} else {
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
			if tpq, err := trade.QueryPayStatus(e.op); err != nil {
				trade.PayNotifyResp(err, resp)
				return "", ""
			} else {
				if tpq.TradeStatus == Success {
					if err := e.dataOp.SetSuccess(time.Now(), tpq.TransactionId, tpq.BuyerLogonId, Paid); err != nil {
						err = fmt.Errorf("服务器内部错误")
						trade.PayNotifyResp(err, resp)
						return "", ""
					}
					trade.PayNotifyResp(nil, resp)
					return ret.OrderPayCode, Success
				} else {
					log.Error(fmt.Errorf("ip：%s,非法请求", utils.ClientIp(req)))
					return "", tpq.TradeStatus
				}
			}
		case Paid:
			trade.PayNotifyResp(nil, resp)
			return ret.OrderPayCode, Success
		default:
			//返回成功，停止通知
			trade.PayNotifyResp(nil, resp)
			return "", ""
		}
	}
}

//提交退款订单（只退款已支付订单）
func Refund(orderPayCode, notifyUrl string, refundAccountGuid string, refundFee int, refundDesc string, submitTime, timeoutExpress time.Time) (string, error) {
	e, err := NewEntityByPayCode(orderPayCode)
	if err != nil {
		return "", InternalError
	}

	switch e.op.Status {
	//支付订单必须为已支付
	case Paid:
		//是否存在已提交或退款中的退款订单
		if exist, err := e.dataOr.ExistRefundingOrSubmitted(e.op.Code); err != nil {
			return "", InternalError
		} else if exist {
			return "", fmt.Errorf("已有退款订单处理中，不能发起订单")
		}

		//退款金额是否超额
		count, countFee := 0, 0
		if count, countFee, err = e.dataOr.GetRefundFee(e.op.Code); err != nil {
			return "", InternalError
		} else if countFee+refundFee > e.op.TotalFee {
			return "", fmt.Errorf("退款总额超过支付金额")
		}

		code := getCode(OrderRefundCode)
		outRefundNo := code
		if err := e.dataOr.Submit(code, orderPayCode, count, notifyUrl, refundAccountGuid, e.op.TradeWay, outRefundNo, refundFee,
			refundDesc, submitTime, timeoutExpress, Submitted); err != nil {
			return "", InternalError
		}

		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			if err := e.dataOr.SetError(time.Now(), "tradeway（交易方式）不存在", Error); err != nil {
				log.Error(err)
				log.Error(fmt.Sprintf("refundcode: %s，tradeway（交易方式）不存在", e.or.Code))
				return "", InternalError
			}
			return "", fmt.Errorf("订单错误")
		}

		_ = e.ReleaseOrderPay()
		e, err := NewEntityByRefundCode(code)
		if err != nil {
			return "", InternalError
		}
		defer func() { _ = e.ReleaseOrderRefund() }()
		if err := trade.Refund(e.op, e.or); err != nil {
			return "", err
		}

		if err := e.dataOr.SetRefunding(Refunding); err != nil {
			return "", InternalError
		}
		return code, nil
	default:
		return "", fmt.Errorf("该支付订单不能发起退款")
	}
}

//退款查询（只查询退款中的订单）
func RefundSuccess(code string) (res Status, err error) {
	e, err := NewEntityByRefundCode(code)
	if err != nil {
		return "", InternalError
	}
	defer func() { _ = e.ReleaseOrderRefund() }()

	switch e.or.Status {
	case Refunding:
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			err := fmt.Errorf("refundcode: %s，tradeway（交易方式）不存在", e.or.Code)
			if err := e.dataOr.SetError(time.Now(), err.Error(), Error); err != nil {
				log.Error(err)
				return "", InternalError
			}
			return "", err
		}

		ret, err := trade.QueryRefundStatus(e.op, e.or)
		if err != nil {
			return "", err
		}
		if ret.TradeStatus == Success {
			if err := e.dataOr.SetRefunded(ret.RefundId, time.Now(), Refunded); err != nil {
				log.Error(err)
				return "", err
			}
		}
		return ret.TradeStatus, nil
	case Refunded:
		return "", nil
	default:
		return "", fmt.Errorf("无效查询")
	}
}

//退款通知
func RefundNotify(way Way, resp http.ResponseWriter, req *http.Request) Status {
	trade := getTrade(way)
	if trade == nil {
		err := fmt.Errorf("tradeway（交易方式）不存在")
		log.Error(err)
		panic(err)
		return ""
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

//将超时订单设置为已取消
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