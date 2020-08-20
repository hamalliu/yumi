package wxpay

import "encoding/xml"

//BizPayURL ...
type BizPayURL struct {
	XMLName xml.Name `xml:"xml"`

	AppID     string `xml:"appid"`      //公众账号ID
	MchID     string `xml:"mch_id"`     //商户号
	TimeStamp string `xml:"time_stamp"` //时间戳
	NonceStr  string `xml:"nonce_str"`  //随机字符串
	ProductID string `xml:"product_id"` //商品ID
	Sign      string `xml:"sign"`       //签名
}

//ReqUnifiedOrder ...
type ReqUnifiedOrder struct {
	XMLName xml.Name `xml:"xml"`

	AppID          string `xml:"appid"`            //公众账号ID
	MchID          string `xml:"mch_id"`           //商户号
	DeviceInfo     string `xml:"device_info"`      //设备号
	NonceStr       string `xml:"nonce_str"`        //随机字符串
	Sign           string `xml:"sign"`             //签名
	SignType       string `xml:"sign_type"`        //签名类型
	Body           string `xml:"body"`             //商品描述
	Detail         string `xml:"detail"`           //商品详情
	Attach         string `xml:"attach"`           //附加数据
	OutTradeNo     string `xml:"out_trade_no"`     //商户订单号
	FeeType        string `xml:"fee_type"`         //标价币种
	TotalFee       int    `xml:"total_fee"`        //标价金额
	SpbillCreateIP string `xml:"spbill_create_ip"` //终端IP
	TimeStart      string `xml:"time_start"`       //交易起始时间
	TimeExpire     string `xml:"time_expire"`      //交易结束时间
	GoodsTag       string `xml:"goods_tag"`        //订单优惠标记
	NotifyURL      string `xml:"notify_url"`       //通知地址
	TradeType      string `xml:"trade_type"`       //交易类型
	ProductID      string `xml:"product_id"`       //商品ID
	LimitPay       string `xml:"limit_pay"`        //指定支付方式
	OpenID         string `xml:"openid"`           //用户标识
	Receipt        string `xml:"receipt"`          //电子发票入口开放标识
	SceneInfo      string `xml:"scene_info"`       //场景信息
}

//RespUnifiedOrder ...
type RespUnifiedOrder struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppID      string `xml:"appid"`        //公众账号ID
	MchID      string `xml:"mch_id"`       //商户号
	DeviceInfo string `xml:"device_info"`  //设备号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	Sign       string `xml:"sign"`         //签名
	ResultCode string `xml:"result_code"`  //业务结果
	ErrCode    string `xml:"err_code"`     //错误代码
	ErrCodeDes string `xml:"err_code_des"` //错误代码描述

	//以下字段在return_code 和result_code都为SUCCESS的时候有返回
	TradeType string `xml:"trade_type"` //交易类型
	PrepayID  string `xml:"prepay_id"`  //预支付交易会话标识
	CodeURL   string `xml:"code_url"`   //二维码链接
	MwebURL   string `xml:"mweb_url"`   //mweb_url为拉起微信支付收银台的中间页面，可通过访问该url来拉起微信客户端，完成支付,mweb_url的有效期为5分钟。
}

//ReqOrderQuery ...
type ReqOrderQuery struct {
	XMLName xml.Name `xml:"xml"`

	AppID         string `xml:"appid"`          //公众账号ID
	MchID         string `xml:"mch_id"`         //商户号
	TransactionID string `xml:"transaction_id"` //微信订单号
	OutTradeNo    string `xml:"out_trade_no"`   //商户订单号
	NonceStr      string `xml:"nonce_str"`      //随机字符串
	Sign          string `xml:"sign"`           //签名
	SignType      string `xml:"sign_type"`      //签名类型
}

