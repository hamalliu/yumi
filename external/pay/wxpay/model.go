package wxpay

import "time"

//二维码下单传入数据
type UnifiedOrder struct {
	//商品描述：商品简单描述，该字段请按照规范传递，
	// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=4_2
	Body string `NATIVE:"true" MWEB:"true" APP:"true" JSAPI:"true"`

	//商品详情：商品详细描述，对于使用单品优惠的商户，该字段必须按照规范上传，
	// https://pay.weixin.qq.com/wiki/doc/api/danpin.php?chapter=9_102&index=2
	Detail string

	//附加数据:附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。
	Attach string `NATIVE:"true" MWEB:"true" APP:"true" JSAPI:"true"`

	//商户订单号：针对商户唯一订单标识
	OutTradeNo string `NATIVE:"true" MWEB:"true" APP:"true" JSAPI:"true"`

	//标价金额：订单总金额，单位为分
	//https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=4_2
	TotalFee int `NATIVE:"true" MWEB:"true" APP:"true" JSAPI:"true"`

	//订单优惠标记，使用代金券或立减优惠功能时需要的参数
	//https://pay.weixin.qq.com/wiki/doc/api/tools/sp_coupon.php?chapter=12_7&index=3
	GoodsTag string

	//通知地址：异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
	NotifyUrl string `NATIVE:"true" MWEB:"true" APP:"true" JSAPI:"true"`

	//商品ID
	ProductId string `NATIVE:"true"`

	//指定支付方式：上传此参数no_credit--可限制用户不能使用信用卡支付
	LimitPay string

	//用户标识
	OpendId string `JSAPI:"true"`

	//用户请求的客户端ip
	SpbillCreateIp string `NATIVE:"true" MWEB:"true" APP:"true" JSAPI:"true"`

	//支付过期时间
	PayExpire time.Time `NATIVE:"true" MWEB:"true" APP:"true" JSAPI:"true"`

	//场景信息：该字段常用于线下活动时的场景信息上报，支持上报实际门店信息，商户也可以按需求自己上报相关信息。
	// 该字段为JSON对象数据，
	// 对象格式为{"store_info":{"id": "门店ID","name": "名称","area_code": "编码","address": "地址" }} ，
	// 字段详细说明
	//字段名	        | 变量名	        必填	   类型	        示例值	            描述
	//==================================================================================================================
	//-门店id	    | id	        否	   String(32)	SZTX001	            门店编号，由商户自定义
	//-门店名称	    | name	        否	   String(64)	腾讯大厦腾大餐厅	    门店名称 ，由商户自定义
	//-门店行政区划码	| area_code	    否	   String(6)	440305	            门店所在地行政区划码，详细见《最新县及县以上行政区划代码》
	//-门店详细地址	| address	    否	   String(128)	科技园中一路腾讯大厦	门店详细地址 ，由商户自定义
	//
	//==================================================================================================================
	//                                       MWEB                                                                     //
	//                                                                                                                //
	//==================================================================================================================
	//该字段用于上报支付的场景信息,针对H5支付有以下三种场景,请根据对应场景上报,H5支付不建议在APP端使用，针对场景1，2请接入APP支付，不然可能会出现兼容性问题
	//
	//1，IOS移动应用
	//{"h5_info": //h5支付固定传"h5_info"
	//    {"type": "",  //场景类型
	//     "app_name": "",  //应用名
	//     "bundle_id": ""  //bundle_id
	//     }
	//}
	//
	//2，安卓移动应用
	//{"h5_info": //h5支付固定传"h5_info"
	//    {"type": "",  //场景类型
	//     "app_name": "",  //应用名
	//     "package_name": ""  //包名
	//     }
	//}
	//
	//3，WAP网站应用
	//{"h5_info": //h5支付固定传"h5_info"
	//   {"type": "",  //场景类型
	//    "wap_url": "",//WAP网站URL地址
	//    "wap_name": ""  //WAP 网站名
	//    }
	//}
	//==================================================================================================================
	SceneInfo string `MWEB:"true"`
}

type ReturnUnifiedOrder struct {
	TradeType string
	PrepayId  string
	MwebUrl   string
	CodeUrl   string
}

//查询订单返回数据
type OrderQuery struct {
	DeviceInfo         string //设备号
	OpenId             string //用户标识
	IsSubscribe        string //是否关注公众账号
	TradeType          string //交易类型
	TradeState         string //交易状态
	BankType           string //付款银行
	TotalFee           int    //标价金额
	SettlementTotalFee int    //应结订单金额
	FeeType            string //标价币种
	CashFee            string //现金支付金额
	CashFeeType        string //现金支付币种
	CouponFee          int    //代金券金额
	CouponCount        int    //代金券使用数量
	CouponType         string //代金券类型
	CouponId           string //代金券ID
	CouponFeen         int    //单个代金券支付金额
	TransactionId      string //微信支付订单号
	OutTradeNo         string //商户订单号
	Attach             string //附加数据
	TimeEnd            string //支付完成时间
	TradeStateDesc     string //交易状态描述
}

