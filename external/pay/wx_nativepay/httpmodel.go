package wx_nativepay

type BizPayUrl struct {
	AppId     string `xml:"appid"`      //公众账号ID
	MchId     string `xml:"mch_id"`     //商户号
	TimeStamp string `xml:"time_stamp"` //时间戳
	NonceStr  string `xml:"nonce_str"`  //随机字符串
	ProductId string `xml:"product_id"` //商品ID
	Sign      string `xml:"sign"`       //签名
}

type ReqUnifiedOrder struct {
	AppId          string `xml:"appid"`            //公众账号ID
	MchId          string `xml:"mch_id"`           //商户号
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
	SpbillCreateIp string `xml:"spbill_create_ip"` //终端IP
	TimeStart      string `xml:"time_start"`       //交易起始时间
	TimeExpire     string `xml:"time_expire"`      //交易结束时间
	GoodsTag       string `xml:"goods_tag"`        //订单优惠标记
	NotifyUrl      string `xml:"notify_url"`       //通知地址
	TradeType      string `xml:"trade_type"`       //交易类型
	ProductId      string `xml:"product_id"`       //商品ID
	LimitPay       string `xml:"limit_pay"`        //指定支付方式
	OpenId         string `xml:"openid"`           //用户标识
	Receipt        string `xml:"receipt"`          //电子发票入口开放标识
	SceneInfo      string `xml:"scene_info"`       //场景信息
}

type RespUnifiedOrder struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppId       string `xml:"appid"`         //公众账号ID
	MchId       string `xml:"mch_id"`        //商户号
	DeviceInfo  string `xml:"device_info"`   //设备号
	NonceStr    string `xml:"nonce_str"`     //随机字符串
	Sign        string `xml:"sign"`          //签名
	ResultCode  string `xml:"result_code"`   //业务结果
	ErrCode     string `xml:"err_code"`      //错误代码
	ErrCodeDesc string `xml:"err_code_desc"` //错误代码描述

	//以下字段在return_code 和result_code都为SUCCESS的时候有返回
	TradeType string `xml:"trade_type"` //交易类型
	PrepayId  string `xml:"prepay_id"`  //预支付交易会话标识
	CodeUrl   string `xml:"code_url"`   //二维码链接
}

type ReqOrderQuery struct {
	AppId         string `xml:"appid"`          //公众账号ID
	MchId         string `xml:"mch_id"`         //商户号
	TransactionId string `xml:"transaction_id"` //微信订单号
	OutTradeNo    string `xml:"out_trade_no"`   //商户订单号
	NonceStr      string `xml:"nonce_str"`      //随机字符串
	Sign          string `xml:"sign"`           //签名
	SignType      string `xml:"sign_type"`      //签名类型
}

type RespOrderQuery struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppId       string `xml:"appid"`         //公众账号ID
	MchId       string `xml:"mch_id"`        //商户号
	NonceStr    string `xml:"nonce_str"`     //随机字符串
	Sign        string `xml:"sign"`          //签名
	ResultCode  string `xml:"result_code"`   //业务结果
	ErrCode     string `xml:"err_code"`      //错误代码
	ErrCodeDesc string `xml:"err_code_desc"` //错误代码描述

	//以下字段在return_code 、result_code、trade_state都为SUCCESS时有返回 ，
	//如trade_state不为 SUCCESS，则只返回out_trade_no（必传）和attach（选传）。
	DeviceInfo         string `xml:"device_info"`          //设备号
	OpenId             string `xml:"openid"`               //用户标识
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
	CouponId           string `xml:"coupon_id_$n"`         //代金券ID
	CouponFeen         int    `xml:"coupon_fee_$n"`        //单个代金券支付金额
	TransactionId      string `xml:"transaction_id"`       //微信支付订单号
	OutTradeNo         string `xml:"out_trade_no"`         //商户订单号
	Attach             string `xml:"attach"`               //附加数据
	TimeEnd            string `xml:"time_end"`             //支付完成时间
	TradeStateDesc     string `xml:"trade_state_desc"`     //交易状态描述
}

type ReqCloseOrder struct {
	AppId      string `xml:"appid"`        //公众账号ID
	MchId      string `xml:"mch_id"`       //商户号
	OutTradeNo string `xml:"out_trade_no"` //商户订单号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	Sign       string `xml:"sign"`         //签名
	SignType   string `xml:"sign_type"`    //签名类型
}