//RespOrderQuery ...
type RespOrderQuery struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppID      string `xml:"appid"`        //公众账号ID
	MchID      string `xml:"mch_id"`       //商户号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	Sign       string `xml:"sign"`         //签名
	ResultCode string `xml:"result_code"`  //业务结果
	ErrCode    string `xml:"err_code"`     //错误代码
	ErrCodeDes string `xml:"err_code_des"` //错误代码描述

	//以下字段在return_code 、result_code、trade_state都为SUCCESS时有返回 ，
	//如trade_state不为 SUCCESS，则只返回out_trade_no（必传）和attach（选传）。
	DeviceInfo         string `xml:"device_info"`          //设备号
	OpenID             string `xml:"openid"`               //用户标识
	IsSubscribe        string `xml:"is_subscribe"`         //是否关注公众账号
	TradeType          string `xml:"trade_type"`           //交易类型
	TradeState         string `xml:"trade_state"`          //交易状态
	BankType           string `xml:"bank_type"`            //付款银行
	TotalFee           int    `xml:"total_fee"`            //标价金额
	SettlementTotalFee int    `xml:"settlement_total_fee"` //应结订单金额
	FeeType            string `xml:"fee_type"`             //标价币种
	CashFee            string `xml:"cash_fee"`             //现金支付金额
	CashFeeType        string `xml:"cash_fee_type"`        //现金支付币种
	CouponFee          int    `xml:"coupon_fee"`           //代金券金额
	CouponCount        int    `xml:"coupon_count"`         //代金券使用数量
	CouponType         string `xml:"coupon_type_$n"`       //代金券类型
	CouponID           string `xml:"coupon_id_$n"`         //代金券ID
	CouponFeen         int    `xml:"coupon_fee_$n"`        //单个代金券支付金额
	TransactionID      string `xml:"transaction_id"`       //微信支付订单号
	OutTradeNo         string `xml:"out_trade_no"`         //商户订单号
	Attach             string `xml:"attach"`               //附加数据
	TimeEnd            string `xml:"time_end"`             //支付完成时间
	TradeStateDesc     string `xml:"trade_state_desc"`     //交易状态描述
}

//ReqCloseOrder ...
type ReqCloseOrder struct {
	XMLName xml.Name `xml:"xml"`

	AppID      string `xml:"appid"`        //公众账号ID
	MchID      string `xml:"mch_id"`       //商户号
	OutTradeNo string `xml:"out_trade_no"` //商户订单号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	Sign       string `xml:"sign"`         //签名
	SignType   string `xml:"sign_type"`    //签名类型
}

//RespCloseOrder ...
type RespCloseOrder struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppID      string `xml:"appid"`        //公众账号ID
	MchID      string `xml:"mch_id"`       //商户号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	Sign       string `xml:"sign"`         //签名
	ResultCode string `xml:"result_code"`  //业务结果
	ResultMsg  string `xml:"result_msg"`   //业务结果描述
	ErrCode    string `xml:"err_code"`     //错误代码
	ErrCodeDes string `xml:"err_code_des"` //错误代码描述
}

//ReqRefund ...
type ReqRefund struct {
	XMLName xml.Name `xml:"xml"`

	AppID         string `xml:"appid"`           //公众账号ID
	MchID         string `xml:"mch_id"`          //商户号
	NonceStr      string `xml:"nonce_str"`       //随机字符串
	Sign          string `xml:"sign"`            //签名
	SignType      string `xml:"sign_type"`       //签名类型
	TransactionID string `xml:"transaction_id"`  //微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`    //商户订单号
	OutRefundNo   string `xml:"out_refund_no"`   //商户退款单号
	TotalFee      int    `xml:"total_fee"`       //订单金额
	RefundFee     int    `xml:"refund_fee"`      //退款金额
	RefundFeeType string `xml:"refund_fee_type"` //退款货币种类
	RefundDesc    string `xml:"refund_desc"`     //退款原因
	RefundAccount string `xml:"refund_account"`  //退款资金来源
	NotifyURL     string `xml:"notify_url"`      //退款结果通知url
}

