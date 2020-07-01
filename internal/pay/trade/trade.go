package trade

import "net/http"

type Status string

const (
	Success  Status = "SUCCESS"  //成功
	Closed   Status = "CLOSED"   //关闭
	Finished Status = "FINISHED" //完成
	ERROR    Status = "ERROR"    //错误（异常）

	NotPay           Status = "NOTPAY"           //未支付
	RefundProcessing Status = "REFUNDPROCESSING" //退款处理中
)

type ReturnPay struct {
	AppId string
	MchId string
	Data  string
}

type ReturnQueryPay struct {
	BuyerLogonId  string
	TransactionId string
	TradeStatus   Status
}

type ReturnPayNotify struct {
	OrderPayCode string
	ReqData      interface{}
}

type ReturnQueryRefund struct {
	RefundId      string
	RefundLogonId string
	TradeStatus   Status
}

type ReturnRefundNotify struct {
	OrderRefundCode string
	ReqData         interface{}
}

type Trade interface {
	Pay(op OrderPay) (ReturnPay, error)
	PayNotifyReq(req *http.Request) (ReturnPayNotify, error)
	PayNotifyCheck(op OrderPay, reqData interface{}) error
	PayNotifyResp(err error, resp http.ResponseWriter)
	QueryPayStatus(op OrderPay) (ReturnQueryPay, error)
	TradeClose(op OrderPay) error
	Refund(op OrderPay, or OrderRefund) error
	QueryRefundStatus(op OrderPay, or OrderRefund) (ReturnQueryRefund, error)
	RefundNotifyReq(req *http.Request) (ReturnRefundNotify, error)
	RefundNotifyCheck(op OrderPay, or OrderRefund, reqData interface{}) error
	RefundNotifyResp(err error, resp http.ResponseWriter)
}

type Way string

var trades map[Way]Trade

func RegisterTrade(way Way, trade Trade) {
	if trades == nil {
		trades = make(map[Way]Trade)
	}
	trades[way] = trade
}

func getTrade(tradeWay Way) Trade {
	return trades[tradeWay]
}

type Merchant interface {
	BuildOutTradeNo()
}
