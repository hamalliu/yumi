package trade

type Status string

const (
	Success  Status = "支付成功"
	NotPay   Status = "未支付"
	Closed   Status = "交易关闭"
	Finished Status = "交易完成"
)

type TradePay struct {
	AppId string
	MchId string
	Data  interface{}
}

type TradePayQuery struct {
	BuyerLogonId  string
	TransactionId string
	TradeStatus   Status
}

type Trade interface {
	Pay(op OrderPay) (TradePay, error)
	QueryPayStatus(op OrderPay) (TradePayQuery, error)
	TradeClose(op OrderPay) error
	Refund(op OrderPay, or OrderRefund) error
	QueryRefundStatus(op OrderPay, or OrderRefund)
}

type TradeWay string

var trades map[TradeWay]Trade

func RegisterTrade(way TradeWay, trade Trade) {
	if trades == nil {
		trades = make(map[TradeWay]Trade)
	}
	trades[way] = trade
}

func getTrade(tradeWay TradeWay) Trade {
	return trades[tradeWay]
}

type Merchant interface {
	BuildOutTradeNo()
}