//RespRefund ...
type RespRefund struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	ResultCode          string `xml:"result_code"`           //业务结果
	ErrCode             string `xml:"err_code"`              //错误代码
	ErrCodeDes          string `xml:"err_code_des"`          //错误代码描述
	AppID               string `xml:"appid"`                 //公众账号ID
	MchID               string `xml:"mch_id"`                //商户号
	NonceStr            string `xml:"nonce_str"`             //随机字符串
	Sign                string `xml:"sign"`                  //签名
	TransactionID       string `xml:"transaction_id"`        //微信支付订单号
	OutTradeNo          string `xml:"out_trade_no"`          //商户订单号
	OutRefundNo         string `xml:"out_refund_no"`         //商户退款单号
	RefundID            string `xml:"refund_id"`             //微信退款单号
	RefundFee           int    `xml:"refund_fee"`            //退款金额
	RefundChannel       string `xml:"refund_channel"`        //退款金额
	SettlementRefundFee int    `xml:"settlement_refund_fee"` //应结退款金额
	TotalFee            int    `xml:"total_fee"`             //标价金额
	SettlementTotalFee  int    `xml:"settlement_total_fee"`  //应结订单金额
	FeeType             string `xml:"fee_type"`              //标价币种
	CashFee             int    `xml:"cash_fee"`              //现金支付金额
	CashFeeType         string `xml:"cash_fee_type"`         //现金支付币种
	CashRefundFee       int    `xml:"cash_refund_fee"`       //现金退款金额
	CouponType          string `xml:"coupon_type_$n"`        //代金券类型
	CouponRefundFee     int    `xml:"coupon_refund_fee"`     //代金券退款总金额
	CouponRefundFeen    int    `xml:"coupon_refund_fee_$n"`  //单个代金券退款金额
	CouponRefundCount   int    `xml:"coupon_refund_count"`   //退款代金券使用数量
	CouponRefundID      string `xml:"coupon_refund_id_$n"`   //退款代金券ID
}

//ReqRefundQuery ...
type ReqRefundQuery struct {
	XMLName xml.Name `xml:"xml"`

	AppID         string `xml:"appid"`          //公众账号ID
	MchID         string `xml:"mch_id"`         //商户号
	NonceStr      string `xml:"nonce_str"`      //随机字符串
	Sign          string `xml:"sign"`           //签名
	SignType      string `xml:"sign_type"`      //签名类型
	TransactionID string `xml:"transaction_id"` //微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`   //商户订单号
	OutRefundNo   string `xml:"out_refund_no"`  //商户退款单号
	RefundID      string `xml:"refund_id"`      //微信退款单号
	Offset        int    `xml:"offset,canzero"` //偏移量
}

//RespRefundQuery ...
type RespRefundQuery struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	ResultCode           string `xml:"result_code"`              //业务结果
	ErrCode              string `xml:"err_code"`                 //错误代码
	ErrCodeDes           string `xml:"err_code_des"`             //错误代码描述
	AppID                string `xml:"appid"`                    //公众账号ID
	MchID                string `xml:"mch_id"`                   //商户号
	NonceStr             string `xml:"nonce_str"`                //随机字符串
	Sign                 string `xml:"sign"`                     //签名
	TotalRefundCount     int    `xml:"total_refund_count"`       //订单总退款次数
	TransactionID        string `xml:"transaction_id	"`          //微信订单号
	OutTradeNo           string `xml:"out_trade_no"`             //商户订单号
	TotalFee             int    `xml:"total_fee"`                //订单金额
	SettlementTotalFee   int    `xml:"settlement_total_fee"`     //应结订单金额
	FeeType              string `xml:"fee_type"`                 //标价币种
	CashFee              int    `xml:"cash_fee"`                 //现金支付金额
	RefundCount          int    `xml:"refund_count"`             //退款笔数
	OutRefundNon         string `xml:"out_refund_no_$n"`         //商户退款单号
	RefundIdn            string `xml:"refund_id_$n"`             //微信退款单号
	RefundChanneln       string `xml:"refund_channel_$n"`        //退款渠道
	RefundFeen           int    `xml:"refund_fee_$n"`            //申请退款金额
	SettlementRefundFeen int    `xml:"settlement_refund_fee_$n"` //退款金额
	CouponTypenm         string `xml:"coupon_type_$n_$m"`        //代金券类型
	CouponRefundFeen     int    `xml:"coupon_refund_fee_$n"`     //总代金券退款金额
	CouponRefundCountn   int    `xml:"coupon_refund_count_$n"`   //退款代金券使用数量
	CouponRefundIdnm     string `xml:"coupon_refund_id_$n_$m"`   //退款代金券ID
	CouponRefundFeenm    int    `xml:"coupon_refund_fee_$n_$m"`  //单个代金券退款金额
	RefundStatusn        string `xml:"refund_status_$n"`         //退款状态
	RefundAccountn       string `xml:"refund_account_$n"`        //退款资金来源
	RefundRecvAccoutn    string `xml:"refund_recv_accout_$n"`    //退款入账账户
	RefundSuccessTimen   string `xml:"refund_success_time_$n"`   //退款成功时间
}

