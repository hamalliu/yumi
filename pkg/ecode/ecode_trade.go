package ecode

/**
 *pay模块错误码（12000-13000）
 */

var (
	//OrderPayTimeout 订单已过期
	OrderPayTimeout = add(120001, "订单已过期")
	//NotSupportTradeWay 不支持的交易方式
	NotSupportTradeWay = add(120002, "不支持的交易方式")
	//InvalidSendPay 无效发起支付
	InvalidSendPay = add(120003, "无效发起支付")
	//InvalidQueryPay 无效支付查询
	InvalidQueryPay = add(120004, "无效支付查询")
	//InvalidCancelOrderPay 无效取消订单
	InvalidCancelOrderPay = add(120005, "无效取消订单")
	//InvalidCloseTrade 无效关闭交易
	InvalidCloseTrade = add(120006, "无效关闭交易")
	//CannotRepeatSendRefund 不能重复发起退款订单
	CannotRepeatSendRefund = add(120007, "不能重复发起退款订单")
	//InvalidQueryRefund 无效退款查询
	InvalidQueryRefund = add(120008, "无效退款查询")
	//InvalidSendRefund 无效退款发起
	InvalidSendRefund = add(120009, "无效退款发起")
)
