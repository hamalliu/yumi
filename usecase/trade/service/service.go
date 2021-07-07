package service

import (
	"net/http"

	"yumi/pkg/status"
	"yumi/usecase/trade/entity"
	"yumi/usecase/trade/thirdpf"
)

// Service ...
type Service struct {
	data   Data
	trades thirdpf.Trades
}

// New a Service object
func New(data Data, trades *thirdpf.Trades) (*Service, error) {
	return &Service{}, nil
}

// CreateOrderPay 提交支付订单
func (s *Service) CreateOrderPay(req CreateOrderPayRequest) (resp CreateOrderPayResponse, err error) {
	attr := req.Attribute()
	trade, err := s.trades.GetTrade(thirdpf.Way(attr.TradeWay))
	if err != nil {
		return
	}
	op := entity.NewOrderPay(&attr, trade)
	err = op.Submit()
	if err != nil {
		return
	}

	// 持久化
	err = s.data.CreateOrderPay(attr)
	if err != nil {
		return
	}

	resp.set(attr)
	return
}

// CancelOrderPay 取消支付订单
func (s *Service) CancelOrderPay(code string) (err error) {
	attr, err := s.data.GetOrderPay(code)
	if err != nil {
		return err
	}
	trade, err := s.trades.GetTrade(thirdpf.Way(attr.TradeWay))
	if err != nil {
		return
	}
	op := entity.NewOrderPay(&attr, trade)
	err = op.Cancel()
	if err != nil {
		return err
	}

	err = s.data.UpdateOrderPay(attr)
	if err != nil {
		return err
	}

	return
}

// Pay 发起支付
func (s *Service) Pay(req PayRequest) (PayResponse, error) {
	resp := PayResponse{}

	attr, err := s.data.GetOrderPay(req.Code)
	if err != nil {
		return resp, err
	}
	trade, err := s.trades.GetTrade(thirdpf.Way(attr.TradeWay))
	if err != nil {
		return resp, err
	}
	op := entity.NewOrderPay(&attr, trade)
	curTrade, err := s.trades.GetTrade(thirdpf.Way(attr.TradeWay))
	if err != nil {
		return resp, err
	}
	resp.Data, err = op.Pay(curTrade, req.TradeWay, req.ClientIP, req.NotifyURL)
	if err != nil {
		return resp, err
	}

	err = s.data.UpdateOrderPay(attr)
	if err != nil {
		return resp, err
	}

	resp.set(attr)
	return resp, nil
}

// QueryPaid 查询支付成功（只查询待支付订单）
func (s *Service) QueryPaid(code string) (paid bool, err error) {
	attr, err := s.data.GetOrderPay(code)
	if err != nil {
		return
	}
	trade, err := s.trades.GetTrade(thirdpf.Way(attr.TradeWay))
	if err != nil {
		return
	}
	op := entity.NewOrderPay(&attr, trade)
	paid, err = op.QueryPaid()
	if err != nil {
		return
	}

	err = s.data.UpdateOrderPay(attr)
	if err != nil {
		return
	}

	return
}

// PayNotify 支付通知(待支付时处理通知)
func (s *Service) PayNotify(tradeWay string, w http.ResponseWriter, req *http.Request) error {
	curTrade, err := s.trades.GetTrade(thirdpf.Way(tradeWay))
	if err != nil {
		curTrade.PayNotifyResp(err, w)
		return err
	}

	//解析通知参数
	ret, err := curTrade.PayNotifyReq(req)
	if err != nil {
		curTrade.PayNotifyResp(err, w)
		return err
	}

	attr, err := s.data.GetOrderPay(ret.OrderPayCode)
	if err != nil {
		curTrade.PayNotifyResp(err, w)
		return err
	}
	//检查通知参数
	if err := curTrade.PayNotifyCheck(attr, ret.ReqData); err != nil {
		curTrade.PayNotifyResp(err, w)
		return err
	}

	// 处理notify
	trade, err := s.trades.GetTrade(thirdpf.Way(attr.TradeWay))
	if err != nil {
		return err
	}
	op := entity.NewOrderPay(&attr, trade)
	paid, err := op.QueryPaid()
	if err != nil {
		curTrade.PayNotifyResp(err, w)
		return err
	}
	if !paid {
		err = status.InvalidRequest()
		curTrade.PayNotifyResp(err, w)
		return err
	}

	// 持久化
	err = s.data.UpdateOrderPay(attr)
	if err != nil {
		curTrade.PayNotifyResp(err, w)
		return err
	}

	curTrade.PayNotifyResp(nil, w)
	return nil
}

