package data

import (
	"yumi/pkg/status"
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/usecase/trade"
	"yumi/usecase/trade/entity"
)

// MysqlDB ...
type MysqlDB struct {
	*mysqlx.Client
}

// New ...
func New(db *mysqlx.Client) *MysqlDB {
	return &MysqlDB{Client: db}
}

var _ trade.Data = &MysqlDB{}

// CreateOrderPay ...
func (db *MysqlDB) CreateOrderPay(attr entity.OrderPayAttribute) error {
	sqlStr := `
		INSERT 
		INTO 
			order_pay 
			("code", "buyer_account_guid", "seller_key", "out_trade_no", "notify_url", "total_fee", "body", "detail", 
			"timeout_express", "submit_time", "status") 
		VALUES 
			(?, ?, ?, ?, ?,  ?, ?, ?, ?, ?, ?)`
	var err error
	if _, err = db.Insert(sqlStr,
		attr.Code, attr.BuyerAccountGUID, attr.SellerKey, attr.OutTradeNo, attr.NotifyURL, attr.TotalFee, attr.Body, attr.Detail,
		attr.TimeoutExpress, attr.SubmitTime, attr.Status); err != nil {
		return status.Internal().WithDetails(err.Error())
	}

	return nil
}

// GetOrderPay ...
func (db *MysqlDB) GetOrderPay(code string) (trade.DataOrderPay, error) {
	if code == "" {
		return nil, status.Internal().WithDetails("code 不能为空")
	}
	sqlStr := `
			SELECT 
				"seq_id" AS "seq_id", 
				ifnull("code", '') AS "code",
				ifnull("trade_way", '') AS "trade_way", 
				ifnull("seller_key", '') AS "seller_key", 
				ifnull("app_id", '') AS "app_id", 
				ifnull("mch_id", '') AS "mch_id", 
				ifnull("transaction_id", '') AS "transaction_id", 
				ifnull("notify_url", '') AS "notify_url", 
				ifnull("buyer_logon_id", '') AS "buyer_logon_id", 
				ifnull("spbill_create_ip", '') AS "spbill_create_ip", 
				ifnull("buyer_account_guid", '') AS "buyer_account_guid", 
				ifnull("total_fee", 0) AS "total_fee",
				ifnull("body", '') AS "body", 
				ifnull("detail", '') AS "detail", 
				ifnull("out_trade_no", '') AS "out_trade_no", 
				ifnull("timeout_express", '') AS "timeout_express", 
				ifnull("pay_expire", '') AS "pay_expire", 
				ifnull("pay_time", '') AS "pay_time", 
				ifnull("cancel_time", '') AS "cancel_time",
				ifnull("error_time", '') AS "error_time",
				ifnull("submit_time", '') AS "submit_time",
				ifnull("status", '') AS "status",
				ifnull("remarks", '') AS "remarks"
			FROM 
				order_pay 
			WHERE 
				"code" = ?
			`
	op := OrderPay{}
	if err := db.Get(&op, sqlStr, code); err != nil {
		return nil, status.Internal().WithDetails(err.Error())
	}
	return &op, nil
}

// CreateOrderRefund ...
func (db *MysqlDB) CreateOrderRefund(attr entity.OrderRefundAttribute) error {
	sqlStr := `
		INSERT 
		INTO 
			order_refund 
			("code", "order_pay_code", "serial_num",  "notify_url", "refund_account_guid", "refund_way", "out_refund_no", 
			"refund_fee", "refund_desc", "submit_time", "timeout_express", "status") 
			VALUES 
			(?, ?, ?, ?, ?,  ?, ?, ?, ?, ?,  ?, ?)`
	if _, err := db.Insert(sqlStr,
		attr.Code, attr.RefundAccountGUID, attr.SerialNum, attr.NotifyURL, attr.RefundAccountGUID, attr.RefundWay, attr.OutRefundNo, 
		attr.RefundFee, attr.RefundDesc, attr.SubmitTime, attr.TimeoutExpress, attr.Status); err != nil {
		return status.Internal().WithDetails(err.Error())
	}

	return nil
}

