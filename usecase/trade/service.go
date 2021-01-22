package trade

import (
	"fmt"
	"net/http"
	"time"

	"yumi/pkg/codes"
	"yumi/pkg/ecode"
	"yumi/pkg/log"
	"yumi/pkg/status"
	"yumi/usecase/trade/entity"
)

const timeFormat = "2006-01-02 15:04:05.999"

// Service ...
type Service struct {
}

// New a Service object
func New() (*Service, error) {
	return &Service{}, nil
}

//CreateOrderPay 提交支付订单
func (s *Service) CreateOrderPay(req CreateOrderPayRequest) (resp CreateOrderPayResponse, err error) {
	data := GetData()

	attr := entity.OrderPayAttribute{}
	req.Attribute(&attr)

	err = data.CreateOrderPay(attr)
	resp.Code = attr.Code

	return
}

//Pay 发起支付
func Pay(req PayRequest) (PayResponse, error) {
	data := GetData()
	resp := PayResponse{}

	dataOp, err := data.GetOrderPay(req.Code)
	if err != nil {
		return resp, err
	}
	attr := dataOp.Attribute()

	op := entity.NewOrderPay(dataOp.Attribute())
	err = op.CanPay()
	if err != nil {
		return resp, err
	}

	switch attr.Status {
	case entity.Submitted:
		goto SendPay
	case entity.WaitPay:
		//查询当前支付状态
		trade := getTrade(Way(req.TradeWay))
		if trade == nil {
			return resp, err 
		}
		//查询之前支付状态
		tpq, err := trade.QueryPayStatus(op)
		if err != nil {
			return resp, err
		}
		if tpq.TradeStatus == StatusTradePlatformNotPay {
			if err := trade.TradeClose(op); err != nil {
				return resp, err 
			}
			goto SendPay
		}
		if tpq.TradeStatus == StatusTradePlatformClosed {
			goto SendPay
		}
		return resp, status
	default:
		return resp, err
	}

SendPay:
	req.Attribute(attr)

	trade := getTrade(Way(req.TradeWay))
	if trade == nil {
		return resp, 
	}

	tp, err := trade.Pay(op)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

//CancelOrderPay 取消支付订单
func CancelOrderPay(code string) (err error) {
	e, err := newEntityByPayCode(code)
	if err != nil {
		return err
	}
	defer func() { _ = e.release() }()

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
	e, err := newEntityByPayCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.release() }()

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
		if tpq.TradeStatus == StatusTradePlatformSuccess {
			if err := e.dataOp.SetSuccess(time.Now(), tpq.TransactionID, tpq.BuyerLogonID, Paid); err != nil {
				return "", err
			}
			return Paid, nil
		}
		return WaitPay, nil
	case Paid:
		return Paid, nil
	case Cancelled:
		return Cancelled, nil
	default:
		return "", ecode.InvalidQueryPay
	}
}

//PayStatus 查询支付状态
func PayStatus(code string) (res StatusTradePlatform, err error) {
	e, err := newEntityByPayCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.release() }()

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
	e, err := newEntityByPayCode(code)
	if err != nil {
		return err
	}
	defer func() { _ = e.release() }()

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
	if trade == nil {
		err := fmt.Errorf("交易方式不存在：%s", way)
		log.Error(err)
		trade.PayNotifyResp(err, resp)
		return "", ""
	}

	//解析通知参数
	ret, err := trade.PayNotifyReq(req)
	if err != nil {
		trade.PayNotifyResp(err, resp)
		return "", ""
	}
	e, err := newEntityByPayCode(ret.OrderPayCode)
	if err != nil {
		log.Error(err)
		trade.PayNotifyResp(err, resp)
		return "", ""
	}
	defer func() { _ = e.release() }()

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
		if tpq.TradeStatus == StatusTradePlatformSuccess {
			if err := e.dataOp.SetSuccess(time.Now(), tpq.TransactionID, tpq.BuyerLogonID, Paid); err != nil {
				err = fmt.Errorf("服务器内部错误")
				trade.PayNotifyResp(err, resp)
				return "", ""
			}
			trade.PayNotifyResp(nil, resp)
			return ret.OrderPayCode, Paid
		}
		log.Error(fmt.Errorf("非法请求"))
		return "", WaitPay
	case Paid:
		trade.PayNotifyResp(nil, resp)
		return ret.OrderPayCode, Paid
	default:
		//返回成功，停止通知
		trade.PayNotifyResp(nil, resp)
		return "", ""
	}
}

//Refund 提交退款订单（只退款已支付订单）
func Refund(orderPayCode, notifyURL string, refundAccountGUID string, refundFee int, refundDesc string, submitTime, timeoutExpress time.Time) (string, error) {
	e, err := newEntityByPayCode(orderPayCode)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.release() }()

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

		// 获取交易方式
		trade := getTrade(e.op.TradeWay)
		if trade == nil {
			_ = e.dataOr.SetError(time.Now(), "tradeway（交易方式）不存在", Error)
			err = fmt.Errorf("refundcode: %s，tradeway（交易方式）不存在", e.or.Code)
			return "", ecode.ServerErr(err)
		}

		// 提交退款订单
		code := getCode(OrderRefundCode)
		outRefundNo := code
		if err := e.dataOr.Submit(code, orderPayCode, count, notifyURL, refundAccountGUID, e.op.TradeWay, outRefundNo, refundFee,
			refundDesc, submitTime, timeoutExpress, Submitted); err != nil {
			return "", err
		}
		e.dataOr, err = newDataOrderRefund(code)
		e.or = e.dataOr.Entity()

		// 退款
		if err := trade.Refund(e.op, e.or); err != nil {
			return "", err
		}

		// 设置退款中
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
	e, err := newEntityByRefundCode(code)
	if err != nil {
		return "", err
	}
	defer func() { _ = e.release() }()

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
		if ret.TradeStatus == StatusTradePlatformSuccess {
			if err := e.dataOr.SetRefunded(ret.RefundID, time.Now(), Refunded); err != nil {
				log.Error(err)
				return "", err
			}
			return Refunded, nil
		}
		return Refunding, nil
	case Refunded:
		return Refunded, nil
	default:
		return "", ecode.InvalidQueryRefund
	}
}

//RefundNotify 退款通知
func RefundNotify(way Way, resp http.ResponseWriter, req *http.Request) Status {
	trade := getTrade(way)
	if trade == nil {
		err := fmt.Errorf("交易方式不存在：%s", way)
		log.Error(err)
		trade.RefundNotifyResp(err, resp)
		return ""
	}

	ret, err := trade.RefundNotifyReq(req)
	if err != nil {
		trade.RefundNotifyResp(err, resp)
		return ""
	}
	e, err := newEntityByRefundCode(ret.OrderRefundCode)
	if err != nil {
		trade.PayNotifyResp(err, resp)
		return ""
	}
	defer func() { _ = e.release() }()

	if err := trade.RefundNotifyCheck(e.op, e.or, ret.ReqData); err != nil {
		trade.PayNotifyResp(err, resp)
		return ""
	}

	status, err := RefundSuccess(ret.OrderRefundCode)
	if err != nil {
		trade.PayNotifyResp(err, resp)
		return ""
	}
	if status == Refunded {
		trade.RefundNotifyResp(nil, resp)
		return Refunded
	}

	return status
}

//SetTimeout 将超时订单设置为已取消
func SetTimeout(code string) error {
	e, err := newEntityByPayCode(code)
	if err != nil {
		return err
	}
	defer func() { _ = e.release() }()

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
