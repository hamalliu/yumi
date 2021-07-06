package entity

import (
	"yumi/pkg/types"
	"yumi/pkg/status"
	"yumi/usecase/trade/entity/internal"
)

// PayExpireSecond 支付过期期限：30分钟
const PayExpireSecond = 30 * 60

//Status 订单状态
type Status string

const (
	//Submitted "已提交" #包含支付订单和退款订单
	Submitted Status = "SUBMITTED"
	//WaitPay "待支付"
	WaitPay Status = "WAIT_PAY"
	//Paid "已支付"
	Paid Status = "PAID"
	//Cancelled "已取消" #包含支付订单和退款订单
	Cancelled Status = "CANCELLED"
	//Refunding "退款中"
	Refunding Status = "REFUNDING"
	//Refunded "已退款"
	Refunded Status = "REFUNDED"
	//Error "错误" #包含支付订单和退款订单
	Error Status = "ERROR"
)

//OrderPayAttribute 支付订单
type OrderPayAttribute struct {
	//唯一编码
	Code string `db:"code"`
	//交易方式
	TradeWay string `db:"trade_way"`
	//卖家key（对本支付系统唯一）
	SellerKey string `db:"seller_key"`
	//开放平台应用唯一id
	AppID string `db:"app_id"`
	//商户id（如果支付方式为支付宝则是seller_id，
	// 收款支付宝账号对应的支付宝唯一用户号。如果为微信则是微信商户号）
	MchID string `db:"mch_id"`
	//支付平台订单号（如果支付方式为支付宝则是支付宝订单号，如果为微信则是微信订单号）
	TransactionID string `db:"transaction_id"`
	//回调url
	NotifyURL string `db:"notify_url"`
	//买家账号id（如果支付方式为支付宝则是买家支付宝账号id，如果为微信则是微信账号id）
	BuyerLogonID string `db:"buyer_logon_id"`
	//用户支付客户端ip
	SpbillCreateIP string `db:"spbill_create_ip"`
	//买家账号guid
	BuyerAccountGUID string `db:"buyer_account_guid"`
	//订单总金额，单位为分
	TotalFee int `db:"total_fee"`
	//商品描述
	Body string `db:"body"`
	//商品详情
	Detail string `db:"detail"`
	//商户订单号
	OutTradeNo string `db:"out_trade_no"`
	//下单时间
	SubmitTime types.Timestamp `db:"submit_time"`
	//订单过期时间
	TimeoutExpress types.Timestamp `db:"timeout_express"`
	//付款时间
	PayTime types.Timestamp `db:"pay_time"`
	//未支付过期时间
	PayExpire types.Timestamp `db:"pay_expire"`
	//取消时间
	CancelTime types.Timestamp `db:"cancel_time"`
	//错误时间
	ErrorTime types.Timestamp `db:"error_time"`
	//状态（已提交（用户已提交但未发起支付），待支付，已支付，已取消）
	Status Status `db:"status"`
	//备注
	Remarks string `db:"remarks"`
}

// OrderPay ...
type OrderPay struct {
	attr *OrderPayAttribute
}

// NewOrderPay ...
func NewOrderPay(attr *OrderPayAttribute) OrderPay {
	return OrderPay{attr: attr}
}

// Submit ...
func (m *OrderPay) Submit() error {
	if m.attr.TimeoutExpress < types.NowTimestamp()+types.Timestamp(PayExpireSecond) {
		return status.InvalidArgument()
	}
	m.attr.Code = internal.GetCode(internal.OrderPayCode)
	m.attr.SubmitTime = types.NowTimestamp()
	return nil
}

// Cancel ...
func (m *OrderPay) Cancel() error {
	if m.attr.Status == Submitted {
		m.setCancelled()
	} else if m.attr.Status == WaitPay {
		thirdpf, err := GetThirdpf(Way(m.attr.TradeWay))
		if err != nil {
			return err
		}
		ret1, err := thirdpf.QueryPayStatus(*m.attr)
		if err != nil {
			return err
		}
		if ret1.TradeStatus == StatusTradePlatformSuccess {
			m.setPaid(ret1.TransactionID, ret1.BuyerLogonID)
			return status.FailedPrecondition().WithMessage(OrderFinishedRefuseCancel)
		}
		err = thirdpf.TradeClose(*m.attr)
		if err != nil {
			return err
		}
		m.setCancelled()
	}
	if m.attr.Status == Cancelled {
		return nil
	}
	return status.InvalidRequest()
}

