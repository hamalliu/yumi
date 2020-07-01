package trade

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"yumi/pkg/log"
	"yumi/utils"
)

var (
	orderPaySync    map[string]*sync.Mutex
	orderRefundSync map[string]*sync.Mutex
)

type Entity struct {
	dataOp DataOrderPay
	op     OrderPay

	dataOr DataOrderRefund
	or     OrderRefund
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
//加载支付订单数据
func NewEntityByPayCode(code string) (*Entity, error) {
	if code == "" {
		dataOp, err := NewDataOrderPay(code)
		if err != nil {
			return nil, err
		}
		e := Entity{dataOp: dataOp}
		e.op = e.dataOp.Data()
		e.dataOr, err = NewDataOrderRefund(code)
		return &e, nil
	}

	entityMutex.Lock()
	if orderPaySync[code] == nil {
		orderPaySync[code] = new(sync.Mutex)
	}
	entityMutex.Unlock()
	orderPaySync[code].Lock()

	dataOp, err := NewDataOrderPay(code)
	if err != nil {
		orderPaySync[code].Unlock()
		return nil, err
	}
	e := Entity{dataOp: dataOp}
	e.op = e.dataOp.Data()

	e.dataOr, err = NewDataOrderRefund("")

	return &e, nil
}

//释放支付订单数据
func (e *Entity) ReleaseOrderPay() error {
	if e.op.Code == "" {
		return nil
	}

	if orderPaySync[e.op.Code] == nil {
		err := fmt.Errorf("无法释放锁，可能造成死锁")
		log.Error(err)
		return err
	}
	orderPaySync[e.op.Code].Unlock()
	return nil
}

//加载退款订单数据
func NewEntityByRefundCode(code string) (*Entity, error) {
	if code == "" {
		panic(fmt.Errorf("code 不能为空"))
	}

	entityMutex.Lock()
	if orderRefundSync[code] == nil {
		orderRefundSync[code] = new(sync.Mutex)
	}
	entityMutex.Unlock()
	orderRefundSync[code].Lock()

	e := Entity{}
	var err error
	e.dataOr, err = NewDataOrderRefund(code)
	if err != nil {
		orderRefundSync[code].Unlock()
		return nil, err
	}
	e.or = e.dataOr.Data()

	if orderPaySync[e.or.OrderPayCode] == nil {
		orderPaySync[e.or.OrderPayCode] = new(sync.Mutex)
	}
	orderPaySync[e.or.OrderPayCode].Lock()
	e.dataOp, err = NewDataOrderPay(e.or.OrderPayCode)
	if err != nil {
		orderPaySync[e.or.OrderPayCode].Unlock()
		orderRefundSync[code].Unlock()
		return nil, err
	}
	e.op = e.dataOp.Data()

	return &e, nil
}

//释放退款订单数据
func (e *Entity) ReleaseOrderRefund() error {
	if e.op.Code == "" {
		return nil
	}

	if orderRefundSync[e.or.Code] == nil {
		err := fmt.Errorf("无法释放锁，可能造成死锁")
		log.Error(orderRefundSync, err)
		return err
	}
	orderRefundSync[e.or.Code].Unlock()

	if orderPaySync[e.or.OrderPayCode] == nil {
		err := fmt.Errorf("无法释放锁，可能造成死锁")
		log.Error(orderPaySync, err)
		return err
	}
	orderPaySync[e.or.OrderPayCode].Unlock()
	return nil
}

func getOutTradeNo() string {
	prefix := strings.ReplaceAll(time.Now().Format("06121545.999999"), ".", "")
	return fmt.Sprintf("%s%s", prefix, utils.CreateRandomStr(10, utils.NUMBER))
}

//生成订单号
type CodeType uint8

const (
	OrderPayCode CodeType = iota
	OrderRefundCode
)

var count uint64

func getCode(codeType CodeType) string {
	prefix := strings.ReplaceAll(time.Now().Format("06121545.999"), ".", "")
	random := utils.CreateRandomStr(3, utils.NUMBER)
	if count >= 100 {
		count = 0
	}
	count++
	return fmt.Sprintf("%s%d%d%s", prefix, codeType, count, random)
}
