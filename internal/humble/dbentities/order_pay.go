package dbentities

import (
	"time"

	"yumi/external/dbc"
	"yumi/internal/entities/orderpay"
	"yumi/utils/internal_error"
)

//支付订单
type Data struct {
	SeqId int64 `db:"seq_id"`
	orderpay.OrderPay
}

func Init() {
	orderpay.RigsterData(&Data{})
}

func (m Data) Clone() orderpay.Data {
	return &Data{}
}

func (m *Data) SubmitOrder(outTradeNo, notifyUrl string, totalFee int, body, detail string, timeoutExpress, submitTime time.Time, code string, status orderpay.Status) error {
	sqlStr := `
		INSERT 
		INTO 
			order_pay 
			("out_trade_no", "notify_url", "total_fee", "body", "detail", "timeout_express", "submit_time") 
		VALUES 
			(?, ?, ?, ?, ?,  ?, ?, ?)`
	if _, err := dbc.Get().Insert(sqlStr,
		outTradeNo, notifyUrl, totalFee, body, detail, timeoutExpress, submitTime); err != nil {
		return internal_error.With(err)
	} else {
		sqlStr = `
		UPDATE
			order_pay 
		SET 
			"code" = ?,
		    "status" = ?
		WHERE 
			"seq_id" = ?`
		if _, err := dbc.Get().Exec(sqlStr, code, status, m.SeqId); err != nil {
			return internal_error.With(err)
		} else {
			return nil
		}
	}
}

//支付成功，更新订单状态（待支付->已支付）
func (m *Data) PaySuccess(payTime time.Time, status orderpay.Status) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"pay_time" = ?, 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, payTime, status, m.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

//关闭订单，更新订单状态（待支付->已提交）
func (m *Data) SetSubmitted(status orderpay.Status) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, status, m.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

//设置订单错误
func (m *Data) SetError(errorTime time.Time, status orderpay.Status) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"error_time" = ?, 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, errorTime, status, m.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

//设置取消订单
func (m *Data) SetCancelled(cancelTime time.Time, status orderpay.Status) error {
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
		return internal_error.With(err)
	}
	return nil
}

//设置支付方式
func (m *Data) SetPayWay(payWay orderpay.TradeWay, appId, mchId string, status orderpay.Status) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"pay_way" = ?, 
			"app_id" = ?, 
			"mch_id" = ?, 
		    "status" = ?
		WHERE 
			"seq_id" = ?`
	if _, err := dbc.Get().Exec(sqlStr, payWay, appId, mchId, status, m.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

//设置订单号
func (m *Data) SetTransactionId(transactionId, buyerLogonId string) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"transaction_id" = ?, 
			"buyer_logon_id" = ? 
		WHERE 
			"seq_id" = ?`
	if _, err := dbc.Get().Exec(sqlStr, transactionId, buyerLogonId, m.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

//加载订单数据
func (m *Data) Load(code string) (orderpay.OrderPay, error) {
	sqlStr := `
			SELECT 
				"seq_id" AS "seqid", 
				ifnull("code", '') AS "code", 
				ifnull("pay_way", '') AS "payway", 
				ifnull("app_id", '') AS "appid", 
				ifnull("mch_id", '') AS "mchid", 
				ifnull("transaction_id", '') AS "trasactionid", 
				ifnull("buyer_logon_id", '') AS "buyerlogonid", 
				ifnull("notify_url", '') AS "notifyurl",
				ifnull("total_fee", 0) AS "totalfee",
				ifnull("body", '') AS "body", 
				ifnull("detail", '') AS "detail", 
				ifnull("out_trade_no", '') AS "outtradeno", 
				ifnull("timeout_express", '') AS "timeoutexpress", 
				ifnull("pay_expire", '') AS "payexpire", 
				ifnull("pay_date", '') AS "paydate", 
				ifnull("cancel_time", '') AS "canceltime",
				ifnull("err_time", '') AS "errtime",
				ifnull("submit_time", '') AS "submittime",
				ifnull("status", '') AS "status"
			FROM 
				order_pay 
			WHERE 
				"code" = ?
			`
	if err := dbc.Get().Get(&m, sqlStr, code); err != nil {
		return m.OrderPay, internal_error.With(err)
	}

	return m.OrderPay, nil
}

//根据开发者appId和商户订单号加载订单数据
func (m *Data) LoadByOutTradeNo(appId, outTradeNo string) (orderpay.OrderPay, error) {
	sqlStr := `
			SELECT 
				"seq_id" AS "seqid", 
				ifnull("code", '') AS "code", 
				ifnull("pay_way", '') AS "payway", 
				ifnull("app_id", '') AS "appid", 
				ifnull("mch_id", '') AS "mchid", 
				ifnull("transaction_id", '') AS "trasactionid", 
				ifnull("buyer_logon_id", '') AS "buyerlogonid", 
				ifnull("notify_url", '') AS "notifyurl",
				ifnull("total_fee", 0) AS "totalfee",
				ifnull("body", '') AS "body", 
				ifnull("detail", '') AS "detail", 
				ifnull("out_trade_no", '') AS "outtradeno", 
				ifnull("timeout_express", '') AS "timeoutexpress", 
				ifnull("pay_expire", '') AS "payexpire", 
				ifnull("pay_date", '') AS "paydate", 
				ifnull("cancel_time", '') AS "canceltime",
				ifnull("err_time", '') AS "errtime",
				ifnull("submit_time", '') AS "submittime",
				ifnull("status", '') AS "status"
			FROM 
				order_pay 
			WHERE 
				"app_id" = ? 
			AND 
			    "out_trade_no = ?"
			`
	if err := dbc.Get().Get(&m, sqlStr, appId, outTradeNo); err != nil {
		return m.OrderPay, internal_error.With(err)
	}

	return m.OrderPay, nil
}
