package shop

import (
	"time"

	"yumi/internal/entities/trade"
)

//提交订单
func SubmitOrderPay(accountGuid, sellerKey string, notifyUrl string, totalFee int, body, detail string,
	timeoutExpress time.Time) error {
	return trade.SubmitOrderPay(accountGuid, sellerKey, notifyUrl, totalFee, body, detail, timeoutExpress)
}

//立即支付
func Pay(code string, tradeWay trade.TradeWay) (interface{}, error) {
	return trade.Pay(code, tradeWay)
}

//支付遇到问题
func PayProblem(code string) error {
	//查询支付是否支付成功，如果未支付，则关闭订单
	status, err := trade.PaySuccess(code)
	if err != nil {
		return err
	}
	if status == trade.NotPay {
		return trade.CloseTrade(code)
	}

	return nil
}

//支付完成
func PayCompleted(code string) error {
	//查询支付是否支付成功，如果未支付，则关闭订单
	status, err := trade.PaySuccess(code)
	if err != nil {
		return err
	}
	if status == trade.NotPay {
		return trade.CloseTrade(code)
	}

	return nil
}

//取消订单
func CancelleOrderPay(code string) error {
	return trade.CancelleOrderPay(code)
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
