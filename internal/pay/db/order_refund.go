package db

import (
	"time"

	"yumi/internal/pay/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/external/dbc"
)

//OrderRefund 退款订单
type OrderRefund struct {
	SeqID int64 `db:"seq_id"`
	trade.OrderRefund
}

//New ...
func (m *OrderRefund) New(code string) (trade.DataOrderRefund, error) {
	if code == "" {
		return &OrderRefund{}, nil
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
	if err := dbc.Get().Get(&or, sqlStr, code); err != nil {
		return &OrderRefund{}, ecode.ServerErr(err)
	}

	return &or, nil
}

//Data ...
func (m *OrderRefund) Data() trade.OrderRefund {
	return m.OrderRefund
}

//RefundCount ...
type RefundCount struct {
	Count     int `db:"count"`
	RefundFee int `db:"refund_fee"`
}

//GetRefundFee 已退款总金额
func (m *OrderRefund) GetRefundFee(orderPayCode string) (int, int, error) {
	sqlStr := `
		SELECT 
			count(seq_id) AS "count", 
			ifnull(sum("refund_fee"), 0) AS "refund_fee" 
		FROM 
			order_refund 
		WHERE 
			"order_pay_code" = ? 
		AND
			"status" = ?`
	fee := RefundCount{}
	if err := dbc.Get().Get(&fee, sqlStr, orderPayCode, trade.Refunded); err != nil {
		return 0, 0, ecode.ServerErr(err)
	}

	return fee.Count, fee.RefundFee, nil
}

//ExistRefundingOrSubmitted 是否正在提起退款
func (m *OrderRefund) ExistRefundingOrSubmitted(orderPayCode string) (bool, error) {
	sqlStr := `
		SELECT 
			seq_id 
		FROM 
			order_refund 
 		WHERE 
			"order_pay_code" = ? 
		AND 
			"status" = ? 
		OR 
			"status" = ?`
	seqID := 0
	if err := dbc.Get().Get(&seqID, sqlStr, orderPayCode, trade.Submitted, trade.Refunding); err != nil {
		return false, ecode.ServerErr(err)
	}
	
	return seqID != 0, nil
}

//Submit 提交订单
func (m *OrderRefund) Submit(code, orderPayCode string, serialNum int, notifyURL string, refundAccountGUID string, refundWay trade.Way,
	outRefundNo string, refundFee int, refundDesc string, submitTime, timeoutExpress time.Time, status trade.OrderStatus) error {
	sqlStr := `
		INSERT 
		INTO 
			order_refund 
			("code", "order_pay_code", "serial_num",  "notify_url", "refund_account_guid", "refund_way", "out_refund_no", 
			"refund_fee", "refund_desc", "submit_time", "timeout_express", "status") 
			VALUES 
			(?, ?, ?, ?, ?,  ?, ?, ?, ?, ?,  ?, ?)`
	if _, err := dbc.Get().Insert(sqlStr, code, orderPayCode, serialNum, notifyURL, refundAccountGUID,
		refundWay, outRefundNo, refundFee, refundDesc, submitTime, timeoutExpress, status); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

//SetSubmitted 更新订单状态（待支付->已提交）
func (m *OrderRefund) SetSubmitted(status trade.OrderStatus) error {
	sqlStr := `
		UPDATE 
			order_refund 
		SET 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, status, m.SeqID); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//SetRefunding 设置退款中
func (m *OrderRefund) SetRefunding(status trade.OrderStatus) error {
	sqlStr := `
		UPDATE 
			order_refund 
		SET 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, status, m.SeqID); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//SetCancelled 设置取消订单
func (m *OrderRefund) SetCancelled(cancelTime time.Time, status trade.OrderStatus) error {
	sqlStr := `
		UPDATE 
			order_refund 
		SET 
			"cancel_time" = ?, 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, cancelTime, status, m.SeqID); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//SetRefunded 设置订单错误
func (m *OrderRefund) SetRefunded(refundID string, refundedTime time.Time, status trade.OrderStatus) error {
	sqlStr := `
		UPDATE 
			order_refund 
		SET 
			"refund_id" = ?, 
			"refunded_time" = ?, 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, refundID, refundedTime, status, m.SeqID); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//SetError 设置订单错误
func (m *OrderRefund) SetError(errorTime time.Time, remarks string, status trade.OrderStatus) error {
	sqlStr := `
		UPDATE 
			order_refund 
		SET 
			"error_time" = ?, 
			"status" = ?, 
			"remarks" = ?
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, errorTime, status, remarks, m.SeqID); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}
