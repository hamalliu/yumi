package db

import (
	"database/sql"
	"errors"
	"time"

	"yumi/pkg/ecode"
	"yumi/pkg/stores/dbc"
	"yumi/usecase/trade"
	"yumi/usecase/trade/entity"
)

//OrderPay 支付订单
type OrderPay struct {
	SeqID int64 `db:"seq_id"`
	entity.OrderPayAttribute
}

//New ...
func (m *OrderPay) New(code string) (entity.OrderPayAttribute, error) {
	if code == "" {
		return &entity.OrderPayAttribute{}, nil
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
	if err := dbc.Get().Get(&op, sqlStr, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ecode.OrdernoDoesNotExist
		}
		return nil, ecode.ServerErr(err)
	}

	return &op, nil
}

//Entity 支付订单数据
func (m *OrderPay) Entity() trade.OrderPay {
	return m.OrderPayAttribute
}

//Submit 提交
func (m *OrderPay) Submit(buyerAccountGUID, sellerKey, outTradeNo, notifyURL string, totalFee int, body, detail string,
	timeoutExpress, submitTime time.Time, code string, status trade.Status) error {
	sqlStr := `
		INSERT 
		INTO 
			order_pay 
			("code", "buyer_account_guid", "seller_key", "out_trade_no", "notify_url", "total_fee", "body", "detail", 
			"timeout_express", "submit_time", "status") 
		VALUES 
			(?, ?, ?, ?, ?,  ?, ?, ?, ?, ?, ?)`
	var err error
	if m.SeqID, err = dbc.Get().Insert(sqlStr,
		code, buyerAccountGUID, sellerKey, outTradeNo, notifyURL, totalFee, body, detail, timeoutExpress, submitTime,
		status); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

//SetWaitPay 设置支付方式
func (m *OrderPay) SetWaitPay(payWay trade.Way, appID, mchID, spbillCreateIP string, payExpire time.Time,
	status trade.Status) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"trade_way" = ?, 
			"app_id" = ?, 
			"mch_id" = ?,
			"spbill_create_ip" = ?,
			"pay_expire" = ?,
		    "status" = ?
		WHERE 
			"seq_id" = ?`
	if _, err := dbc.Get().Exec(sqlStr, payWay, appID, mchID, spbillCreateIP, payExpire, status, m.SeqID); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//SetSuccess 支付成功，更新订单状态（待支付->已支付）
func (m *OrderPay) SetSuccess(payTime time.Time, transactionID, buyerLogonID string, status trade.Status) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"pay_time" = ?,
			"transaction_id" = ?, 
			"buyer_logon_id" = ?,
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, payTime, transactionID, buyerLogonID, status, m.SeqID); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//SetCancelled 设置取消订单
func (m *OrderPay) SetCancelled(cancelTime time.Time, status trade.Status) error {
	sqlStr := `
		UPDATE 
			order_pay 
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

//SetError 设置订单错误
func (m *OrderPay) SetError(errorTime time.Time, remarks string, status trade.Status) error {
	sqlStr := `
		UPDATE 
			order_pay 
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

//SetOutTradeNo 设置商户订单号
func (m *OrderPay) SetOutTradeNo(outTradeNo, notifyURL string) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"out_trade_no" = ?, 
			"notify_url" = ? 
		WHERE 
			"seq_id" = ?`
	if _, err := dbc.Get().Exec(sqlStr, outTradeNo, notifyURL, m.SeqID); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//SetSubmitted 关闭订单，更新订单状态（待支付->已提交）
func (m *OrderPay) SetSubmitted(status trade.Status) error {
	sqlStr := `
		UPDATE 
			order_pay 
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

//GetOrderPayCodesSubmittedAndWaitPay ...
func GetOrderPayCodesSubmittedAndWaitPay() ([]string, error) {
	sqlStr := `
		SELECT 
			"code" 
		FROM 
			order_pay 
		WHERE 
			"status" = ? 
		OR 
			"status" = ?`
	codes := []string{}
	if err := dbc.Get().Select(&codes, sqlStr, trade.Submitted, trade.WaitPay); err != nil {
		return codes, ecode.ServerErr(err)
	}

	return codes, nil
}
