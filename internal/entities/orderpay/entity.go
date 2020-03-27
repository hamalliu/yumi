package orderpay

import (
	"fmt"
	"sync"
	"time"

	"yumi/external/pay"
	"yumi/utils/internal_error"
)

var orderSync map[string]*sync.Mutex

type Entity struct {
	Data Data
	OrderPay
}

//======================================================================================================================
/**
 * 业务对象对象基本操作
 * 供业务对象接口调用，对外不开放
 */

//提交订单
func (e *Entity) submitOrder(notifyUrl string, totalFee int, body, detail, accountGuid, code string, timeoutExpress time.Time) error {
	return e.Data.SubmitOrder(getOutTradeNo(), notifyUrl, totalFee, body, detail, timeoutExpress, time.Now(), code, Submitted)
}

//支付成功，更新订单状态（待支付->已支付）
func (e *Entity) paySuccess(paytime string) error {
	return e.Data.PaySuccess(time.Now(), Paid)
}

//关闭订单，更新订单状态（待支付->已提交）
func (e *Entity) setSubmitted() error {
	return e.Data.SetSubmitted(Submitted)
}

//设置订单错误
func (e *Entity) setError() error {
	return e.Data.SetError(time.Now(), Error)
}

//设置取消订单
func (e *Entity) setCancelled() error {
	return e.Data.SetCancelled(time.Now(), Cancelled)
}

//设置支付方式
func (e *Entity) setPayWay(payWay TradeWay, appId, mchId string) error {
	return e.Data.SetPayWay(payWay, appId, mchId, WaitPay)
}

//设置订单号
func (e *Entity) setTransactionId(transactionId, buyerLogonId string) error {
	return e.Data.SetTransactionId(transactionId, buyerLogonId)
}

//加载订单数据
func (e *Entity) load(code string) error {
	if op, err := NewData().Load(code); err != nil {
		return err
	} else {
		e.OrderPay = op
	}

	if orderSync[e.Code] == nil {
		orderSync[e.Code] = new(sync.Mutex)
	}
	orderSync[e.Code].Lock()
	return nil
}

//释放数据
func (e *Entity) release() error {
	if orderSync[e.Code] == nil {
		err := fmt.Errorf("无法释放锁，可能造成死锁")
		return internal_error.With(err)
	}
	orderSync[e.Code].Unlock()
	return nil
}

//根据开发者appId和商户订单号加载订单数据
func (e *Entity) loadByOutTradeNo(appId, outTradeNo string) error {
	if op, err := NewData().LoadByOutTradeNo(appId, outTradeNo); err != nil {
		return err
	} else {
		e.OrderPay = op
	}

	if orderSync[e.Code] == nil {
		orderSync[e.Code] = new(sync.Mutex)
	}
	orderSync[e.Code].Lock()
	return nil
}

func getOutTradeNo() string {
	return fmt.Sprintf("%s%s", time.Now().Format("060102150405"), pay.CreateRandomStr(10, pay.NUMBER))
}

func getCode(seqId int64) string {
	//TODO
	return fmt.Sprintf("%s%d", time.Now().Format("060102"), seqId)
}