// Pay ...
func (m *OrderPay) Pay(tradeWay, notifyURL, clientIP string) (string, error) {
	if m.attr.Status == Submitted {
		now := types.NowTimestamp()
		if now > m.attr.TimeoutExpress {
			return "", status.DeadlineExceeded().WithMessage(OrderTimeout)
		}

		thirdpf, err := GetThirdpf(Way(tradeWay))
		if err != nil {
			return "", status.Internal().WithDetails(err)
		}

		m.setWaitPay(tradeWay, notifyURL, clientIP)
		ret, err := thirdpf.Pay(*m.attr)
		if err != nil {
			return "", status.Internal().WithDetails(err)
		}

		return ret.Data, nil
	} else if m.attr.Status == WaitPay {
		thirdpf1, err := GetThirdpf(Way(m.attr.TradeWay))
		if err != nil {
			return "", status.Internal().WithDetails(err)
		}
		ret1, err := thirdpf1.QueryPayStatus(*m.attr)
		if err != nil {
			return "", status.Internal().WithDetails(err)
		}
		if ret1.TradeStatus == StatusTradePlatformSuccess {
			// 设置已支付
			m.setPaid(ret1.TransactionID, ret1.BuyerLogonID)
			return "", status.AlreadyExists().WithMessage(OrderAlreadyExists)
		} else if ret1.TradeStatus == StatusTradePlatformNotPay {
			// TODO: 直接返回之前数据
		} else if ret1.TradeStatus == StatusTradePlatformClosed {
			// 重新下单
			thirdpf2, err := GetThirdpf(Way(tradeWay))
			if err != nil {
				return "", status.Internal().WithDetails(err)
			}
			m.setWaitPay(tradeWay, notifyURL, clientIP)
			ret2, err := thirdpf2.Pay(*m.attr)
			if err != nil {
				return "", status.Internal().WithDetails(err)
			}
			return ret2.Data, nil
		}
	}

	return "", status.InvalidRequest()
}

// QueryPaid ...
func (m *OrderPay) QueryPaid() (bool, error) {
	thirdpf, err := GetThirdpf(Way(m.attr.TradeWay))
	if err != nil {
		return false, err
	}

	if m.attr.Status == WaitPay {
		//查询支付状态
		tpq, err := thirdpf.QueryPayStatus(*m.attr)
		if err != nil {
			return false, err
		}
		if tpq.TradeStatus == StatusTradePlatformSuccess {
			m.setPaid(tpq.TransactionID, tpq.BuyerLogonID)
			return true, nil
		}
		if tpq.TradeStatus == StatusTradePlatformNotPay {
			return false, nil
		}
	} else if m.attr.Status == Paid {
		return true, nil
	}

	return false, status.InvalidRequest()
}

// PayNotifyReq ...
func (m *OrderPay) PayNotifyReq() {

}

// setWaitPay 设置待支付
func (m *OrderPay) setWaitPay(tradeWay, notifyURL, clientIP string) {
	m.attr.OutTradeNo = internal.GetOutTradeNo()
	m.attr.TradeWay = tradeWay
	m.attr.NotifyURL = notifyURL
	m.attr.SpbillCreateIP = clientIP
	m.attr.PayExpire = types.NowTimestamp() + PayExpireSecond
	m.attr.Status = WaitPay
}

// setPaid 支付成功，更新订单状态（待支付->已支付）
func (m *OrderPay) setPaid(transactionID, buyerLogonID string) {
	m.attr.BuyerLogonID = buyerLogonID
	m.attr.TransactionID = transactionID
	m.attr.Status = Paid
	m.attr.PayTime = types.NowTimestamp()
}

// setCancelled 设置取消订单
func (m *OrderPay) setCancelled() {
	m.attr.Status = Cancelled
	m.attr.CancelTime = types.NowTimestamp()
}