//ReqDownloadBill ...
type ReqDownloadBill struct {
	XMLName xml.Name `xml:"xml"`

	AppID    string `xml:"appid"`     //公众账号ID
	MchID    string `xml:"mch_id"`    //商户号
	NonceStr string `xml:"nonce_str"` //随机字符串
	Sign     string `xml:"sign"`      //签名
	BillDate string `xml:"bill_date"` //对账单日期
	BillType string `xml:"bill_type"` //账单类型
	TarType  string `xml:"tar_type"`  //压缩账单
}

//RespDownloadBill ...
type RespDownloadBill struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息
	ErrCode    string `xml:"err_code"`    //错误代码
}

//ReqDownloadFundFlow ...
type ReqDownloadFundFlow struct {
	//TODO
}

//RespDownloadFundFlow ...
type RespDownloadFundFlow struct {
	//TODO
}

//ReqPrepayNotify ...
type ReqPrepayNotify struct {
	XMLName xml.Name `xml:"xml"`

	AppID       string `xml:"appid"`        //公众账号ID
	MchID       string `xml:"mch_id"`       //商户号
	OpenID      string `xml:"openid"`       //用户标识
	IsSubscribe string `xml:"is_subscribe"` //是否关注公众账号
	NonceStr    string `xml:"nonce_str"`    //随机字符串
	ProductID   string `xml:"product_id"`   //商品ID
	Sign        string `xml:"sign"`         //签名
}

