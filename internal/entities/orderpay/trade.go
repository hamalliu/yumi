package orderpay

import (
	"yumi/external/pay/ali_pagepay"
)

type TradeStatus string

const (
	TradeStatusSuccess  TradeStatus = "支付成功"
	TradeStatusNotPay   TradeStatus = "未支付"
	TradeStatusClosed   TradeStatus = "交易关闭"
	TradeStatusFinished TradeStatus = "交易完成"
)

type Trade interface {
	Pay(e *Entity) ([]byte, error)
	QueryPayStatus(mch ali_pagepay.Merchant, e *Entity) (TradeStatus, error)
	TradeClose(mch ali_pagepay.Merchant, e *Entity) error
	PayNotify(mch ali_pagepay.Merchant, rawQuery string)
	Refund()
	QueryRefundStatus()
	RefundNotify()
}

type TradeWay string

var trades map[TradeWay]Trade

func RegisterTrade(way TradeWay, trade Trade) {
	if trades == nil {
		trades = make(map[TradeWay]Trade)
	}
	trades[way] = trade
}

type Merchant interface {
	BuildOutTradeNo()
}