// // Refund 提交退款订单（只退款已支付订单）
// func Refund(orderPayCode, notifyURL string, refundAccountGUID string, refundFee int, refundDesc string, submitTime, timeoutExpress time.Time) (string, error) {
// 	e, err := newEntityByPayCode(orderPayCode)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer func() { _ = e.release() }()

// 	switch e.op.Status {
// 	//支付订单必须为已支付
// 	case Paid:
// 		//是否存在已提交或退款中的退款订单
// 		if exist, err := e.dataOr.ExistRefundingOrSubmitted(e.op.Code); err != nil {
// 			return "", err
// 		} else if exist {
// 			return "", ecode.CannotRepeatSendRefund
// 		}

// 		//退款金额是否超额
// 		count, countFee := 0, 0
// 		if count, countFee, err = e.dataOr.GetRefundFee(e.op.Code); err != nil {
// 			return "", err
// 		} else if countFee+refundFee > e.op.TotalFee {
// 			err = fmt.Errorf("退款总额超过支付金额")
// 			return "", ecode.ServerErr(err)
// 		}

// 		// 获取交易方式
// 		trade := getTrade(e.op.TradeWay)
// 		if trade == nil {
// 			_ = e.dataOr.SetError(time.Now(), "tradeway（交易方式）不存在", Error)
// 			err = fmt.Errorf("refundcode: %s，tradeway（交易方式）不存在", e.or.Code)
// 			return "", ecode.ServerErr(err)
// 		}

// 		// 提交退款订单
// 		code := getCode(OrderRefundCode)
// 		outRefundNo := code
// 		if err := e.dataOr.Submit(code, orderPayCode, count, notifyURL, refundAccountGUID, e.op.TradeWay, outRefundNo, refundFee,
// 			refundDesc, submitTime, timeoutExpress, Submitted); err != nil {
// 			return "", err
// 		}
// 		e.dataOr, err = newDataOrderRefund(code)
// 		e.or = e.dataOr.Entity()

// 		// 退款
// 		if err := trade.Refund(e.op, e.or); err != nil {
// 			return "", err
// 		}

// 		// 设置退款中
// 		if err := e.dataOr.SetRefunding(Refunding); err != nil {
// 			return "", err
// 		}
// 		return code, nil
// 	default:
// 		return "", ecode.InvalidSendRefund
// 	}
// }

// // RefundSuccess 退款查询（只查询退款中的订单）
// func RefundSuccess(code string) (res Status, err error) {
// 	e, err := newEntityByRefundCode(code)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer func() { _ = e.release() }()

// 	switch e.or.Status {
// 	case Refunding:
// 		trade := getTrade(e.op.TradeWay)
// 		if trade == nil {
// 			_ = e.dataOr.SetError(time.Now(), "tradeway（交易方式）不存在", Error)
// 			err = fmt.Errorf("refundcode: %s，tradeway（交易方式）不存在", e.or.Code)
// 			return "", ecode.ServerErr(err)
// 		}

// 		ret, err := trade.QueryRefundStatus(e.op, e.or)
// 		if err != nil {
// 			return "", err
// 		}
// 		if ret.TradeStatus == StatusTradePlatformSuccess {
// 			if err := e.dataOr.SetRefunded(ret.RefundID, time.Now(), Refunded); err != nil {
// 				log.Error(err)
// 				return "", err
// 			}
// 			return Refunded, nil
// 		}
// 		return Refunding, nil
// 	case Refunded:
// 		return Refunded, nil
// 	default:
// 		return "", ecode.InvalidQueryRefund
// 	}
// }

// // RefundNotify 退款通知
// func RefundNotify(way Way, resp http.ResponseWriter, req *http.Request) Status {
// 	trade := getTrade(way)
// 	if trade == nil {
// 		err := fmt.Errorf("交易方式不存在：%s", way)
// 		log.Error(err)
// 		trade.RefundNotifyResp(err, resp)
// 		return ""
// 	}

// 	ret, err := trade.RefundNotifyReq(req)
// 	if err != nil {
// 		trade.RefundNotifyResp(err, resp)
// 		return ""
// 	}
// 	e, err := newEntityByRefundCode(ret.OrderRefundCode)
// 	if err != nil {
// 		trade.PayNotifyResp(err, resp)
// 		return ""
// 	}
// 	defer func() { _ = e.release() }()

// 	if err := trade.RefundNotifyCheck(e.op, e.or, ret.ReqData); err != nil {
// 		trade.PayNotifyResp(err, resp)
// 		return ""
// 	}

// 	status, err := RefundSuccess(ret.OrderRefundCode)
// 	if err != nil {
// 		trade.PayNotifyResp(err, resp)
// 		return ""
// 	}
// 	if status == Refunded {
// 		trade.RefundNotifyResp(nil, resp)
// 		return Refunded
// 	}

// 	return status
// }
