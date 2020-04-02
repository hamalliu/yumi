package dbentities

import (
	"time"

	"yumi/internal/entities/trade"
)

//退款订单
type OrderRefund struct {
	SeqId int64 `db:"seq_id"`
	trade.OrderRefund
}

func (m OrderRefund) Clone() trade.DataOrderRefund {
	return &OrderRefund{}
}

//提交订单
func (m *OrderRefund) Submit() error {
	return nil
}

//关闭订单，更新订单状态（待支付->已提交）
func (m *OrderRefund) SetSubmitted(status trade.OrderStatus) error {
	return nil
}

//设置订单错误
func (m *OrderRefund) SetRefunding(refundingTime time.Time, status trade.OrderStatus) error {
	return nil
}

//设置取消订单
func (m *OrderRefund) SetCancelled(cancelTime time.Time, status trade.OrderStatus) error {
	return nil
}

//设置订单错误
func (m *OrderRefund) SetRefunded(errorTime time.Time, status trade.OrderStatus) error {
	return nil
}

//设置订单错误
func (m *OrderRefund) SetError(errorTime time.Time, status trade.OrderStatus) error {
	return nil
}

//加载订单数据
func (m *OrderRefund) Load(code string) (trade.OrderRefund, error) {
	return trade.OrderRefund{}, nil
}

//根据开发者appId和商户订单号加载订单数据
func (m *OrderRefund) LoadByOutTradeNo(appId, outTradeNo string) (trade.OrderRefund, error) {
	return trade.OrderRefund{}, nil
}
