package orderpay

import "time"

type Status string

const (
	Submitted Status = "已提交"
	WaitPay   Status = "待支付"
	Paid      Status = "已支付"
	Cancelled Status = "已取消"
	Refunding Status = "退款中"
	Refunded  Status = "已退款"

	Error Status = "错误"
)

//支付订单
type OrderPay struct {
	Code string `db:"code"` //唯一编码

	//交易方式
	TradeWay TradeWay `db:"pay_way"`
	//卖家key（）
	SellerKey string `db:"seller_key"`

	//开放平台应用唯一id
	AppId string `db:"appid"`
	//商户id（如果支付方式为支付宝则是seller_id，
	// 收款支付宝账号对应的支付宝唯一用户号。如果为微信则是微信商户号）
	MchId string `db:"mch_id"`
	//支付平台订单号（如果支付方式为支付宝则是支付宝订单号，如果为微信则是微信订单号）
	TransactionId string `db:"transaction_id"`
	//回调url
	NotifyUrl string `db:"notify_url"`
	//买家账号id（如果支付方式为支付宝则是买家支付宝账号id，如果为微信则是微信账号id）
	BuyerLogonId string `db:"buyer_logon_id"`

	BuyerAccountGuid string    `db:"buyer_account_guid"` //买家账号guid
	TotalFee         int       `db:"total_fee"`          //订单总金额，单位为分
	Body             string    `db:"body"`               //商品描述
	Detail           string    `db:"detail"`             //商品详情
	OutTradeNo       string    `db:"out_trade_no"`       //商户订单号
	TimeoutExpress   time.Time `db:"timeout_express"`    //最晚付款时间，逾期将关闭交易
	PayExpire        time.Time `db:"pay_expire"`         //未支付过期时间
	PayTime          time.Time `db:"pay_time"`           //付款时间
	CancelTime       time.Time `db:"cancel_time"`        //取消时间
	ErrorTime        time.Time `db:"error_time"`         //错误时间
	SubmitTime       time.Time `db:"submit_time"`        //下单时间
	Status           Status    `db:"status"`             //状态（已提交（用户已提交但未发起支付），待支付，已支付，已取消）
}

//支付完成后生成商品订单，（如果退款一定是全额退款）
type OrderGoods struct {
	BuyerAccountGuid string //买家账号guid
	Code             string //唯一编码
	GoodsCode        string //商品编码
	OrderPayCode     string //支付订单编码
	Amount           int    //商品价格，单位分
	Body             string //商品描述
	Detail           string //商品详情
	RefundExpire     string //最晚退款时间，逾期不能退款
	RefundDate       string //退款时间
	Status           string //状态（待支付，已支付，已退款，已结束）
}

//======================================================================================================================
type Data interface {
	Clone() Data
	//提交订单
	SubmitOrder(sellerKey, outTradeNo, notifyUrl string, totalFee int, body, detail string, timeoutExpress, submitTime time.Time, code string, status Status) error
	//支付成功，更新订单状态（待支付->已支付）
	PaySuccess(payTime time.Time, status Status) error
	//关闭订单，更新订单状态（待支付->已提交）
	SetSubmitted(status Status) error
	//设置订单错误
	SetError(errorTime time.Time, status Status) error
	//设置取消订单
	SetCancelled(cancelTime time.Time, status Status) error
	//设置支付方式
	SetPayWay(payWay TradeWay, appId, mchId string, status Status) error
	//设置订单号
	SetTransactionId(transactionId, buyerLogonId string) error
	//加载订单数据
	Load(code string) (OrderPay, error)
	//根据开发者appId和商户订单号加载订单数据
	LoadByOutTradeNo(appId, outTradeNo string) (OrderPay, error)
}

var data Data

func RigsterData(src Data) {
	data = src
}

func NewData() Data {
	return data.Clone()
}
