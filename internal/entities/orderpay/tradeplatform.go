package orderpay

type TradeStatus string

const (
	TradeStatusSuccess  TradeStatus = "支付成功"
	TradeStatusNotPay   TradeStatus = "未支付"
	TradeStatusClosed   TradeStatus = "交易关闭"
	TradeStatusFinished TradeStatus = "交易完成"
)

type TradePay struct {
	AppId string
	MchId string
	Data  interface{}
}

type TradePayQuery struct {
	BuyerLogonId  string
	TransactionId string
	TradeStatus   TradeStatus
}

type Trade interface {
	Pay(e *Entity) (TradePay, error)
	QueryPayStatus(e *Entity) (TradePayQuery, error)
	TradeClose(e *Entity) error
	Refund(e *Entity) error
	QueryRefundStatus(e *Entity)
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
