package alipay

import "time"

//PagePay ...
type PagePay struct {
	OutTradeNo     string `require:"true"` //商户订单号
	ProductCode    string `require:"true"` //销售产品码，与支付宝签约的产品码名称。固定：FAST_INSTANT_TRADE_PAY
	TotalAmount    string `require:"true"` //订单总金额
	Subject        string `require:"true"` //订单标题
	Body           string //订单描述
	GoodsType      string //商品主类型 :0-虚拟类商品,1-实物类商品
	NotifyURL      string
	AppAuthToken   string
	PassbackParams string
	PayExpire      time.Time
}

//PagePayReturn ...
type PagePayReturn struct {
	TradeNo         string //支付宝交易号
	OutTradeNo      string //商户订单号
	SellerID        string //收款支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字
	TotalAmount     string //交易金额
	MerchantOrderNo string //商户原始订单号
	PagePayHTML     []byte //支付宝支付页面
}

//TradeQuery ...
type TradeQuery struct {
	OutTradeNo   string //商家订单号
	TradeNo      string //支付宝交易号
	AppAuthToken string
}

//TradeQueryReturn ...
type TradeQueryReturn struct {
	TradeNo      string `json:"trade_no"`       //支付宝交易号
	OutTradeNo   string `json:"out_trade_no"`   //商家订单号
	BuyerlogonID string `json:"buyer_logon_id"` //买家支付宝账号
	TradeStatus  string `json:"trade_status"`   //交易状态：WAIT_BUYER_PAY（交易创建，等待买家付款）、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）、TRADE_SUCCESS（交易支付成功）、TRADE_FINISHED（交易结束，不可退款）
	TotalAmount  string `json:"total_amount"`   //交易的订单金额，单位为元，两位小数。
	//标价币种，该参数的值为支付时传入的trans_currency，
	// 支持英镑：GBP、港币：HKD、美元：USD、新加坡元：SGD、日元：JPY、
	// 加拿大元：CAD、澳元：AUD、欧元：EUR、新西兰元：NZD、韩元：KRW、
	// 泰铢：THB、瑞士法郎：CHF、瑞典克朗：SEK、丹麦克朗：DKK、
	// 挪威克朗：NOK、马来西亚林吉特：MYR、印尼卢比：IDR、
	// 菲律宾比索：PHP、毛里求斯卢比：MUR、以色列新谢克尔：ILS、
	// 斯里兰卡卢比：LKR、俄罗斯卢布：RUB、阿联酋迪拉姆：AED、
	// 捷克克朗：CZK、南非兰特：ZAR、人民币：CNY、新台币：TWD。
	// 当trans_currency 和 settle_currency 不一致时，
	// trans_currency支持人民币：CNY、新台币：TWD
	TransCurrency       string        `json:"trans_currency"`
	SettleCurrency      string        `json:"settle_currency"`        //订单结算币种
	SettleAmount        string        `json:"settle_amount"`          //结算币种订单金额
	PayCurrency         string        `json:"pay_currency"`           //订单支付币种
	PayAmount           string        `json:"pay_amount"`             //支付币种订单金额
	SettleTransRate     string        `json:"settle_trans_rate"`      //结算币种兑换标价币种汇率
	TransPayRate        string        `json:"trans_pay_rate"`         //标价币种兑换支付币种汇率
	BuyerPayAmount      string        `json:"buyer_pay_amount"`       //买家实付金额，单位为元，两位小数。
	PointAmount         string        `json:"point_amount"`           //积分支付的金额，单位为元，两位小数。
	InvoiceAmount       string        `json:"invoice_amount"`         //交易中用户支付的可开具发票的金额，单位为元，两位小数。
	SendPayDate         string        `json:"send_pay_date"`          //本次交易打款给卖家的时间
	ReceiptAmount       string        `json:"receipt_amount"`         //实收金额，单位为元，两位小数。
	StoreID             string        `json:"store_id"`               //商户门店编号
	TerminalID          string        `json:"terminal_id"`            //商户机具终端编号
	FundBillList        TradeFundBill `json:"fund_bill_list"`         //交易支付使用的资金渠道
	StoreName           string        `json:"store_name"`             //请求交易支付中的商户店铺的名称
	BuyerUserID         string        `json:"buyer_user_id"`          //买家在支付宝的用户id
	ChargeAmount        string        `json:"charge_amount"`          //该笔交易针对收款方的收费金额；
	ChargeFlags         string        `json:"charge_flags"`           //费率活动标识，当交易享受活动优惠费率时，返回该活动的标识；
	SettlementID        string        `json:"settlement_id"`          //支付清算编号，用于清算对账使用；
	AuthTradePayMode    string        `json:"auth_trade_pay_mode"`    //预授权支付模式，该参数仅在信用预授权支付场景下返回。信用预授权支付：CREDIT_PREAUTH_PAY
	BuyerUserType       string        `json:"buyer_user_type"`        //买家用户类型。CORPORATE:企业用户；PRIVATE:个人用户。
	MdiscountAmount     string        `json:"mdiscount_amount"`       //商家优惠金额
	DiscountAmount      string        `json:"discount_amount"`        //平台优惠金额
	BuyerUserName       string        `json:"buyer_user_name"`        //买家名称；买家为个人用户时为买家姓名，买家为企业用户时为企业名称；
	Subject             string        `json:"subject"`                //订单标题；
	Body                string        `json:"body"`                   //订单描述;
	AlipaySubMerchantID string        `json:"alipay_sub_merchant_id"` //间连商户在支付宝端的商户编号；
	ExtInfos            string        `json:"ext_infos"`              //交易额外信息，特殊场景下与支付宝约定返回。
}