//RespPrepayNotify ...
type RespPrepayNotify struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`  //返回状态码
	ReturnMsg  string `xml:"return_msg"`   //返回信息
	AppID      string `xml:"appid"`        //公众账号ID
	MchID      string `xml:"mch_id"`       //商户号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	PrepayID   string `xml:"prepay_id"`    //预支付ID
	ResultCode string `xml:"result_code"`  //业务结果
	ErrCodeDes string `xml:"err_code_des"` //错误描述
	Sign       string `xml:"sign"`         //签名
}

//ReqPayNotify ...
type ReqPayNotify struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppID              string `xml:"appid"`                //公众账号ID
	MchID              string `xml:"mch_id"`               //商户号
	DeviceInfo         string `xml:"device_info"`          //设备号
	NonceStr           string `xml:"nonce_str"`            //随机字符串
	Sign               string `xml:"sign"`                 //签名
	SignType           string `xml:"sign_type"`            //签名类型
	ResultCode         string `xml:"result_code"`          //业务结果
	ErrCode            string `xml:"err_code"`             //错误代码
	ErrCodeDes         string `xml:"err_code_des"`         //错误代码描述
	OpenID             string `xml:"openid"`               //用户标识
	IsSubscribe        string `xml:"is_subscribe"`         //是否关注公众账号
	TradeType          string `xml:"trade_type"`           //交易类型
	BankType           string `xml:"bank_type"`            //付款银行
	TotalFee           int    `xml:"total_fee"`            //订单金额
	SettlementTotalFee int    `xml:"settlement_total_fee"` //应结订单金额
	FreeType           string `xml:"fee_type"`             //货币种类
	CashFee            string `xml:"cash_fee"`             //现金支付金额
	CashFeeType        string `xml:"cash_fee_type"`        //现金支付币种
	CouponFee          int    `xml:"coupon_fee"`           //总代金券金额
	CouponCount        int    `xml:"coupon_count"`         //代金券使用数量
	CouponType         string `xml:"coupon_type_$n"`       //代金券类型
	CouponID           string `xml:"coupon_id_$n"`         //代金券ID
	CouponFeen         int    `xml:"coupon_fee_$n"`        //单个代金券支付金额
	TransactionID      string `xml:"transaction_id"`       //微信订单号
	OutTradeNo         string `xml:"out_trade_no"`         //商户订单号
	Attach             string `xml:"attach"`               //附加数据
	TimeEnd            string `xml:"time_end"`             //支付完成时间
}

//RespPayNotify ...
type RespPayNotify struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息
}

//ReqPayitilReport ...
type ReqPayitilReport struct {
	//TODO
}

//RespPayitilReport ...
type RespPayitilReport struct {
	//TODO
}

//ReqShortURL ...
type ReqShortURL struct {
	XMLName xml.Name `xml:"xml"`

	AppID    string `xml:"appid"`     //公众账号ID
	MchID    string `xml:"mch_id"`    //商户号
	LongURL  string `xml:"long_url"`  //长URL链接
	NonceStr string `xml:"nonce_str"` //随机字符串
	Sign     string `xml:"sign"`      //签名
	SignType string `xml:"sign_type"` //签名类型
}

//RespShortURL ...
type RespShortURL struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppID      string `xml:"appid"`       //公众账号ID
	MchID      string `xml:"mch_id"`      //商户号
	NonceStr   string `xml:"nonce_str"`   //随机字符串
	Sign       string `xml:"sign"`        //签名
	ResultCode string `xml:"result_code"` //业务结果
	ErrCode    string `xml:"err_code"`    //错误代码
	ShortURL   string `xml:"short_url"`   //短URL链接
}

//ReqRefundNotify ...
type ReqRefundNotify struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppID    string `xml:"appid"`     //公众账号ID
	MchID    string `xml:"mch_id"`    //退款的商户号
	NonceStr string `xml:"nonce_str"` //随机字符串
	ReqInfo  string `xml:"req_info"`  //加密信息

	DecryptReqInfo ReqRefundNotifyEncryptInfo `xml:"-"`
}

//ReqRefundNotifyEncryptInfo ...
type ReqRefundNotifyEncryptInfo struct {
	XMLName xml.Name `xml:"root"`

	TransactionID       string `xml:"transaction_id"`        //微信订单号
	OutTradeNo          string `xml:"out_trade_no"`          //商户订单号
	RefundID            string `xml:"refund_id"`             //微信退款单号
	OutRefundNo         string `xml:"out_refund_no"`         //商户退款单号
	TotalFee            int    `xml:"total_fee"`             //订单金额
	SettlementTotalFee  int    `xml:"settlement_total_fee"`  //应结订单金额
	RefundFee           string `xml:"refund_fee"`            //申请退款金额
	SettlementRefundFee int    `xml:"settlement_refund_fee"` //退款金额
	RefundStatusn       string `xml:"refund_status_$n"`      //退款状态
	SuccessTime         string `xml:"success_time"`          //退款成功时间
	RefundRecvAccout    string `xml:"refund_recv_accout"`    // 退款入账账户
	RefundAccount       string `xml:"refund_account"`        //退款资金来源
	RefundRequestSource string `xml:"refund_request_source"` //退款发起来源
}

//RespRefundNotify ...
type RespRefundNotify struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息
}

//ReqBatchQueryComment ...
type ReqBatchQueryComment struct {
	//TODO
}

//RespBatchQueryComment ...
type RespBatchQueryComment struct {
	//TODO
}
