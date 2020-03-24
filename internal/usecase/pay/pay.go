package pay

import (
	"fmt"
	"time"

	"yumi/external/dbc"
	"yumi/external/pay"
	"yumi/external/pay/ali_pagepay"
	"yumi/external/pay/wx_nativepay"
	"yumi/internal/entities/orderpay"
	"yumi/utils/internal_error"
	"yumi/utils/log"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

const (
	PayWayAliPagePay   = "ALIPAGEPAY"
	PayWayWxNative1Pay = "WXNATIVE1PAY"
)

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
		if _, err := dbc.Get().Exec(sqlStr, getCode(seqId), orderpay.OrderPayStatusSubmitted, seqId); err != nil {
			return internal_error.With(err)
		}
	}

	return nil
}

//立即支付
func Pay(code string) error {
	order := orderpay.OrderPay{}
	if err := order.Load(code); err != nil {
		return err
	}
	defer order.Release()

	//订单是否过期
	if order.PayExpire < time.Now().Format(TimeFormat) {
		return fmt.Errorf("订单已过期不能发起支付")
	}

	//如果未支付则关闭订单，将状态置为已提交
	switch order.Status {
	case orderpay.OrderPayStatusSubmitted:
		return nil

	case orderpay.OrderPayStatusWaitPay:
		switch order.PayWay {
		case PayWayAliPagePay:
			if mch, err := getAliApp(order.AppId); err != nil {
				return err
			} else {
				if err := GetAliPagePay().PayProblem(mch, order); err != nil {
					return err
				}
				return nil
			}
		case PayWayWxNative1Pay:
			if mch, err := getWxApp(order.AppId); err != nil {
				return err
			} else {
				if err := GetWxNative1().PayProblem(mch, order); err != nil {
					return err
				}
				return nil
			}
		default:
			//第一次发起支付
			return nil
		}

	case orderpay.OrderPayStatusPaid:
		return fmt.Errorf("该订单已支付")

	case orderpay.OrderPayStatusCancelled:
		return fmt.Errorf("该订单已取消")

	default:
		return fmt.Errorf("该订单不能发起支付")
	}
}

//支付遇到问题
func PayProblem(code string) error {
	//查询支付结果，如果未支付，则关闭订单，如果成功则更新订单
	order := orderpay.OrderPay{}
	if err := order.Load(code); err != nil {
		return err
	}
	defer order.Release()

	//订单是否过期
	if order.PayExpire < time.Now().Format(TimeFormat) {
		return fmt.Errorf("订单已过期不能发起支付")
	}

	//如果未支付则关闭订单，将状态置为已提交
	switch order.Status {
	case orderpay.OrderPayStatusSubmitted:
		return nil

	case orderpay.OrderPayStatusWaitPay:
		switch order.PayWay {
		case PayWayAliPagePay:
			if mch, err := getAliApp(order.AppId); err != nil {
				return err
			} else {
				if err := GetAliPagePay().PayProblem(mch, order); err != nil {
					return err
				}
				return nil
			}
		case PayWayWxNative1Pay:
			if mch, err := getWxApp(order.AppId); err != nil {
				return err
			} else {
				if err := GetWxNative1().PayProblem(mch, order); err != nil {
					return err
				}
				return nil
			}
		default:
			//第一次发起支付
			return nil
		}

	case orderpay.OrderPayStatusPaid:
		return nil

	case orderpay.OrderPayStatusCancelled:
		return nil

	default:
		return nil
	}
}

//支付完成
func PayCompleted(code string) error {
	//查询支付结果，如果成功则更新订单
	//查询支付结果，如果未支付，则关闭订单，如果成功则更新订单
	order := orderpay.OrderPay{}
	if err := order.Load(code); err != nil {
		return err
	}
	defer order.Release()

	//订单是否过期
	if order.PayExpire < time.Now().Format(TimeFormat) {
		return fmt.Errorf("订单已过期不能发起支付")
	}

	//如果未支付则关闭订单，将状态置为已提交
	switch order.Status {
	case orderpay.OrderPayStatusSubmitted:
		return nil

	case orderpay.OrderPayStatusWaitPay:
		switch order.PayWay {
		case PayWayAliPagePay:
			if mch, err := getAliApp(order.AppId); err != nil {
				return err
			} else {
				if err := GetAliPagePay().PayProblem(mch, order); err != nil {
					return err
				}
				return nil
			}
		case PayWayWxNative1Pay:
			if mch, err := getWxApp(order.AppId); err != nil {
				return err
			} else {
				if err := GetWxNative1().PayProblem(mch, order); err != nil {
					return err
				}
				return nil
			}
		default:
			//第一次发起支付
			return nil
		}

	case orderpay.OrderPayStatusPaid:
		return nil

	case orderpay.OrderPayStatusCancelled:
		return nil

	default:
		return nil
	}
}

//取消订单
func CancellOrder(code string) error {
	order := orderpay.OrderPay{}
	if err := order.Load(code); err != nil {
		return err
	}
	defer order.Release()

	if order.Status == orderpay.OrderPayStatusSubmitted {
		//将状态置为已取消
		return order.SetCancelled()
	}
	if order.Status == orderpay.OrderPayStatusWaitPay {
		//关闭订单，将状态置为已取消
		switch order.PayWay {
		case PayWayAliPagePay:
			if mch, err := getAliApp(order.AppId); err != nil {
				return err
			} else {
				return GetAliPagePay().PayProblem(mch, order)
			}
		case PayWayWxNative1Pay:
			if mch, err := getWxApp(order.AppId); err != nil {
				return err
			} else {
				return GetWxNative1().PayProblem(mch, order)
			}
		default:
			//待支付状态未设置支付方式，内部错误需检查程序。
			//设置订单为错误状态
			if err := order.SetError(); err != nil {
				return err
			}
			log.Critical("待支付状态未设置支付方式，内部错误需检查程序")
			return fmt.Errorf("该订单出错，需要内部核查")
		}
	}

	//不能取消订单
	return fmt.Errorf("该订单%s不能取消", order.Status)
}

//提交退款(只根据商品退款)
func Refund(goodsCode string) error {
	//TODO

	return nil
}

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