type RespCloseOrder struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppId       string `xml:"appid"`         //公众账号ID
	MchId       string `xml:"mch_id"`        //商户号
	NonceStr    string `xml:"nonce_str"`     //随机字符串
	Sign        string `xml:"sign"`          //签名
	ResultCode  string `xml:"result_code"`   //业务结果
	ResultMsg   string `xml:"result_msg"`    //业务结果描述
	ErrCode     string `xml:"err_code"`      //错误代码
	ErrCodeDesc string `xml:"err_code_desc"` //错误代码描述
}

type ReqRefund struct {
	AppId         string `xml:"appid"`           //公众账号ID
	MchId         string `xml:"mch_id"`          //商户号
	NonceStr      string `xml:"nonce_str"`       //随机字符串
	Sign          string `xml:"sign"`            //签名
	SignType      string `xml:"sign_type"`       //签名类型
	TransactionId string `xml:"transaction_id"`  //微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`    //商户订单号
	OutRefundNo   string `xml:"out_refund_no"`   //商户退款单号
	TotalFee      string `xml:"total_fee"`       //订单金额
	RefundFee     string `xml:"refund_fee"`      //退款金额
	RefundFeeType string `xml:"refund_fee_type"` //退款货币种类
	RefundDesc    string `xml:"refund_desc"`     //退款原因
	RefundAccount string `xml:"refund_account"`  //退款资金来源
	NotifyUrl     string `xml:"notify_url"`      //退款结果通知url
}

type RespRefund struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	ResultCode          string `xml:"result_code"`           //业务结果
	ErrCode             string `xml:"err_code"`              //错误代码
	ErrCodeDesc         string `xml:"err_code_desc"`         //错误代码描述
	AppId               string `xml:"appid"`                 //公众账号ID
	MchId               string `xml:"mch_id"`                //商户号
	NonceStr            string `xml:"nonce_str"`             //随机字符串
	Sign                string `xml:"sign"`                  //签名
	SignType            string `xml:"sign_type"`             //签名类型
	TransactionId       string `xml:"transaction_id"`        //微信支付订单号
	OutTradeNo          string `xml:"out_trade_no"`          //商户订单号
	OutRefundNo         string `xml:"out_refund_no"`         //商户退款单号
	RefundId            string `xml:"refund_id"`             //微信退款单号
	RefundFee           int    `xml:"refund_fee"`            //退款金额
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
	CouponRefundId      string `xml:"coupon_refund_id_$n"`   //退款代金券ID
}

type ReqRefundQuery struct {
	AppId         string `xml:"appid"`          //公众账号ID
	MchId         string `xml:"mch_id"`         //商户号
	NonceStr      string `xml:"nonce_str"`      //随机字符串
	Sign          string `xml:"sign"`           //签名
	SignType      string `xml:"sign_type"`      //签名类型
	TransactionId string `xml:"transaction_id"` //微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`   //商户订单号
	OutRefundNo   string `xml:"out_refund_no"`  //商户退款单号
	RefundId      string `xml:"refund_id"`      //微信退款单号
	Offset        int    `xml:"offset"`         //偏移量
}

type RespRefundQuery struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	ResultCode           string `xml:"result_code"`              //业务结果
	ErrCode              string `xml:"err_code"`                 //错误代码
	ErrCodeDesc          string `xml:"err_code_desc"`            //错误代码描述
	AppId                string `xml:"appid"`                    //公众账号ID
	MchId                string `xml:"mch_id"`                   //商户号
	NonceStr             string `xml:"nonce_str"`                //随机字符串
	Sign                 string `xml:"sign"`                     //签名
	TotalRefundCount     int    `xml:"total_refund_count"`       //订单总退款次数
	TransactionId        string `xml:"transaction_id	"`       //微信订单号
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

type ReqDownloadBill struct {
	AppId    string `xml:"appid"`     //公众账号ID
	MchId    string `xml:"mch_id"`    //商户号
	NonceStr string `xml:"nonce_str"` //随机字符串
	Sign     string `xml:"sign"`      //签名
	BillDate string `xml:"bill_date"` //对账单日期
	BillType string `xml:"bill_type"` //账单类型
	TarType  string `xml:"tar_type"`  //压缩账单
}

