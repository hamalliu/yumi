package trade

import "net/http"

//Status ...
type Status string

const (
	//Success ...
	Success  Status = "SUCCESS"  //成功
	//Closed ...
	Closed   Status = "CLOSED"   //关闭
	//Finished ...
	Finished Status = "FINISHED" //完成
	//ERROR ...
	ERROR    Status = "ERROR"    //错误（异常）

	//NotPay ...
	NotPay           Status = "NOTPAY"           //未支付
	//RefundProcessing ...
	RefundProcessing Status = "REFUNDPROCESSING" //退款处理中
)

//ReturnPay ...
type ReturnPay struct {
	AppID string
	MchID string
	Data  string
}

//ReturnQueryPay ...
type ReturnQueryPay struct {
	BuyerLogonID  string
	TransactionID string
	TradeStatus   Status
}

//ReturnPayNotify ...
type ReturnPayNotify struct {
	OrderPayCode string
	ReqData      interface{}
}

//ReturnQueryRefund ...
type ReturnQueryRefund struct {
	RefundID      string
	RefundLogonID string
	TradeStatus   Status
}

//ReturnRefundNotify ...
type ReturnRefundNotify struct {
	OrderRefundCode string
	ReqData         interface{}
}

//Trade ...
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

//Way ...
type Way string

var trades map[Way]Trade

//RegisterTrade ...
func RegisterTrade(way Way, trade Trade) {
	if trades == nil {
		trades = make(map[Way]Trade)
	}
	trades[way] = trade
}

func getTrade(tradeWay Way) Trade {
	return trades[tradeWay]
}

//Merchant ...
type Merchant interface {
	BuildOutTradeNo()
}
