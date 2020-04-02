package trade

import (
	"fmt"
	"sync"
	"time"

	"yumi/external/pay"
	"yumi/utils/internal_error"
)

var orderSync map[string]*sync.Mutex

type Entity struct {
	dataOp DataOrderPay
	op     OrderPay

	dataOr DataOrderRefund
	or     OrderRefund
}

//======================================================================================================================
/**
 * 业务对象对象基本操作
 * 供业务对象接口调用，对外不开放
 */
//加载支付订单数据
func (e *Entity) loadOrderPay(code string) error {
	if op, err := NewDataOrderPay().Load(code); err != nil {
		return err
	} else {
		e.op = op
	}

	if orderSync[e.op.Code] == nil {
		orderSync[e.op.Code] = new(sync.Mutex)
	}
	orderSync[e.op.Code].Lock()
	return nil
}

//释放支付订单数据
func (e *Entity) releaseOrderPay() error {
	if orderSync[e.op.Code] == nil {
		err := fmt.Errorf("无法释放锁，可能造成死锁")
		return internal_error.With(err)
	}
	orderSync[e.op.Code].Unlock()
	return nil
}

//根据开发者appId和商户订单号加载支付订单数据
func (e *Entity) loadOrderPayByOutTradeNo(appId, outTradeNo string) error {
	if op, err := NewDataOrderPay().LoadByOutTradeNo(appId, outTradeNo); err != nil {
		return err
	} else {
		e.op = op
	}

	if orderSync[e.op.Code] == nil {
		orderSync[e.op.Code] = new(sync.Mutex)
	}
	orderSync[e.op.Code].Lock()
	return nil
}

//加载退款订单数据
func (e *Entity) loadOrderRefund(code string) error {
	if or, err := NewDataOrderRefund().Load(code); err != nil {
		return err
	} else {
		e.or = or
	}
	if op, err := NewDataOrderPay().Load(code); err != nil {
		return err
	} else {
		e.op = op
	}

	if orderSync[e.or.OrderPayCode] == nil {
		orderSync[e.or.OrderPayCode] = new(sync.Mutex)
	}
	orderSync[e.or.OrderPayCode].Lock()
	return nil
}

//释放退款订单数据
func (e *Entity) releaseOrderRefund() error {
	if orderSync[e.or.OrderPayCode] == nil {
		err := fmt.Errorf("无法释放锁，可能造成死锁")
		return internal_error.With(err)
	}
	orderSync[e.or.OrderPayCode].Unlock()
	return nil
}

//根据开发者appId和商户订单号加载退款订单数据
func (e *Entity) loadOrderRefundByOutTradeNo(appId, outTradeNo string) error {
	if or, err := NewDataOrderRefund().LoadByOutTradeNo(appId, outTradeNo); err != nil {
		return err
	} else {
		e.or = or
	}
	if op, err := NewDataOrderPay().LoadByOutTradeNo(appId, outTradeNo); err != nil {
		return err
	} else {
		e.op = op
	}

	if orderSync[e.or.OrderPayCode] == nil {
		orderSync[e.or.OrderPayCode] = new(sync.Mutex)
	}
	orderSync[e.or.OrderPayCode].Lock()
	return nil
}

func getOutTradeNo() string {
	return fmt.Sprintf("%s%s", time.Now().Format("060102150405"), pay.CreateRandomStr(10, pay.NUMBER))
}

func getCode(seqId int64) string {
	//TODO
	return fmt.Sprintf("%s%d", time.Now().Format("060102"), seqId)
}
