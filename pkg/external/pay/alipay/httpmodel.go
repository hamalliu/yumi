package alipay

type ReqPublicPrameter struct {
	AppId        string `json:"app_id"`     //支付宝分配给开发者的应用ID
	Method       string `json:"method"`     //接口名称
	Format       string `json:"format"`     //仅支持JSON
	ReturnUrl    string `json:"return_url"` //HTTP/HTTPS开头字符串
	CharSet      string `json:"charset"`    //请求使用的编码格式，如utf-8,gbk,gb2312等
	SignType     string `json:"sign_type"`  //商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign         string `json:"sign"`       //商户请求参数的签名串
	Timestamp    string `json:"timestamp"`  //发送请求的时间
	Version      string `json:"version"`    //调用的接口版本，固定为：1.0
	NotifyUrl    string `json:"notify_url"` //支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	AppAuthToken string `json:"app_auth_token"`
	BizContent   string `json:"biz_content"` //请求参数的集合
}

type RespPublicPrameter struct {
	Code    string `json:"code"`     //网关返回码
	Msg     string `json:"msg"`      //网关返回码描述
	SubCode string `json:"sub_code"` //业务返回码
	SubMsg  string `json:"sub_msg"`  //业务返回码描述
	Sign    string `json:"sign"`     //签名
}

type ReqPagePay struct {
	OutTradeNo         string         `json:"out_trade_no"`         //商户订单号
	ProductCode        string         `json:"product_code"`         //销售产品码:FAST_INSTANT_TRADE_PAY
	TotalAmount        string         `json:"total_amount"`         //订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]。
	Subject            string         `json:"subject"`              //订单标题
	Body               string         `json:"body"`                 //订单描述
	TimeExpire         string         `json:"time_expire"`          //绝对超时时间
	GoodsDetail        []GoodsDetail  `json:"goods_detail"`         //订单包含的商品列表信息，json格式，其它说明详见商品明细说明
	PassbackParams     string         `json:"passback_params"`      //公用回传参数
	ExtendParams       []ExtendParams `json:"extend_params"`        //业务扩展参数
	GoodsType          string         `json:"goods_type"`           //商品主类型 :0-虚拟类商品,1-实物类商品 注：虚拟类商品不支持使用花呗渠道
	TimeoutExpress     string         `json:"timeout_express"`      //该笔订单允许的最晚付款时间，逾期将关闭交易。
	PromoParams        string         `json:"promo_params"`         //优惠参数
	MerchantOrderNo    string         `json:"merchant_order_no"`    //商户原始订单号，最大长度限制32位
	EnablePayChannels  string         `json:"enable_pay_channels"`  //可用渠道,用户只能在指定渠道范围内支付，多个渠道以逗号分割
	StoreId            string         `json:"store_id"`             //商户门店编号
	DisablePayChannels string         `json:"disable_pay_channels"` //禁用渠道
	QrPayMode          string         `json:"qr_pay_mode"`          //PC扫码支付的方式，支持前置模式和跳转模式。
	QrcodeWidth        string         `json:"qrcode_width"`         //商户自定义二维码宽度
	InvoiceInfo        []InvoiceInfo  `json:"invoice_info"`         //开票信息
	IntegrationType    string         `json:"integration_type"`     //请求后页面的集成方式。	取值范围：1. ALIAPP：支付宝钱包内 2. PCWEB：PC端访问 默认值为PCWEB。
	RequestFromUrl     string         `json:"request_from_url"`     //请求来源地址。如果使用ALIAPP的集成方式，用户中途取消支付会返回该地址。
	BusinessParams     string         `json:"business_params"`      //商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
}

