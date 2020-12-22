package db

import "yumi/usecase/trade"

// Init trade db
func Init() {
	trade.InitDataOrderPay(&OrderPay{})
	trade.InitDataOrderRefund(&OrderRefund{})
}
