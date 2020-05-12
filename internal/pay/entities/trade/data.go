package trade

import (
	"time"
)

type OrderStatus string

const (
	Submitted OrderStatus = "SUBMITTED" //"已提交" //包含支付订单和退款订单
	WaitPay   OrderStatus = "WAIT_PAY"  //"待支付"
	Paid      OrderStatus = "PAID"      //"已支付"
	Cancelled OrderStatus = "CANCELLED" //"已取消" //包含支付订单和退款订单
	Refunding OrderStatus = "REFUNDING" //"退款中"
	Refunded  OrderStatus = "REFUNDED"  //"已退款"

	Error OrderStatus = "ERROR" //"错误" //包含支付订单和退款订单
)

//支付订单
type OrderPay struct {
	Code string `db:"code"` //唯一编码

	//交易方式
	TradeWay Way `db:"trade_way"`
	//卖家key（对本支付系统唯一）
	SellerKey string `db:"seller_key"`

	//开放平台应用唯一id
	AppId string `db:"app_id"`
	//商户id（如果支付方式为支付宝则是seller_id，
	// 收款支付宝账号对应的支付宝唯一用户号。如果为微信则是微信商户号）
	MchId string `db:"mch_id"`
	//支付平台订单号（如果支付方式为支付宝则是支付宝订单号，如果为微信则是微信订单号）
	TransactionId string `db:"transaction_id"`
	//回调url
	NotifyUrl string `db:"notify_url"`
	//买家账号id（如果支付方式为支付宝则是买家支付宝账号id，如果为微信则是微信账号id）
	BuyerLogonId string `db:"buyer_logon_id"`
	//用户支付客户端ip
	SpbillCreateIp string `db:"spbill_create_ip"`

	BuyerAccountGuid string      `db:"buyer_account_guid"` //买家账号guid
	TotalFee         int         `db:"total_fee"`          //订单总金额，单位为分
	Body             string      `db:"body"`               //商品描述
	Detail           string      `db:"detail"`             //商品详情
	OutTradeNo       string      `db:"out_trade_no"`       //商户订单号
	SubmitTime       time.Time   `db:"submit_time"`        //下单时间
	TimeoutExpress   time.Time   `db:"timeout_express"`    //订单过期时间
	PayTime          time.Time   `db:"pay_time"`           //付款时间
	PayExpire        time.Time   `db:"pay_expire"`         //未支付过期时间
	CancelTime       time.Time   `db:"cancel_time"`        //取消时间
	ErrorTime        time.Time   `db:"error_time"`         //错误时间
	Status           OrderStatus `db:"status"`             //状态（已提交（用户已提交但未发起支付），待支付，已支付，已取消）
	Remarks          string      `db:"remarks"`            //备注
}

//支付数据接口
type DataOrderPay interface {
	New(code string) (DataOrderPay, error)
	Data() OrderPay

	//提交订单
	Submit(buyerAccountGuid, sellerKey, outTradeNo, notifyUrl string, totalFee int, body, detail string, timeoutExpress,
		submitTime time.Time, code string, status OrderStatus) error
	//设置待支付
	SetWaitPay(payWay Way, appId, mchId, spbillCreateIp string, payExpire time.Time, status OrderStatus) error
	//支付成功，更新订单状态（待支付->已支付）
	SetSuccess(payTime time.Time, transactionId, buyerLogonId string, status OrderStatus) error
	//设置取消订单
	SetCancelled(cancelTime time.Time, status OrderStatus) error
	//设置订单错误
	SetError(errorTime time.Time, remarks string, status OrderStatus) error

	//设置订单号
	SetOutTradeNo(outTradeNo, notifyUrl string) error
	//更新订单状态（待支付->已提交）
	SetSubmitted(status OrderStatus) error
}

var op DataOrderPay

func RigsterDataOrderPay(src DataOrderPay) {
	op = src
}

func NewDataOrderPay(code string) (DataOrderPay, error) {
	return op.New(code)
}

//======================================================================================================================

type OrderRefund struct {
	Code         string `db:"code"`           //唯一编码
	OrderPayCode string `db:"order_pay_code"` //支付订单编码

	//回调url
	NotifyUrl string `db:"notify_url"`

	SerialNum         int         `db:"serial_num"`          //序号（对于支付订单的序号）
	RefundAccountGuid string      `db:"refund_account_guid"` //退款账号guid
	RefundWay         Way         `db:"refund_way"`          //必须和支付方式保持一致
	RefundId          string      `db:"refund_id"`           //三方支付平台退款单号
	OutRefundNo       string      `db:"out_refund_no"`       //商户退款单号
	RefundFee         int         `db:"refund_fee"`          //退款金额
	RefundDesc        string      `db:"refund_desc"`         //退款原因
	RefundedTime      time.Time   `db:"refunded_time"`       //退款时间
	SubmitTime        time.Time   `db:"submit_time"`         //提交时间
	TimeoutExpress    time.Time   `db:"timeout_express"`     //订单过期时间
	CancelTime        time.Time   `db:"cancel_time"`         //取消订单时间
	Status            OrderStatus `db:"status"`              //状态（已提交（用户已提交但未发起支付），退款中，已退款）
	Remarks           string      `db:"remarks"`             //备注
}

//退款数据接口
type DataOrderRefund interface {
	New(code string) (DataOrderRefund, error)
	Data() OrderRefund

	//已退款总金额
	GetRefundFee(orderPayCode string) (int, int, error)
	//是否正在提起退款
	ExistRefundingOrSubmitted(orderPayCode string) (bool, error)
	//提交订单
	Submit(code, orderPayCode string, serialNum int, notifyUrl string, refundAccountGuid string, refundWay Way,
		outRefundNo string, refundFee int, refundDesc string, submitTime, timeoutExpress time.Time, status OrderStatus) error
	//更新订单状态（待支付->已提交）
	SetSubmitted(status OrderStatus) error
	//设置退款中
	SetRefunding(status OrderStatus) error
	//设置取消订单
	SetCancelled(cancelTime time.Time, status OrderStatus) error
	//设置已退款
	SetRefunded(refundId string, refundedTime time.Time, status OrderStatus) error
	//设置订单错误
	SetError(errorTime time.Time, remarks string, status OrderStatus) error
}

var or DataOrderRefund

func RigsterDataOrderRefund(src DataOrderRefund) {
	or = src
}

func NewDataOrderRefund(code string) (DataOrderRefund, error) {
	return or.New(code)
}