type GoodsDetail struct {
	GoodsId        string `json:"goods_id"`        //商品的编号
	AlipayGoodsId  string `json:"alipay_goods_id"` //支付宝定义的统一商品编号
	GoodsName      string `json:"goods_name"`      //商品名称
	Quantity       int    `json:"quantity"`        //商品数量
	Price          string `json:"price"`           //商品单价，单位为元
	GoodsCategory  string `json:"goods_category"`  //商品类目
	CategoriesTree string `json:"categories_tree"` //商品类目树，从商品类目根节点到叶子节点的类目id组成，类目id值使用|分割
	Body           string `json:"body"`            //商品描述信息
	ShowUrl        string `json:"show_url"`        //商品的展示地址
}

type ExtendParams struct {
	SysServiceProviderId string `json:"sys_service_provider_id"` //系统商编号
	HbFqNum              string `json:"hb_fq_num"`               //使用花呗分期要进行的分期数
	HbGqSellerPercent    string `json:"hb_fq_seller_percent"`    //使用花呗分期需要卖家承担的手续费比例的百分值，传入100代表100%
	IndustryRefluxInfo   string `json:"industry_reflux_info"`    //行业数据回流信息
	CardType             string `json:"card_type"`               //卡类型
}

type InvoiceInfo struct {
	KeyInfo KeyInfo `json:"key_info"` //开票关键信息
	Details string  `json:"details"`  //开票内容
}

type KeyInfo struct {
	IsSupportInvoice    bool   `json:"is_support_invoice"`    //该交易是否支持开票
	InvoiceMerchantName string `json:"invoice_merchant_name"` //开票商户名称：商户品牌简称|商户门店简称
	TaxNum              string `json:"tax_num"`               //税号
}

type RespPagePay struct {
	RespPublicPrameter
	TradeNo         string `json:"trade_no"`          //支付宝交易号
	OutTradeNo      string `json:"out_trade_no"`      //商户订单号
	SellerId        string `json:"seller_id"`         //收款支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字
	TotalAmount     string `json:"total_amount"`      //交易金额
	MerchantOrderNo string `json:"merchant_order_no"` //商户原始订单号
}

//
type ReqRefund struct {
	OutTradeNo     string        `json:"out_trade_no"`    //商户订单号
	TradeNo        string        `json:"trade_no"`        //支付宝交易号
	RefundAmount   string        `json:"refund_amount"`   //需要退款的金额，该金额不能大于订单金额,单位为元
	RefundCurrency string        `json:"refund_currency"` //订单退款币种信息
	RefundReason   string        `json:"refund_reason"`   //退款的原因说明
	OutRequestNo   string        `json:"out_request_no"`  //标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	OperatorId     string        `json:"operator_id"`     //商户的操作员编号
	StoreId        string        `json:"store_id"`        //商户的门店编号
	TerminalId     string        `json:"terminal_id"`     //商户的终端编号
	GoodsDetail    []GoodsDetail `json:"goods_detail"`    //订单包含的商品列表信息，json格式，其它说明详见商品明细说明
}

//
type RespRefund struct {
	RespPublicPrameter
	RefundReturn
}

//资金渠道
type TradeFundBill struct {
	FundChannel string `json:"fund_channel"` //交易使用的资金渠道
	BankCode    string `json:"bank_code"`    //银行卡支付时的银行代码
	Amount      string `json:"amount"`       //该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount"`  //渠道实际付款金额
	FundType    string `json:"fund_type"`    //渠道所使用的资金类型
}

//退回的前置资产列表
type PresetPayToolInfo struct {
	Amount         []string `json:"amount"`           //前置资产金额
	AssertTypeCode string   `json:"assert_type_code"` //前置资产类型编码，和收单支付传入的preset_pay_tool里面的类型编码保持一致。
}

type ReqRefundQuery struct {
	TradeNo      string `json:"trade_no"`       //支付宝交易号，和商户订单号不能同时为空
	OutTradeNo   string `json:"out_trade_no"`   //订单支付时传入的商户订单号,和支付宝交易号不能同时为空
	OutRequestNo string `json:"out_request_no"` //请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的外部交易号
	OrgPid       string `json:"org_pid"`        //银行间联模式下有用，其它场景请不要使用；
}