//申请退款传入数据
type Refund struct {
	TransactionId string //微信支付订单号
	OutTradeNo    string //商户订单号
	OutRefundNo   string `require:"true"` //商户退款单号
	TotalFee      int    `require:"true"` //订单金额
	RefundFee     int    `require:"true"` //退款金额
	RefundFeeType string //退款货币种类
	RefundDesc    string //退款原因
	RefundAccount string //退款资金来源
	NotifyUrl     string //退款结果通知url
	CertP12       []byte `require:"true"`
}

//申请退款返回数据
type RefundReturn struct {
	TransactionId       string //微信支付订单号
	OutTradeNo          string //商户订单号
	OutRefundNo         string //商户退款单号
	RefundId            string //微信退款单号
	RefundFee           int    //退款金额
	SettlementRefundFee int    //应结退款金额
	TotalFee            int    //标价金额
	SettlementTotalFee  int    //应结订单金额
	FeeType             string //标价币种
	CashFee             int    //现金支付金额
	CashFeeType         string //现金支付币种
	CashRefundFee       int    //现金退款金额
	CouponType          string //代金券类型
	CouponRefundFee     int    //代金券退款总金额
	CouponRefundFeen    int    //单个代金券退款金额
	CouponRefundCount   int    //退款代金券使用数量
	CouponRefundId      string //退款代金券ID
}

//退款查询传入数据
type RefundQuery struct {
	TransactionId string //微信支付订单号
	OutTradeNo    string //商户订单号
	OutRefundNo   string //商户退款单号
	RefundId      string //微信退款单号
	Offset        int    //偏移量
}

//退款查询返回数据
type RefundQueryReturn struct {
	TotalRefundCount     int    //订单总退款次数
	TransactionId        string //微信订单号
	OutTradeNo           string //商户订单号
	TotalFee             int    //订单金额
	SettlementTotalFee   int    //应结订单金额
	FeeType              string //标价币种
	CashFee              int    //现金支付金额
	RefundCount          int    //退款笔数
	OutRefundNon         string //商户退款单号
	RefundIdn            string //微信退款单号
	RefundChanneln       string //退款渠道
	RefundFeen           int    //申请退款金额
	SettlementRefundFeen int    //退款金额
	CouponTypenm         string //代金券类型
	CouponRefundFeen     int    //总代金券退款金额
	CouponRefundCountn   int    //退款代金券使用数量
	CouponRefundIdnm     string //退款代金券ID
	CouponRefundFeenm    int    //单个代金券退款金额
	RefundStatusn        string //退款状态
	RefundAccountn       string //退款资金来源
	RefundRecvAccoutn    string //退款入账账户
	RefundSuccessTimen   string //退款成功时间
}

//下载对账单传入数据
type DownLoadBill struct {
	//对账单日期：下载对账单的日期，格式：20140603
	BillDate string `require:"true"`

	//账单类型：
	// ALL（默认值），返回当日所有订单信息（不含充值退款订单）
	//
	//SUCCESS，返回当日成功支付的订单（不含充值退款订单）
	//
	//REFUND，返回当日退款订单（不含充值退款订单）
	//
	//RECHARGE_REFUND，返回当日充值退款订单
	BillType string

	//压缩账单：非必传参数，固定值：GZIP，返回格式为.gzip的压缩包账单。
	// 不传则默认为数据流形式。
	TarType string
}

//账单
type DownloadBillReturn struct {
	TradeDate            string //交易时间
	AppId                string //公众账号ID
	MchId                string //商户号
	SubMchId             string //子商户号
	DeviceInfo           string //设备号
	TransactionId        string //微信订单号
	OutTradeNo           string //商户订单号
	OpenId               string //用户标识
	TradeType            string //交易类型
	TradeState           string //交易状态
	BankType             string //付款银行
	FeeType              string //标价币种
	TotalFee             string //总金额
	CouponFee            string //代金券或立减优惠金额
	RefundApplyTimen     string //退款申请时间
	RefundSuccessTimen   string //退款成功时间
	RefundIdn            string //微信退款单号
	OutRefundNon         string //商户退款单号
	SettlementRefundFeen string //退款金额
	RefundCouponFee      string //代金券或立减优惠退款金额
	RefundType           string //退款类型
	RefundState          string //退款状态
	ProductName          string //商品名称
	ProductBar           string //商品数据包
	HandlingFee          string //手续费
	Rate                 string //费率
}

//账单统计数据
type DownloadBillStatisticsReturn struct {
	TotalTransactions      string //总交易单数
	TotalTransactionValue  string //总交易额
	TotalRefundValue       string //总退款金额
	TotalCouponRefundValue string //总代金券或立减优惠退款金额
	TotalHandlingFeeValue  string //手续费总金额
}
