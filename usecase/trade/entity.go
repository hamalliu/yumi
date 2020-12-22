package trade

import (
	"sync"
	
	"yumi/pkg/ecode"
)

var (
	orderPaySync    map[string]*sync.Mutex
	orderRefundSync map[string]*sync.Mutex
)

//Entity 业务对象
type Entity struct {
	// 支付订单业务对象
	op OrderPay
	// 支付订单数据接口
	dataOp DataOrderPay

	// 退款订单业务对象
	or OrderRefund
	// 退款订单数据接口
	dataOr DataOrderRefund
}

func init() {
	if orderPaySync == nil {
		orderPaySync = make(map[string]*sync.Mutex)
	}
	if orderRefundSync == nil {
		orderRefundSync = make(map[string]*sync.Mutex)
	}
}

var entityMutex sync.Mutex

/**
 * 业务对象对象基本操作
 * 供业务对象接口调用，对外不开放
 */

//emptyEntity 空的业务对象，主要用于新的订单
func emptyEntity() (*Entity, error) {
	var err error
	e := Entity{}

	e.dataOp, err = newDataOrderPay("")
	if err != nil {
		return nil, err
	}
	e.op = e.dataOp.Entity()
	e.dataOr, err = newDataOrderRefund("")
	e.or = e.dataOr.Entity()
	return &e, nil
}

//newEntityByPayCode 加载支付订单数据
func newEntityByPayCode(code string) (*Entity, error) {
	if code == "" {
		return nil, ecode.OrdernoDoesNotExist
	}

	entityMutex.Lock()
	if orderPaySync[code] == nil {
		orderPaySync[code] = new(sync.Mutex)
	}
	entityMutex.Unlock()
	orderPaySync[code].Lock()

	var err error
	e := Entity{}
	e.dataOp, err = newDataOrderPay(code)
	if err != nil {
		orderPaySync[code].Unlock()
		return nil, err
	}
	e.op = e.dataOp.Entity()
	e.dataOr, err = newDataOrderRefund("")
	e.or = e.dataOr.Entity()

	return &e, nil
}

//newEntityByRefundCode 加载退款订单数据
func newEntityByRefundCode(code string) (*Entity, error) {
	if code == "" {
		return nil, ecode.OrdernoDoesNotExist
	}

	entityMutex.Lock()
	if orderRefundSync[code] == nil {
		orderRefundSync[code] = new(sync.Mutex)
	}
	entityMutex.Unlock()
	orderRefundSync[code].Lock()

	// 加载退款订单
	var err error
	e := Entity{}
	e.dataOr, err = newDataOrderRefund(code)
	if err != nil {
		orderRefundSync[code].Unlock()
		return nil, err
	}
	e.or = e.dataOr.Entity()

	// 需要加载支付订单，因为有关验证需要支付相关数据
	if orderPaySync[e.or.OrderPayCode] == nil {
		orderPaySync[e.or.OrderPayCode] = new(sync.Mutex)
	}
	orderPaySync[e.or.OrderPayCode].Lock()
	e.dataOp, err = newDataOrderPay(e.or.OrderPayCode)
	if err != nil {
		orderPaySync[e.or.OrderPayCode].Unlock()
		orderRefundSync[code].Unlock()
		return nil, err
	}
	e.op = e.dataOp.Entity()

	return &e, nil
}

//release 释放订单数据
func (e *Entity) release() error {
	if e == nil {
		return nil
	}

	if e.op.Code != "" && orderPaySync[e.op.Code] != nil {
		orderPaySync[e.op.Code].Unlock()
	}
	if e.or.OrderPayCode != "" && orderPaySync[e.or.OrderPayCode] != nil {
		orderPaySync[e.or.OrderPayCode].Unlock()
	}
	if e.or.Code != "" && orderRefundSync[e.or.Code] != nil {
		orderRefundSync[e.or.Code].Unlock()
	}

	return nil
}