type RespDownloadBill struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息
	ErrCode    string `xml:"err_code"`    //错误代码
}

type ReqDownloadFundFlow struct {
	//TODO
}

type RespDownloadFundFlow struct {
	//TODO
}

type ReqPrepayNotify struct {
	AppId       string `xml:"appid"`        //公众账号ID
	MchId       string `xml:"mch_id"`       //商户号
	OpenId      string `xml:"openid"`       //用户标识
	IsSubscribe string `xml:"is_subscribe"` //是否关注公众账号
	NonceStr    string `xml:"nonce_str"`    //随机字符串
	ProductId   string `xml:"product_id"`   //商品ID
	Sign        string `xml:"sign"`         //签名
}

type RespPrepayNotify struct {
	ReturnCode string `xml:"return_code"`  //返回状态码
	ReturnMsg  string `xml:"return_msg"`   //返回信息
	AppId      string `xml:"appid"`        //公众账号ID
	MchId      string `xml:"mch_id"`       //商户号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	PrepayId   string `xml:"prepay_id"`    //预支付ID
	ResultCode string `xml:"result_code"`  //业务结果
	ErrCodeDes string `xml:"err_code_des"` //错误描述
	Sign       string `xml:"sign"`         //签名
}

type ReqPayNotify struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppId              string `xml:"appid"`                //公众账号ID
	MchId              string `xml:"mch_id"`               //商户号
	DeviceInfo         string `xml:"device_info"`          //设备号
	NonceStr           string `xml:"nonce_str"`            //随机字符串
	Sign               string `xml:"sign"`                 //签名
	SignType           string `xml:"sign_type"`            //签名类型
	ResultCode         string `xml:"result_code"`          //业务结果
	ErrCode            string `xml:"err_code"`             //错误代码
	ErrCodeDesc        string `xml:"err_code_desc"`        //错误代码描述
	OpenId             string `xml:"openid"`               //用户标识
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
	CouponId           string `xml:"coupon_id_$n"`         //代金券ID
	CouponFeen         int    `xml:"coupon_fee_$n"`        //单个代金券支付金额
	TransactionId      string `xml:"transaction_id"`       //微信订单号
	OutTradeNo         string `xml:"out_trade_no"`         //商户订单号
	Attach             string `xml:"attach"`               //附加数据
	TimeEnd            string `xml:"time_end"`             //支付完成时间
}

type RespPayNotify struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息
}

type ReqPayitilReport struct {
	//TODO
}

type RespPayitilReport struct {
	//TODO
}

type ReqShortUrl struct {
	AppId    string `xml:"appid"`     //公众账号ID
	MchId    string `xml:"mch_id"`    //商户号
	LongUrl  string `xml:"long_url"`  //长URL链接
	NonceStr string `xml:"nonce_str"` //随机字符串
	Sign     string `xml:"sign"`      //签名
	SignType string `xml:"sign_type"` //签名类型
}

type RespShortUrl struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppId      string `xml:"appid"`       //公众账号ID
	MchId      string `xml:"mch_id"`      //商户号
	NonceStr   string `xml:"nonce_str"`   //随机字符串
	Sign       string `xml:"sign"`        //签名
	ResultCode string `xml:"result_code"` //业务结果
	ErrCode    string `xml:"err_code"`    //错误代码
	ShortUrl   string `xml:"short_url"`   //短URL链接
}

type ReqRefundNotify struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppId    string `xml:"appid"`     //公众账号ID
	MchId    string `xml:"mch_id"`    //退款的商户号
	NonceStr string `xml:"nonce_str"` //随机字符串
	ReqInfo  string `xml:"req_info"`  //加密信息
}

type ReqRefundNotifyEncryptInfo struct {
	TransactionId       string `xml:"transaction_id"`        //微信订单号
	OutTradeNo          string `xml:"out_trade_no"`          //商户订单号
	RefundId            string `xml:"refund_id"`             //微信退款单号
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

type RespRefundNotify struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息
}

type ReqBatchQueryComment struct {
	//TODO
}

type RespBatchQueryComment struct {
	//TODO
}
