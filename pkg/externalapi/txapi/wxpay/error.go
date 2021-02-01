package wxpay

import "errors"

var (
	ErrInvalidRequest       = errors.New("参数错误")
	ErrNoAuth               = errors.New("商户无此接口权限，或异常IP请求不予受理")
	ErrNotEnouth            = errors.New("余额不足")
	ErrOrderPaid            = errors.New("商户订单已支付")
	ErrOrderClose           = errors.New("订单已关闭")
	ErrSyetemError          = errors.New("系统错误")
	ErrAppidNotExist        = errors.New("APPID不存在")
	ErrMchidNotExist        = errors.New("MCHID不存在")
	ErrAppidMchidNotMatch   = errors.New("appid和mch_id不匹配")
	ErrLackParams           = errors.New("缺少参数")
	ErrOutTradeNoUsed       = errors.New("商户订单号重复")
	ErrSignError            = errors.New("签名错误")
	ErrXMLFormatError       = errors.New("XML格式错误")
	ErrRequirePostMethod    = errors.New("请使用post方法")
	ErrPostDataEmpty        = errors.New("post数据为空")
	ErrNotUtf8              = errors.New("编码格式错误")
	ErrOrderNotExist        = errors.New("此交易订单号不存在")
	ErrBizerrNeedRetry      = errors.New("退款业务流程错误，需要商户触发重试来解决")
	ErrTradeOverdue         = errors.New("订单已经超过退款期限")
	ErrError                = errors.New("业务错误")
	ErrUserAccountAbnormal  = errors.New("退款请求失败")
	ErrInvalidReqTooMuch    = errors.New("无效请求过多")
	ErrInvalidTransactionID = errors.New("无效transaction_id")
	ErrParamError           = errors.New("参数错误")
	ErrFrequencyLimited     = errors.New("频率限制")
	ErrRefundNotExist       = errors.New("退款订单查询失败")
	ErrNoComment            = errors.New("对应的时间段没有用户的评论数据")
	ErrTimeExpire           = errors.New("拉取的时间超过3个月")
)

var errMap = map[string]error{
	"INVALID_REQUEST":       ErrInvalidRequest,
	"NOAUTH":                ErrNoAuth,
	"NOTENOUGH":             ErrNotEnouth,
	"ORDERPAID":             ErrOrderPaid,
	"ORDERCLOSED":           ErrOrderClose,
	"SYSTEMERROR":           ErrSyetemError,
	"APPID_NOT_EXIST":       ErrAppidNotExist,
	"MCHID_NOT_EXIST":       ErrMchidNotExist,
	"APPID_MCHID_NOT_MATCH": ErrAppidMchidNotMatch,
	"LACK_PARAMS":           ErrLackParams,
	"OUT_TRADE_NO_USED":     ErrOutTradeNoUsed,
	"SIGNERROR":             ErrSignError,
	"XML_FORMAT_ERROR":      ErrXMLFormatError,
	"REQUIRE_POST_METHOD":   ErrRequirePostMethod,
	"POST_DATA_EMPTY":       ErrPostDataEmpty,
	"NOT_UTF8":              ErrNotUtf8,
	"ORDERNOTEXIST":         ErrOrderNotExist,
	"BIZERR_NEED_RETRY":     ErrBizerrNeedRetry,
	"TRADE_OVERDUE":         ErrTradeOverdue,
	"ERROR":                 ErrError,
	"USER_ACCOUNT_ABNORMAL": ErrUserAccountAbnormal,
	"INVALID_REQ_TOO_MUCH":  ErrInvalidReqTooMuch,
	"INVALID_TRANSACTIONID": ErrInvalidTransactionID,
	"PARAM_ERROR":           ErrParamError,
	"FREQUENCY_LIMITED":     ErrFrequencyLimited,
	"REFUNDNOTEXIST":        ErrRefundNotExist,
	"NO_COMMENT":            ErrNoComment,
	"TIME_EXPIRE":           ErrTimeExpire,
}
