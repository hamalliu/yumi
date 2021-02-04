package entity

import "time"

// OrderRefundAttribute 退款订单
type OrderRefundAttribute struct {
	//唯一编码
	Code string `db:"code"`
	//支付订单编码
	OrderPayCode string `db:"order_pay_code"`

	//回调url
	NotifyURL string `db:"notify_url"`

	//序号（对于支付订单的序号）
	SerialNum int `db:"serial_num"`
	//退款账号guid
	RefundAccountGUID string `db:"refund_account_guid"`
	//必须和支付方式保持一致
	RefundWay string `db:"refund_way"`
	//三方支付平台退款单号
	RefundID string `db:"refund_id"`
	//商户退款单号
	OutRefundNo string `db:"out_refund_no"`
	//退款金额
	RefundFee int `db:"refund_fee"`
	//退款原因
	RefundDesc string `db:"refund_desc"`
	//退款时间
	RefundedTime time.Time `db:"refunded_time"`
	//提交时间
	SubmitTime time.Time `db:"submit_time"`
	//订单过期时间
	TimeoutExpress time.Time `db:"timeout_express"`
	//取消订单时间
	CancelTime time.Time `db:"cancel_time"`
	//状态（已提交（用户已提交但未发起支付），退款中，已退款）
	Status Status `db:"status"`
	//备注
	Remarks string `db:"remarks"`
}

// OrderRefund ...
type OrderRefund struct {
	attr *OrderRefundAttribute
}
