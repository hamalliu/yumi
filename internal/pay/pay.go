package pay

import (
	"fmt"
	"net/http"
	"time"

	"yumi/internal/pay/db"
	"yumi/internal/pay/trade"
)

//提交订单
func SubmitOrderPay(accountGuid, sellerKey string, totalFee int, body, detail string,
	timeoutExpress time.Time) (string, error) {
	if totalFee <= 0 {
		return "", fmt.Errorf("金额必须大于0")
	}
	return trade.SubmitOrderPay(accountGuid, sellerKey, totalFee, body, detail, timeoutExpress)
}

//立即支付
func Pay(code string, tradeWay trade.Way, clientIp, notifyUrl string, payExpire time.Time) (string, error) {
	return trade.Pay(code, tradeWay, clientIp, notifyUrl, payExpire)
}

//轮循查询支付结果(发货的时候必须处理重复的成功通知)
func RoundRobinQueryPaySuccess(code string) (bool, error) {
	status, err := trade.PaySuccess(code)
	if err != nil {
		return false, err
	}
	if status == trade.Success {
		return true, nil
	}

	return false, nil
}

//支付遇到问题
func PayProblem(code string) (bool, error) {
	//查询支付是否支付成功，如果未支付，则关闭订单
	status, err := trade.PaySuccess(code)
	if err != nil {
		return false, err
	}
	if status == trade.NotPay {
		err := trade.CloseTrade(code)
		if err != nil {
			return false, err
		}
	}

	if status == trade.Success {
		return true, nil
	}

	return false, nil
}

//支付完成
func PayCompleted(code string) (bool, error) {
	//查询支付是否支付成功，如果未支付，则关闭订单
	status, err := trade.PaySuccess(code)
	if err != nil {
		return false, err
	}
	if status == trade.NotPay {
		err := trade.CloseTrade(code)
		if err != nil {
			return false, err
		}
	} else if status == trade.Success {
		return true, nil
	}

	return false, nil
}

//取消订单
func CancelOrderPay(code string) error {
	return trade.CancelOrderPay(code)
}

//支付通知(发货的时候必须处理重复的成功通知)
func PayNotify(tradeWay trade.Way, resp http.ResponseWriter, req *http.Request) (string, bool) {
	orderPayCode, tradeStatus := trade.PayNotify(tradeWay, resp, req)
	if tradeStatus == trade.Success {
		return orderPayCode, true
	}

	return orderPayCode, false
}

//提交退款
func Refund(orderPayCode, notifyUrl string, refundAccountGuid string, refundFee int, refundDesc string, submitTime,
	timeoutExpress time.Time) (string, error) {
	return trade.Refund(orderPayCode, notifyUrl, refundAccountGuid, refundFee, refundDesc, submitTime, timeoutExpress)
}

//查询退款
func RefundSuccess(code string) (bool, error) {
	status, err := trade.RefundSuccess(code)
	if err != nil {
		return false, err
	}
	if status == trade.Success {
		return true, nil
	}
	return false, nil
}

//退款通知
func RefundNotify(tradeWay trade.Way, resp http.ResponseWriter, req *http.Request) bool {
	tradeStatus := trade.RefundNotify(tradeWay, resp, req)
	if tradeStatus == trade.Success {
		return true
	}

	return false
}

//纠正异常订单
func CorrectPaySuccess() ([]string, error) {
	codes, err := db.GetOrderPayCodesSubmittedAndWaitPay()
	if err != nil {
		return nil, err
	}

	shipCodes := []string{}
	for _, code := range codes {
		status, err := trade.PaySuccess(code)
		if err != nil {
			return nil, err
		}
		if status == trade.NotPay {
			err := trade.CloseTrade(code)
			if err != nil {
				return nil, err
			}
		}

		if status == trade.Success {
			shipCodes = append(shipCodes, code)
		} else {
			if err := trade.SetTimeout(code); err != nil {
				return nil, err
			}
		}
	}
	return shipCodes, nil
}
