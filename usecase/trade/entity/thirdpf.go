package entity

import (
	"net/http"
)

//StatusTradePlatform ...
type StatusTradePlatform string

const (
	//StatusTradePlatformSuccess 成功
	StatusTradePlatformSuccess StatusTradePlatform = "SUCCESS"
	//StatusTradePlatformClosed 关闭
	StatusTradePlatformClosed StatusTradePlatform = "CLOSED"
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
	Pay(op OrderPayAttribute) (ReturnPay, error)
	// 支付通知提供三个接口，以应对不同支付平台的接口差异
	// 处理请求
	PayNotifyReq(req *http.Request) (ReturnPayNotify, error)
	// 检查参数
	PayNotifyCheck(op OrderPayAttribute, reqData interface{}) error
	// 应答
	PayNotifyResp(err error, resp http.ResponseWriter)
	// 查询支付状态
	QueryPayStatus(op OrderPayAttribute) (ReturnQueryPay, error)
	// 关闭交易
	TradeClose(op OrderPayAttribute) error

	// 退款
	Refund(op OrderPayAttribute, or OrderRefundAttribute) error
	// 查询退款状态
	QueryRefundStatus(op OrderPayAttribute, or OrderRefundAttribute) (ReturnQueryRefund, error)
	// 退款通知提供三个接口，以应对不同支付平台的接口差异
	// 处理请求
	RefundNotifyReq(req *http.Request) (ReturnRefundNotify, error)
	// 检查参数
	RefundNotifyCheck(op OrderPayAttribute, or OrderRefundAttribute, reqData interface{}) error
	// 应答
	RefundNotifyResp(err error, resp http.ResponseWriter)
}
