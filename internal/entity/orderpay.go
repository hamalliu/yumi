package entity

import (
	"fmt"
	"sync"
	"yumi/external/dbc"
	"yumi/utils/internal_error"
)

const (
	OrderPayStatusSubmitted = "已提交"
	OrderPayStatusWaitPay   = "待支付"
	OrderPayStatusPaid      = "已支付"
	OrderPayStatusCancelled = "已取消"

	OrderPayStatusError = "错误"
)

var orderSync map[string]*sync.Mutex

//支付订单
type OrderPay struct {
	SeqId int64  `db:"seqid" json:"seqid"`
	Code  string `db:"code" json:"code"` //唯一编码

	//支付方式
	PayWay string `db:"payway" json:"payway"`

	//开放平台应用唯一id
	AppId string `db:"appid" json:"appid"`

	//商户id（如果支付方式为支付宝则是seller_id，
	// 收款支付宝账号对应的支付宝唯一用户号。如果为微信则是微信商户号）
	MchId string `db:"mchid" json:"mchid"`

	//支付平台订单号（如果支付方式为支付宝则是支付宝订单号，如果为微信则是微信订单号）
	TransactionId string `db:"transactionid" json:"transactionid"`

	//买家账号id（如果支付方式为支付宝则是买家支付宝账号id，如果为微信则是微信账号id）
	BuyerLogonId string `db:"buyerlogonid" json:"buyerlogonid"`

	//买家账号id（如果支付方式为支付宝则是买家支付宝账号id，如果为微信则是微信账号id）
	NotifyUrl string `db:"notifyurl" json:"notifyurl"`

	BuyerAccountGuid string `db:"buyeruserguid" json:"buyeruserguid"`   //买家账号guid
	TotalFee         int    `db:"totalfee" json:"totalfee"`             //订单总金额，单位为分
	Body             string `db:"body" json:"body"`                     //商品描述
	Detail           string `db:"detail" json:"detail"`                 //商品详情
	OutTradeNo       string `db:"outtradeno" json:"outtradeno"`         //商户订单号
	TimeoutExpress   string `db:"timeoutexpress" json:"timeoutexpress"` //最晚付款时间，逾期将关闭交易
	PayExpire        string `db:"payexpire" json:"payexpire"`           //未支付过期时间
	PayTime          string `db:"paytime" json:"paytime"`               //付款时间
	CancelTime       string `db:"canceltime" json:"canceltime"`         //取消时间
	ErrorTime        string `db:"errortime" json:"errortime"`           //错误时间
	SubmitTime       string `db:"submittime" json:"submittime"`         //下单时间
	Status           string `db:"status" json:"status"`                 //状态（已提交（用户已提交但未发起支付），待支付，已支付，已取消）
}

//支付完成后生成商品订单，（如果退款一定是全额退款）
type OrderGoods struct {
	SeqId            int64  `db:"seqId" json:"seqid"`
	BuyerAccountGuid string `db:"buyeruserguid" json:"buyeruserguid"` //买家账号guid
	Code             string `db:"code" json:"code"`                   //唯一编码
	GoodsCode        string `db:"goodscode" json:"goodscode"`         //商品编码
	OrderPayCode     string `db:"orderpaycode" json:"orderpaycode"`   //支付订单编码
	Amount           int    `db:"amount" json:"amount"`               //商品价格，单位分
	Body             string `db:"body" json:"body"`                   //商品描述
	Detail           string `db:"detail" json:"detail"`               //商品详情
	RefundExpire     string `db:"refundexpire" json:"refundexpire"`   //最晚退款时间，逾期不能退款
	RefundDate       string `db:"refunddate" json:"refunddate"`       //退款时间
	Status           string `db:"status" json:"status"`               //状态（待支付，已支付，已退款，已结束）
}

//支付成功，更新订单状态（待支付->已支付）
func (op *OrderPay) PaySuccess(paytime string) error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"pay_time" = ?, 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, paytime, OrderPayStatusPaid, op.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

//关闭订单，更新订单状态（待支付->已提交）
func (op *OrderPay) PayClose() error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, OrderPayStatusSubmitted, op.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

func (op *OrderPay) SetError() error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"error_time" = sysdate(), 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, OrderPayStatusError, op.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

func (op *OrderPay) SetCancelled() error {
	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"cancel_time" = sysdate(), 
			"status" = ? 
		WHERE 
			"seq_id" = ?
		`
	if _, err := dbc.Get().Exec(sqlStr, OrderPayStatusCancelled, op.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

func (op *OrderPay) SetPayWay(payWay, appId, mchId string) error {
	op.PayWay = payWay
	op.AppId = appId
	op.MchId = mchId

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
	if _, err := dbc.Get().Exec(sqlStr,
		payWay, appId, mchId, OrderPayStatusWaitPay, op.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

func (op *OrderPay) SetTransactionId(transactionId, buyerLogonId string) error {
	op.TransactionId = transactionId
	op.BuyerLogonId = buyerLogonId

	sqlStr := `
		UPDATE 
			order_pay 
		SET 
			"transaction_id" = ?, 
			"buyer_logon_id" = ? 
		WHERE 
			"seq_id" = ?`
	if _, err := dbc.Get().Exec(sqlStr, transactionId, buyerLogonId, op.SeqId); err != nil {
		return internal_error.With(err)
	}
	return nil
}

//加载订单数据
func (op *OrderPay) Load(code string) error {
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
	if err := dbc.Get().Get(op, sqlStr, code); err != nil {
		return internal_error.With(err)
	}

	if orderSync[op.Code] == nil {
		orderSync[op.Code] = new(sync.Mutex)
	}
	orderSync[op.Code].Lock()
	return nil
}

//释放数据
func (op *OrderPay) Release() error {
	if orderSync[op.Code] == nil {
		err := fmt.Errorf("无法释放锁，可能造成死锁")
		return internal_error.With(err)
	}
	orderSync[op.Code].Unlock()
	return nil
}

//根据开发者appId和商户订单号加载订单数据
func (op *OrderPay) LoadByOutTradeNo(appId, outTradeNo string) error {
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
	if err := dbc.Get().Get(op, sqlStr, appId, outTradeNo); err != nil {
		return internal_error.With(err)
	}

	if orderSync[op.Code] == nil {
		orderSync[op.Code] = new(sync.Mutex)
	}
	orderSync[op.Code].Lock()
	return nil
}
