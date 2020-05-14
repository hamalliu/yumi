package dbentities

import (
	"time"

	"yumi/internal/pay/entities/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/external/dbc"
)

//支付订单
type OrderPay struct {
	SeqId int64 `db:"seq_id"`
	trade.OrderPay
}

func (m *OrderPay) New(code string) (trade.DataOrderPay, error) {
	if code == "" {
		return &OrderPay{}, nil
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
		return nil, ecode.ServerErr(err)
	}

	return &op, nil
}

//支付订单数据
func (m *OrderPay) Data() trade.OrderPay {
	return m.OrderPay
}

//提交
func (m *OrderPay) Submit(buyerAccountGuid, sellerKey, outTradeNo, notifyUrl string, totalFee int, body, detail string,
	timeoutExpress, submitTime time.Time, code string, status trade.OrderStatus) error {
	sqlStr := `
		INSERT 
		INTO 
			order_pay 
			("code", "buyer_account_guid", "seller_key", "out_trade_no", "notify_url", "total_fee", "body", "detail", 
			"timeout_express", "submit_time", "status") 
		VALUES 
			(?, ?, ?, ?, ?,  ?, ?, ?, ?, ?, ?)`
	var err error
	if m.SeqId, err = dbc.Get().Insert(sqlStr,
		code, buyerAccountGuid, sellerKey, outTradeNo, notifyUrl, totalFee, body, detail, timeoutExpress, submitTime,
		status); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

//设置支付方式
func (m *OrderPay) SetWaitPay(payWay trade.Way, appId, mchId, spbillCreateIp string, payExpire time.Time,
	status trade.OrderStatus) error {
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
	if _, err := dbc.Get().Exec(sqlStr, payWay, appId, mchId, spbillCreateIp, payExpire, status, m.SeqId); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//支付成功，更新订单状态（待支付->已支付）
func (m *OrderPay) SetSuccess(payTime time.Time, transactionId, buyerLogonId string, status trade.OrderStatus) error {
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
	if _, err := dbc.Get().Exec(sqlStr, payTime, transactionId, buyerLogonId, status, m.SeqId); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//设置取消订单
func (m *OrderPay) SetCancelled(cancelTime time.Time, status trade.OrderStatus) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"cancel_time" = ?, 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, cancelTime, status, m.SeqId); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//设置订单错误
func (m *OrderPay) SetError(errorTime time.Time, remarks string, status trade.OrderStatus) error {
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
	if _, err := dbc.Get().Exec(sqlStr, errorTime, status, remarks, m.SeqId); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//设置商户订单号
func (m *OrderPay) SetOutTradeNo(outTradeNo, notifyUrl string) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"out_trade_no" = ?, 
			"notify_url" = ? 
		WHERE 
			"seq_id" = ?`
	if _, err := dbc.Get().Exec(sqlStr, outTradeNo, notifyUrl, m.SeqId); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

//关闭订单，更新订单状态（待支付->已提交）
func (m *OrderPay) SetSubmitted(status trade.OrderStatus) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, status, m.SeqId); err != nil {
		return ecode.ServerErr(err)
	}
	return nil
}

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
