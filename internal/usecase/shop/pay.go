package shop

import (
	"time"

	"yumi/internal/entities/orderpay"
)

//提交订单
func SubmitOrder(sellerKey string, notifyUrl string, totalFee int, body, detail, accountGuid, code string, timeoutExpress time.Time) error {
	return orderpay.SubmitOrder(sellerKey, notifyUrl, totalFee, body, detail, accountGuid, code, timeoutExpress)
}

//立即支付
func Pay(code string, tradeWay orderpay.TradeWay) (interface{}, error) {
	return orderpay.Pay(code, tradeWay)
}

//支付遇到问题
func PayProblem(code string) error {
	//查询支付是否支付成功，如果未支付，则关闭订单
	status, err := orderpay.PaySuccess(code)
	if err != nil {
		return err
	}
	if status == orderpay.TradeStatusNotPay {
		return orderpay.CloseTrade(code)
	}

	return nil
}

//支付完成
func PayCompleted(code string) error {
	//查询支付是否支付成功，如果未支付，则关闭订单
	status, err := orderpay.PaySuccess(code)
	if err != nil {
		return err
	}
	if status == orderpay.TradeStatusNotPay {
		return orderpay.CloseTrade(code)
	}

	return nil
}

//取消订单
func CancellOrder(code string) error {
	return orderpay.CancelleOrder(code)
}

//提交退款
func Refund(goodsCode string) error {
	//TODO

	return nil
}

//查询退款
func QueryRefund(code string) error {

	return nil
}
