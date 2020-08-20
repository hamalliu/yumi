package pay

import (
	"testing"
	"time"

	"yumi/internal/pay/trade"
	"yumi/internal/pay/tradeplatform"
)

func TestSubmitOrderPay(t *testing.T) {
	accountGuid := "liuxin@guid"
	sellerKey := "zzyq_account_001"
	totalFee := 1
	body := "商品描述"
	detail := "商品详情"
	timeoutExpress := time.Now().Add(time.Minute * 30)

	if code, err := SubmitOrderPay(accountGuid, sellerKey, totalFee, body, detail, timeoutExpress); err != nil {
		t.Error(err)
	} else {
		t.Log(code)
	}
}

func TestPay(t *testing.T) {
	code := "2051114151401906"
	notifyUrl := "http:120.24.183.196:20192/signin"

	t.Log(code)
	if intf, err := Pay(code, tradeplatform.WxPayNATIVE2, "125.71.211.25", notifyUrl, time.Now().Add(time.Minute*30)); err != nil {
		t.Error(err)
	} else {
		t.Log(intf)
	}

	status, err := trade.PayStatus(code)
	if err != nil {
		t.Error(err)
	}

	t.Log(status)
}

func TestPayProblem(t *testing.T) {

	code := "205616102203701156"

	if _, err := PayProblem(code); err != nil {
		t.Error(err)
		return
	}

	status, err := trade.PayStatus(code)
	if err != nil {
		t.Error(err)
	}

	t.Log(status)
}

func TestPayCompleted(t *testing.T) {
	code := "20424172163301438"

	if _, err := PayCompleted(code); err != nil {
		t.Error(err)
		return
	}

	status, err := trade.PayStatus(code)
	if err != nil {
		t.Error(err)
	}

	t.Log(status)
}

func TestCancelOrderPay(t *testing.T) {
	code := "20424172163301438"

	if err := CancelOrderPay(code); err != nil {
		t.Error(err)
		return
	}

	status, err := trade.PayStatus(code)
	if err != nil {
		t.Error(err)
	}

	t.Log(status)
}

func TestPayNotify(t *testing.T) {
	//TODO
}

func TestRefund(t *testing.T) {
	orderPayCode := "205616102203701156"
	notifyUrl := "https://weixin.qq.com/notify/"
	refundAccountGuid := "liuxin@guid"
	refundFee := 1
	refundDesc := "测试退款"
	now := time.Now()
	code, err := Refund(orderPayCode, notifyUrl, refundAccountGuid, refundFee, refundDesc, now, now.Add(time.Minute*30))
	if err != nil {
		time.Sleep(time.Second * 1)
		t.Error(code, err)
	}

	t.Log(code)
}

func TestRefundSuccess(t *testing.T) {
	code := "205616122788111101"
	success, err := RefundSuccess(code)
	if err != nil {
		time.Sleep(time.Second * 1)
		t.Error(err)
	}
	if success {
		t.Log("success")
	}
}
