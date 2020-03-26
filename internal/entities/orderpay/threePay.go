package orderpay

type ThreePay interface {
}

type Merchant interface {
	BuildOutTradeNo()
}
