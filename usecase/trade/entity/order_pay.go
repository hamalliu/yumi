package entity

import (
	"yumi/pkg/types"
)

//Status 订单状态
type Status string

const (
	//Submitted "已提交" #包含支付订单和退款订单
	Submitted Status = "SUBMITTED"
	//WaitPay "待支付"
	WaitPay Status = "WAIT_PAY"
	//Paid "已支付"
	Paid Status = "PAID"
	//Cancelled "已取消" #包含支付订单和退款订单
	Cancelled Status = "CANCELLED"
	//Refunding "退款中"
	Refunding Status = "REFUNDING"
	//Refunded "已退款"
	Refunded Status = "REFUNDED"
	//Error "错误" #包含支付订单和退款订单
	Error Status = "ERROR"
)

//OrderPayAttribute 支付订单
type OrderPayAttribute struct {
	//唯一编码
	Code string `db:"code"`
	//交易方式
	TradeWay string `db:"trade_way"`
	//卖家key（对本支付系统唯一）
	SellerKey string `db:"seller_key"`
	//开放平台应用唯一id
	AppID string `db:"app_id"`
	//商户id（如果支付方式为支付宝则是seller_id，
	// 收款支付宝账号对应的支付宝唯一用户号。如果为微信则是微信商户号）
	MchID string `db:"mch_id"`
	//支付平台订单号（如果支付方式为支付宝则是支付宝订单号，如果为微信则是微信订单号）
	TransactionID string `db:"transaction_id"`
	//回调url
	NotifyURL string `db:"notify_url"`
	//买家账号id（如果支付方式为支付宝则是买家支付宝账号id，如果为微信则是微信账号id）
	BuyerLogonID string `db:"buyer_logon_id"`
	//用户支付客户端ip
	SpbillCreateIP string `db:"spbill_create_ip"`
	//买家账号guid
	BuyerAccountGUID string `db:"buyer_account_guid"`
	//订单总金额，单位为分
	TotalFee int `db:"total_fee"`
	//商品描述
	Body string `db:"body"`
	//商品详情
	Detail string `db:"detail"`
	//商户订单号
	OutTradeNo string `db:"out_trade_no"`
	//下单时间
	SubmitTime types.Timestamp `db:"submit_time"`
	//订单过期时间
	TimeoutExpress types.Timestamp `db:"timeout_express"`
	//付款时间
	PayTime types.Timestamp `db:"pay_time"`
	//未支付过期时间
	PayExpire types.Timestamp `db:"pay_expire"`
	//取消时间
	CancelTime types.Timestamp `db:"cancel_time"`
	//错误时间
	ErrorTime types.Timestamp `db:"error_time"`
	//状态（已提交（用户已提交但未发起支付），待支付，已支付，已取消）
	Status Status `db:"status"`
	//备注
	Remarks string `db:"remarks"`
}

// OrderPay ...
type OrderPay struct {
	attr *OrderPayAttribute
}
// NewOrderPay ...
func NewOrderPay(attr *OrderPayAttribute) *OrderPay {
	return &OrderPay{attr: attr}
}

// CanPay ...
func (m *OrderPay) CanPay() error {
	return nil
}

// SetWaitPay 设置待支付
func (m *OrderPay) SetWaitPay() {
	m.attr.Status = WaitPay
	return
}

// SetSuccess 支付成功，更新订单状态（待支付->已支付）
func (m *OrderPay) SetSuccess(payTime types.Timestamp, transactionID, buyerLogonID string, status Status) {
	return
}

// SetCancelled 设置取消订单
func (m *OrderPay) SetCancelled(cancelTime types.Timestamp, status Status) {
	return
}

// SetError 设置订单错误
func (m *OrderPay) SetError(errorTime types.Timestamp, remarks string, status Status) {
	return
}

// SetOutTradeNo 设置订单号
func (m *OrderPay) SetOutTradeNo(outTradeNo string) {
	return
}

// SetSubmitted 更新订单状态（待支付->已提交）
func (m *OrderPay) SetSubmitted(status Status) {
	return
}