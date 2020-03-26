package orderpay

import (
	"fmt"
	"sync"
	"time"

	"yumi/external/dbc"
	"yumi/external/pay"
	"yumi/external/pay/ali_pagepay"
	"yumi/external/pay/wx_nativepay"
	"yumi/utils/internal_error"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

const (
	PayWayAliPagePay   = "ALIPAGEPAY"
	PayWayWxNative1Pay = "WXNATIVE1PAY"
)

type OrderPayStatus string

const (
	OrderPayStatusSubmitted OrderPayStatus = "已提交"
	OrderPayStatusWaitPay   OrderPayStatus = "待支付"
	OrderPayStatusPaid      OrderPayStatus = "已支付"
	OrderPayStatusCancelled OrderPayStatus = "已取消"
	OrderPayStatusRefundind OrderPayStatus = "退款中"
	OrderPayStatusRefunded  OrderPayStatus = "已退款"

	OrderPayStatusError OrderPayStatus = "错误"
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

	BuyerAccountGuid string         `db:"buyeruserguid" json:"buyeruserguid"`   //买家账号guid
	TotalFee         int            `db:"totalfee" json:"totalfee"`             //订单总金额，单位为分
	Body             string         `db:"body" json:"body"`                     //商品描述
	Detail           string         `db:"detail" json:"detail"`                 //商品详情
	OutTradeNo       string         `db:"outtradeno" json:"outtradeno"`         //商户订单号
	TimeoutExpress   string         `db:"timeoutexpress" json:"timeoutexpress"` //最晚付款时间，逾期将关闭交易
	PayExpire        string         `db:"payexpire" json:"payexpire"`           //未支付过期时间
	PayTime          string         `db:"paytime" json:"paytime"`               //付款时间
	CancelTime       string         `db:"canceltime" json:"canceltime"`         //取消时间
	ErrorTime        string         `db:"errortime" json:"errortime"`           //错误时间
	SubmitTime       string         `db:"submittime" json:"submittime"`         //下单时间
	Status           OrderPayStatus `db:"status" json:"status"`                 //状态（已提交（用户已提交但未发起支付），待支付，已支付，已取消）
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

//提交订单
func SubmitOrder(notifyUrl string, totalFee int, accountGuid, body, detail, timeoutExpress string) error {
	sqlStr := `
		INSERT 
		INTO 
			order_pay 
			("out_trade_no", "notify_url", "total_fee", "body", "detail", "timeout_express", "submit_time") 
		VALUES 
			(?, ?, ?, ?, ?,  ?, ?, ?)`
	if seqId, err := dbc.Get().Insert(sqlStr,
		getOutTradeNo(), notifyUrl, totalFee, body, detail, timeoutExpress, time.Now().Format(TimeFormat)); err != nil {
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
		if _, err := dbc.Get().Exec(sqlStr, getCode(seqId), OrderPayStatusSubmitted, seqId); err != nil {
			return internal_error.With(err)
		}
	}

	return nil
}

//支付成功，更新订单状态（待支付->已支付）
func (op *OrderPay) PaySuccess(paytime string) error {
	op.Status = OrderPayStatusPaid
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
func (op *OrderPay) SetSubmitted() error {
	op.Status = OrderPayStatusSubmitted
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

//设置订单错误
func (op *OrderPay) SetError() error {
	op.Status = OrderPayStatusError
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

//设置取消订单
func (op *OrderPay) SetCancelled() error {
	op.Status = OrderPayStatusCancelled
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

//设置支付方式
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

//设置订单号
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

//======================================================================================================================
type TradeStatus string

const (
	TradeStatusSuccess  = "支付成功"
	TradeStatusNotPay   = "未支付"
	TradeStatusClosed   = "交易关闭"
	TradeStatusFinished = "交易完成"
)

type PayResult struct {
	AliPayHtml  []byte
	WxPayBizUrl string
}

//发起支付（只发起已提交订单）
func (op *OrderPay) Pay() (res PayResult, err error) {
	//订单是否过期
	if op.PayExpire < time.Now().Format(TimeFormat) {
		return res, fmt.Errorf("订单已过期，不能发起支付")
	}

	switch op.Status {
	case OrderPayStatusSubmitted:
		switch op.PayWay {
		case PayWayAliPagePay:
			if mch, err := getAliApp(op.AppId); err != nil {
				return res, err
			} else {
				if htmlByte, err := GetAliPagePay().Pay(mch, op); err != nil {
					return res, err
				} else {
					res.AliPayHtml = htmlByte
					return res, nil
				}

			}
		case PayWayWxNative1Pay:
			if mch, err := getWxApp(op.AppId); err != nil {
				return res, err
			} else {
				if bizUrl, err := GetWxNative1().Pay(mch, op); err != nil {
					return res, err
				} else {
					res.WxPayBizUrl = bizUrl
					return res, nil
				}
			}
		default:
			err := fmt.Errorf("不支持的支付方式")
			return res, internal_error.Critical(err)
		}

	case OrderPayStatusWaitPay:
		return res, fmt.Errorf("该订单待支付，不能重复发起支付")

	case OrderPayStatusPaid, OrderPayStatusRefundind, OrderPayStatusRefunded, OrderPayStatusCancelled:
		return res, fmt.Errorf("不能发起支付")

	default:
		err := fmt.Errorf("该订单状态错误")
		return res, internal_error.Critical(err)
	}
}

//查询支付状态（只查询待支付订单）
func (op *OrderPay) QueryPayStatus() (TradeStatus, error) {
	switch op.Status {
	case OrderPayStatusSubmitted, OrderPayStatusRefundind, OrderPayStatusRefunded:
		return "", fmt.Errorf("无效查询")

	case OrderPayStatusWaitPay:
		switch op.PayWay {
		case PayWayAliPagePay:
			if mch, err := getAliApp(op.AppId); err != nil {
				return "", err
			} else {
				return GetAliPagePay().QueryPayStatus(mch, op)
			}
		case PayWayWxNative1Pay:
			if mch, err := getWxApp(op.AppId); err != nil {
				return "", err
			} else {
				return GetWxNative1().QueryPayStatus(mch, op)
			}
		default:
			err := fmt.Errorf("不支持的支付方式")
			return "", internal_error.Critical(err)
		}

	case OrderPayStatusPaid:
		return TradeStatusSuccess, nil

	case OrderPayStatusCancelled:
		return TradeStatusClosed, nil

	default:
		err := fmt.Errorf("该订单状态错误")
		return "", internal_error.Critical(err)
	}
}

//关闭交易（只关闭待支付订单）
func (op *OrderPay) CloseTrade() error {
	switch op.Status {
	case OrderPayStatusSubmitted:
		return fmt.Errorf("该订单未发起支付")

	case OrderPayStatusWaitPay:
		//TODO
		return nil
	case OrderPayStatusPaid, OrderPayStatusCancelled, OrderPayStatusRefundind, OrderPayStatusRefunded:
		return nil

	default:
		err := fmt.Errorf("该订单状态错误")
		return internal_error.Critical(err)
	}
}

//TODO 退款（只退款已支付订单）
func (op *OrderPay) Refund() error {
	return nil
}

//TODO 退款查询（只查询退款中的订单）
func (op *OrderPay) QueryRefundStatus() error {
	return nil
}

//======================================================================================================================
func getWxApp(code string) (wx_nativepay.Merchant, error) {
	mch := wx_nativepay.Merchant{}
	if ret, err := GetWxApp(code); err != nil {
		return mch, err
	} else {
		mch.AppId = ret.AppId
		mch.MchId = ret.MchId
		mch.PrivateKey = ret.PrivateKey
		return mch, nil
	}
}

func getAliApp(code string) (ali_pagepay.Merchant, error) {
	mch := ali_pagepay.Merchant{}
	if ret, err := GetAliApp(code); err != nil {
		return mch, err
	} else {
		mch.AppId = ret.AppId
		mch.PublicKey = ret.PublicKey
		mch.PrivateKey = ret.PrivateKey
		return mch, nil
	}
}

func getOutTradeNo() string {
	return fmt.Sprintf("%s%s", time.Now().Format("060102150405"), pay.CreateRandomStr(10, pay.NUMBER))
}

func getCode(seqId int64) string {
	return fmt.Sprintf("%s%d", time.Now().Format("060102"), seqId)
}