//Refund ...
type Refund struct {
	AppAuthToken string
	OutTradeNo   string //商户订单号
	TradeNo      string //支付宝交易号
	RefundAmount string `require:"true"` //需要退款的金额，该金额不能大于订单金额,单位为元
	RefundReason string `require:"true"` //退款的原因说明
	OutRequestNo string `require:"true"` //标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	OperatorID   string //商户的操作员编号
	StoreID      string //商户的门店编号
	TerminalID   string //商户的终端编号
}

//RefundReturn ...
type RefundReturn struct {
	TradeNo                      string            `json:"trade_no"`                        //支付宝交易号
	OutTradeNo                   string            `json:"out_trade_no"`                    //商户订单号
	BuyerLogonID                 string            `json:"buyer_logon_id"`                  //用户的登录id
	FundChange                   string            `json:"fund_change"`                     //本次退款是否发生了资金变化
	RefundFee                    string            `json:"refund_fee"`                      //退款总金额
	RefundCurrency               string            `json:"refund_currency"`                 //订单退款币种信息
	GmtRefundPay                 string            `json:"gmt_refund_pay"`                  //退款支付时间
	RefundDetailItemList         TradeFundBill     `json:"refund_detail_item_list"`         //退款使用的资金渠道
	StoreName                    string            `json:"store_name"`                      //交易在支付时候的门店名称
	BuyerUserID                  string            `json:"buyer_user_id"`                   //买家在支付宝的用户id
	RefundPresetPayToolInfo      PresetPayToolInfo `json:"refund_preset_paytool_list"`      //退回的前置资产列表
	RefundSettlementID           string            `json:"refund_settlement_id"`            //退款清算编号，用于清算对账使用；
	PresentRefundBuyerAmount     string            `json:"present_refund_buyer_amount"`     //本次退款金额中买家退款金额
	PresentRefundDiscountAmount  string            `json:"present_refund_discount_amount"`  //本次退款金额中平台优惠退款金额
	PresentRefundMdiscountAmount string            `json:"present_refund_mdiscount_amount"` //本次退款金额中商家优惠退款金额
}

//RefundQuery ...
type RefundQuery struct {
	AppAuthToken string
	TradeNo      string //支付宝交易号
	OutTradeNo   string //商户订单号
	OutRequestNo string `require:"true"` //退款请求号
}

//RefundQueryReturn ...
type RefundQueryReturn struct {
	TradeNo                      string        `json:"trade_no"`                        //支付宝交易号
	OutTradeNo                   string        `json:"out_trade_no"`                    //创建交易传入的商户订单号
	OutRequestNo                 string        `json:"out_request_no"`                  //本笔退款对应的退款请求号
	RefundReason                 string        `json:"refund_reason"`                   //发起退款时，传入的退款原因
	TotalAmount                  string        `json:"total_amount"`                    //该笔退款所对应的交易的订单金额
	RefundAmount                 string        `json:"refund_amount"`                   //本次退款请求，对应的退款金额
	GmtRefundPay                 string        `json:"gmt_refund_pay"`                  //退款时间；
	RefundDetailItemList         TradeFundBill `json:"refund_detail_item_list"`         //本次退款使用的资金渠道；
	SendBackFee                  string        `json:"send_back_fee"`                   //本次商户实际退回金额；
	RefundSettlementID           string        `json:"refund_settlement_id"`            //退款清算编号，用于清算对账使用；只在银行间联交易场景下返回该信息；
	PresentRefundBuyerAmount     string        `json:"present_refund_buyer_amount"`     //本次退款金额中买家退款金额
	PresentRefundDiscountAmount  string        `json:"present_refund_discount_amount"`  //本次退款金额中平台优惠退款金额
	PresentRefundMdiscountAmount string        `json:"present_refund_mdiscount_amount"` //本次退款金额中商家优惠退款金额
}