// GetOrderRefund ...
func (db *MysqlDB) GetOrderRefund(code string) (trade.DataOrderRefund, error) {
	if code == "" {
		return nil, status.Internal().WithDetails("code 不能为空")
	}
	sqlStr := `
		SELECT 
			"seq_id", 
			"code", 
			"order_pay_code", 
			ifnull("serial_num", 0) AS "serial_num",
			ifnull("notify_url", '') AS "notify_url", 
			ifnull("refund_account_guid", '') AS "refund_account_guid", 
			ifnull("refund_way", '') AS "refund_way", 
			ifnull("refund_id", '') AS "refund_id", 
			ifnull("out_refund_no", '') AS "out_refund_no", 
			ifnull("refund_fee", 0) AS "refund_fee", 
			ifnull("refund_desc", '') AS "refund_desc", 
			ifnull("refunded_time", '') AS "refunded_time", 
			ifnull("submit_time", '') AS "submit_time", 
			ifnull("cancel_time", '') AS "cancel_time", 
			ifnull("status", '') AS "status",
			ifnull("remarks", '') AS "remarks"
		FROM 	
			order_refund 
		WHERE 
			code = ?
			`
	or := OrderRefund{}
	if err := db.Get(&or, sqlStr, code); err != nil {
		return nil, status.Internal().WithDetails(err.Error())
	}

	return &or, nil
}

//OrderPay 支付订单
type OrderPay struct {
	db *MysqlDB `db:"-"`

	SeqID int64 `db:"seq_id"`
	entity.OrderPayAttribute
}

var _ trade.DataOrderPay = &OrderPay{}

// Attribute ...
func (m *OrderPay) Attribute() entity.OrderPayAttribute {
	return m.OrderPayAttribute
}

//Update 设置支付方式
func (m *OrderPay) Update() error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"trade_way" = ?, 
			"seller_key" = ?, 
			"app_id" = ?, 
			"mch_id" = ?, 
			"transaction_id" = ?, 
			"notify_url" = ?, 
			"buyer_logon_id" = ?, 
			"spbill_create_ip" = ?, 
			"buyer_account_guid" = ?, 
			"total_fee" = ?,
			"body" = ?, 
			"detail" = ?, 
			"out_trade_no" = ?, 
			"timeout_express" = ?, 
			"pay_expire" = ?, 
			"pay_time" = ?, 
			"cancel_time" = ?,
			"error_time" = ?,
			"submit_time" = ?,
			"status" = ?,
			"remarks" = ?
		WHERE 
			"seq_id" = ?`
	if _, err := m.db.Exec(sqlStr, m.TradeWay, m.SellerKey, m.AppID, m.MchID, m.TransactionID, m.NotifyURL, m.BuyerLogonID, 
		m.SpbillCreateIP, m.BuyerAccountGUID, m.TotalFee, m.Body, m.Detail, m.OutTradeNo, m.TimeoutExpress, m.PayExpire, 
		m.PayTime, m.CancelTime, m.ErrorTime, m.SubmitTime, m.Status, m.Remarks, m.SeqID); err != nil {
			return status.Internal().WithDetails(err.Error())
	}
	return nil
}

//OrderRefund 退款订单
type OrderRefund struct {
	db *MysqlDB `db:"-"`

	SeqID int64 `db:"seq_id"`
	entity.OrderRefundAttribute
}

var _ trade.DataOrderRefund= &OrderRefund{}

// Attribute ...
func (m *OrderRefund) Attribute() entity.OrderRefundAttribute {
	return m.OrderRefundAttribute
}

//Update 设置支付方式
func (m *OrderRefund) Update() error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"serial_num" = ?,
			"notify_url" = ?, 
			"refund_account_guid" = ?, 
			"refund_way" = ?, 
			"refund_id" = ?, 
			"out_refund_no" = ?, 
			"refund_fee" = ?, 
			"refund_desc" = ?, 
			"refunded_time" = ?, 
			"submit_time" = ?, 
			"cancel_time" = ?, 
			"status" = ?,
			"remarks" = ?
		WHERE 
			"seq_id" = ?`
	if _, err := m.db.Exec(sqlStr, m.SerialNum, m.NotifyURL, m.RefundAccountGUID, m.RefundWay, m.RefundID, m.OutRefundNo,
		m.RefundFee, m.RefundDesc, m.RefundedTime, m.SubmitTime, m.CancelTime, m.Status, m.Remarks, m.SeqID); err != nil {
		return status.Internal().WithDetails(err.Error())
	}
	return nil
}
