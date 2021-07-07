package data

import (
	"github.com/pkg/errors"

	"yumi/usecase/trade/entity"
)

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
	if _, err = db.curd.Insert(sqlStr,
		attr.Code, attr.BuyerAccountGUID, attr.SellerKey, attr.OutTradeNo, attr.NotifyURL, attr.TotalFee, attr.Body, attr.Detail,
		attr.TimeoutExpress, attr.SubmitTime, attr.Status); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

//UpdateOrderPay 设置支付方式
func (db *MysqlDB) UpdateOrderPay(attr entity.OrderPayAttribute) error {
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
			"code" = ?`
	if _, err := db.curd.Exec(sqlStr, attr.TradeWay, attr.SellerKey, attr.AppID, attr.MchID, attr.TransactionID, attr.NotifyURL, attr.BuyerLogonID,
		attr.SpbillCreateIP, attr.BuyerAccountGUID, attr.TotalFee, attr.Body, attr.Detail, attr.OutTradeNo, attr.TimeoutExpress, attr.PayExpire,
		attr.PayTime, attr.CancelTime, attr.ErrorTime, attr.SubmitTime, attr.Status, attr.Remarks, attr.Code); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetOrderPay ...
func (db *MysqlDB) GetOrderPay(code string) (op entity.OrderPayAttribute, err error) {
	if code == "" {
		return
	}
	sqlStr := `
			SELECT 
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
	if err = db.curd.Get(&op, sqlStr, code); err != nil {
		return op, errors.WithStack(err)
	}
	return op, nil
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
	if _, err := db.curd.Insert(sqlStr,
		attr.Code, attr.RefundAccountGUID, attr.SerialNum, attr.NotifyURL, attr.RefundAccountGUID, attr.RefundWay, attr.OutRefundNo,
		attr.RefundFee, attr.RefundDesc, attr.SubmitTime, attr.TimeoutExpress, attr.Status); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

//UpdateOrderRefund 设置支付方式
func (db *MysqlDB) UpdateOrderRefund(attr entity.OrderRefundAttribute) error {
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
			"code" = ?`
	if _, err := db.curd.Exec(sqlStr, attr.SerialNum, attr.NotifyURL, attr.RefundAccountGUID, attr.RefundWay, attr.RefundID, attr.OutRefundNo,
		attr.RefundFee, attr.RefundDesc, attr.RefundedTime, attr.SubmitTime, attr.CancelTime, attr.Status, attr.Remarks, attr.Code); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetOrderRefund ...
func (db *MysqlDB) GetOrderRefund(code string) (or entity.OrderRefundAttribute, err error) {
	if code == "" {
		return
	}
	sqlStr := `
		SELECT 
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
	if err := db.curd.Get(&or, sqlStr, code); err != nil {
		return or, errors.WithStack(err)
	}

	return or, nil
}