type RespRefundQuery struct {
	RespPublicPrameter
	RefundQueryReturn
}

type ReqTradeQuery struct {
	OutTradeNo   string   `json:"out_trade_no"`  //商家订单号
	TradeNo      string   `json:"trade_no"`      //支付宝交易号
	OrgPid       string   `json:"org_pid"`       //银行间联模式下有用，其它场景请不要使用；
	QueryOptions []string `json:"query_options"` //查询选项，商户通过上送该字段来定制查询返回信息
}

type RespTradeQuery struct {
	RespPublicPrameter
	TradeQueryReturn
}

type ReqTradeClose struct {
	TradeNo    string `json:"trade_no"`     //该交易在支付宝系统中的交易流水号。
	OutTradeNo string `json:"out_trade_no"` //订单支付时传入的商户订单号,和支付宝交易号不能同时为空。
	OperatorId string `json:"operator_id"`  //卖家端自定义的的操作员 ID
}

type RespTradeClose struct {
	RespPublicPrameter
	TradeCloseReturn
}

type ReqBillDownloadUrlQuery struct {
	BillType string `json:"bill_type"` //账单类型，商户通过接口或商户经开放平台授权后其所属服务商通过接口可以获取以下账单类型：trade、signcustomer；trade指商户基于支付宝交易收单的业务账单；signcustomer是指基于商户支付宝余额收入及支出等资金变动的帐务账单。
	BillDate string `json:"bill_date"` //账单时间：日账单格式为yyyy-MM-dd，最早可下载2016年1月1日开始的日账单；月账单格式为yyyy-MM，最早可下载2016年1月开始的月账单。
}

type RespBillDownloadUrlQuery struct {
	RespPublicPrameter
	BillDownloadUrlQueryReturn
}

type ReqNotify struct {
	NotifyTime string `json:"notify_time"` //通知时间
	NotifyType string `json:"notify_type"` //通知类型
	NotifyId   string `json:"notify_id"`   //通知校验id
	Charset    string `json:"charset"`     //编码格式
	Version    string `json:"version"`     //接口版本
	SignType   string `json:"sign_type"`   //签名类型
	Sign       string `json:"sign"`        //签名
	AuthAppId  string `json:"auth_app_id"` //授权方的appid

	TradeNo           string `json:"trade_no"`            //支付宝交易号
	AppId             string `json:"app_id"`              //开发者的app_id
	OutTradeNo        string `json:"out_trade_no"`        //商户订单号
	OutBizNo          string `json:"out_biz_no"`          //商户业务号
	BuyerId           string `json:"buyer_id"`            //买家支付宝用户号
	SellerId          string `json:"seller_id"`           //卖家支付宝用户号
	TradeStatus       string `json:"trade_status"`        //交易状态
	TotalAmount       string `json:"total_amount"`        //订单金额
	ReceiptAmount     string `json:"receipt_amount"`      //实收金额
	InvoiceAmount     string `json:"invoice_amount"`      //开票金额
	BuyerPayAmount    string `json:"buyer_pay_amount"`    //付款金额
	PointAmount       string `json:"point_amount"`        //集分宝金额
	RefundFee         string `json:"refund_fee"`          //总退款金额
	Subject           string `json:"subject"`             //订单标题
	Body              string `json:"body"`                //商品描述
	GmtCreate         string `json:"gmt_create"`          //交易创建时间
	GmtPayment        string `json:"gmt_payment"`         //交易付款时间
	GmtRefund         string `json:"gmt_refund"`          //交易退款时间
	GmtClose          string `json:"gmt_close"`           //交易结束时间
	FundBillList      string `json:"fund_bill_list"`      //支付金额信息
	VoucherDetailList string `json:"voucher_detail_list"` //优惠券信息
	PassbackParams    string `json:"passback_params"`     //回传参数
}
