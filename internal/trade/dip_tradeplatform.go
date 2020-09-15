package trade

import "net/http"

//StatusTradePlatform ...
type StatusTradePlatform string

const (
	//StatusTradePlatformSuccess 成功
	StatusTradePlatformSuccess StatusTradePlatform = "SUCCESS"
	//StatusTradePlatformClosed 关闭
	StatusTradePlatformClosed StatusTradePlatform = "CLOSED"
	//StatusTradePlatformFinished 完成
	StatusTradePlatformFinished StatusTradePlatform = "FINISHED"
	//StatusTradePlatformError 错误（异常）
	StatusTradePlatformError StatusTradePlatform = "ERROR"

	//StatusTradePlatformNotPay 未支付
	StatusTradePlatformNotPay StatusTradePlatform = "NOTPAY"
	//StatusTradePlatformRefundProcessing 退款处理中
	StatusTradePlatformRefundProcessing StatusTradePlatform = "REFUNDPROCESSING"
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
	TradeStatus   StatusTradePlatform
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
	TradeStatus   StatusTradePlatform
}

//ReturnRefundNotify ...
type ReturnRefundNotify struct {
	OrderRefundCode string
	ReqData         interface{}
}

//Trade ...
type Trade interface {
	// 发起支付
	Pay(op OrderPay) (ReturnPay, error)
	// 支付通知提供三个接口，以应对不同支付平台的接口差异
	// 处理请求
	PayNotifyReq(req *http.Request) (ReturnPayNotify, error)
	// 检查参数
	PayNotifyCheck(op OrderPay, reqData interface{}) error
	// 应答
	PayNotifyResp(err error, resp http.ResponseWriter)
	// 查询支付状态
	QueryPayStatus(op OrderPay) (ReturnQueryPay, error)
	// 关闭交易
	TradeClose(op OrderPay) error
	// 退款
	Refund(op OrderPay, or OrderRefund) error
	// 查询退款状态
	QueryRefundStatus(op OrderPay, or OrderRefund) (ReturnQueryRefund, error)
	// 退款通知提供三个接口，以应对不同支付平台的接口差异
	// 处理请求
	RefundNotifyReq(req *http.Request) (ReturnRefundNotify, error)
	// 检查参数
	RefundNotifyCheck(op OrderPay, or OrderRefund, reqData interface{}) error
	// 应答
	RefundNotifyResp(err error, resp http.ResponseWriter)
}

//Way ...
type Way string

var trades map[Way]Trade

//RegisterTrade 注册交易方式
func RegisterTrade(way Way, trade Trade) {
	if trades == nil {
		trades = make(map[Way]Trade)
	}
	trades[way] = trade
}

func getTrade(tradeWay Way) Trade {
	return trades[tradeWay]
}

//Merchant 商户应该实现的接口
type Merchant interface {
	BuildOutTradeNo()
}