//TradeClose ...
type TradeClose struct {
	AppAuthToken string
	TradeNo      string `json:"trade_no"`     //该交易在支付宝系统中的交易流水号。
	OutTradeNo   string `json:"out_trade_no"` //订单支付时传入的商户订单号,和支付宝交易号不能同时为空。
	OperatorID   string `json:"operator_id"`  //卖家端自定义的的操作员 ID
}

//TradeCloseReturn ...
type TradeCloseReturn struct {
	TradeNo    string `json:"trade_no"`     //支付宝交易号
	OutTradeNo string `json:"out_trade_no"` //创建交易传入的商户订单号
}

//BillDownloadURLQuery ...
type BillDownloadURLQuery struct {
	AppAuthToken string
	BillType     string //账单类型，商户通过接口或商户经开放平台授权后其所属服务商通过接口可以获取以下账单类型：trade、signcustomer；trade指商户基于支付宝交易收单的业务账单；signcustomer是指基于商户支付宝余额收入及支出等资金变动的帐务账单。
	BillDate     string //账单时间：日账单格式为yyyy-MM-dd，最早可下载2016年1月1日开始的日账单；月账单格式为yyyy-MM，最早可下载2016年1月开始的月账单。
}

//BillDownloadURLQueryReturn ...
type BillDownloadURLQueryReturn struct {
	BillDownloadURL string `json:"bill_download_url"` //账单下载地址链接，获取连接后30秒后未下载，链接地址失效.
}

//PayNotify ...
type PayNotify struct {
	NotifyTime string `json:"notify_time"` //通知时间
	NotifyType string `json:"notify_type"` //通知类型
	NotifyID   string `json:"notify_id"`   //通知校验id
	Charset    string `json:"charset"`     //编码格式
	Version    string `json:"version"`     //调用接口版本
	SignType   string `json:"sign_type"`   //签名类型
	Sign       string `json:"sign"`        //签名
	AuthAppID  string `json:"auth_app_id"` //授权方的app_id
	PayBusinessNotify
}

//PayBusinessNotify ...
type PayBusinessNotify struct {
	TradeNo           string  `json:"trade_no"`            //支付宝交易号
	AppID             string  `json:"app_id"`              //开发者的app_id
	OutTradeNo        string  `json:"out_trade_no"`        //商户订单号
	OutBizNo          string  `json:"out_biz_no"`          //商户业务号
	BuyerID           string  `json:"buyer_id"`            //买家支付宝用户号
	SellerID          string  `json:"seller_id"`           //卖家支付宝用户号
	TradeStatus       string  `json:"trade_status"`        //交易状态
	TotalAmount       float64 `json:"total_amount"`        //订单金额
	ReceiptAmount     float64 `json:"receipt_amount"`      //实收金额
	InvoiceAmount     float64 `json:"invoice_amount"`      //开票金额
	BuyerPayAmount    float64 `json:"buyer_pay_amount"`    //付款金额
	PointAmount       float64 `json:"point_amount"`        //集分宝金额
	RefundFee         float64 `json:"refund_fee"`          //总退款金额
	Subject           string  `json:"subject"`             //订单标题
	Body              string  `json:"body"`                //商品描述
	GmtCreate         string  `json:"gmt_create"`          //交易创建时间
	GmtPayment        string  `json:"gmt_payment"`         //交易付款时间
	GmtRefund         string  `json:"gmt_refund"`          //交易退款时间
	GmtClose          string  `json:"gmt_close"`           //交易结束时间
	FundBillList      string  `json:"fund_bill_list"`      //支付金额信息
	VoucherDetailList string  `json:"voucher_detail_list"` //优惠券信息
	PassbackParams    string  `json:"passback_params"`     //回传参数
}
