package alipay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

const (
	//ReturnCodeSuccess ...
	ReturnCodeSuccess = "10000"
)

const (
	//TradeStatusWaitBuyerPay 交易创建，等待买家付款
	TradeStatusWaitBuyerPay = "WAIT_BUYER_PAY"
	//TradeStatusCloseed 未付款交易超时关闭，或支付完成后全额退款
	TradeStatusCloseed = "TRADE_CLOSED"
	//TradeStatusSuccess 交易支付成功
	TradeStatusSuccess = "TRADE_SUCCESS"
	//TradeStatusFinished 交易结束，不可退款
	TradeStatusFinished = "TRADE_FINISHED"
)

//PayAPI ...
type PayAPI struct {
	GateWay string

	//统一收单下单并支付页面接口method
	UnifiedOrderMethod string

	//统一收单交易退款接口method
	RefundMethod string

	//统一收单交易退款查询method
	RefundQueryMethod string

	//统一收单线下交易查询method
	TradeQueryMethod string

	//统一收单交易关闭接口method
	TradeCloseMethod string

	//查询对账单下载地址method
	BillDownloadURLQueryMethod string

	//请求超时时间
	Timeout int

	////订单有效期（分钟）
	//OrderValidity int
}

var payapi = PayAPI{
	GateWay:                    "https://openapi.alipay.com/gateway.do",
	UnifiedOrderMethod:         "alipay.trade.page.pay",
	RefundMethod:               "alipay.trade.refund",
	RefundQueryMethod:          "alipay.trade.fastpay.refund.query",
	TradeQueryMethod:           "alipay.trade.query",
	TradeCloseMethod:           "alipay.trade.close",
	BillDownloadURLQueryMethod: "alipay.data.dataservice.bill.downloadurl.query",
	Timeout:                    15,
	//OrderValidity:              30,
}

//GetDefault ...
func GetDefault() PayAPI {
	return payapi
}

var req = ReqPublicPrameter{
	Format:    "JSON",
	ReturnURL: "", //TODO
	CharSet:   "utf-8",
	SignType:  "RSA2",
	Version:   "1.0",
}

//UnifiedOrder 统一收单下单并支付页面接口
func (p PayAPI) UnifiedOrder(mch Merchant, order PagePay) (PagePayReturn, error) {
	retn := PagePayReturn{}

	reqModel := req
	reqModel.AppID = mch.AppID
	reqModel.Method = p.UnifiedOrderMethod
	reqModel.Timestamp = time.Now().Format(timeFormat)
	reqModel.NotifyURL = order.NotifyURL
	reqModel.AppAuthToken = order.AppAuthToken

	reqOrder := ReqPagePay{
		OutTradeNo:        order.OutTradeNo,
		ProductCode:       order.ProductCode,
		TotalAmount:       order.TotalAmount,
		Subject:           order.Subject,
		Body:              order.Body,
		TimeExpire:        order.PayExpire.Format(timeFormat),
		GoodsType:         order.GoodsType,
		TimeoutExpress:    fmt.Sprintf("%dm", 30),
		EnablePayChannels: "balance,moneyFund,bankPay,debitCardExpress",
		QrPayMode:         "2",
		PassbackParams:    order.PassbackParams,
	}

	dataBytes, err := json.Marshal(&reqOrder)
	if err != nil {
		return retn, err
	}
	reqModel.BizContent = string(dataBytes)

	//加签
	sign, err := BuildSign(reqModel, mch.PrivateKey)
	if err != nil {
		return retn, err
	}
	req.Sign = sign

	//发起请求
	respModel := RespPagePay{}
	if _, err := request(&respModel, http.MethodPost, p.GateWay, &reqModel); err != nil {
		return retn, err
	}

	if respModel.Code == ReturnCodeSuccess {
		//验签
		if err := RespVerify(respModel, respModel.Sign, mch.PublicKey); err != nil {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.SubCode != "" {
			return retn, fmt.Errorf("%s", respModel.SubMsg)
		}

		//成功
		retn.TradeNo = respModel.TradeNo
		retn.OutTradeNo = respModel.OutTradeNo
		retn.SellerID = respModel.SellerID
		retn.TotalAmount = respModel.TotalAmount
		retn.MerchantOrderNo = respModel.MerchantOrderNo

		return retn, nil
	}
	//失败
	return retn, fmt.Errorf("%s", respModel.Msg)
}

//TradeQuery 统一收单线下交易查询
func (p PayAPI) TradeQuery(mch Merchant, query TradeQuery) (TradeQueryReturn, error) {
	retn := TradeQueryReturn{}

	if query.TradeNo == "" && query.OutTradeNo == "" {
		return retn, fmt.Errorf("不能同时为空")
	}

	reqModel := req
	reqModel.AppID = mch.AppID
	reqModel.Method = p.TradeQueryMethod
	reqModel.Timestamp = time.Now().Format(timeFormat)
	reqModel.AppAuthToken = query.AppAuthToken

	reqQuery := ReqTradeQuery{
		OutTradeNo: query.OutTradeNo,
		TradeNo:    query.TradeNo,
	}

	dataBytes, err := json.Marshal(&reqQuery)
	if err != nil {
		return retn, err
	}
	reqModel.BizContent = string(dataBytes)

	//加签
	sign, err := BuildSign(reqModel, mch.PrivateKey)
	if err != nil {
		return retn, err
	}
	req.Sign = sign

	//发起请求
	respModel := RespTradeQuery{}
	if _, err := request(&respModel, http.MethodPost, p.GateWay, &reqModel); err != nil {
		return retn, err
	}

	if respModel.Code == ReturnCodeSuccess {
		//验签
		if err := RespVerify(respModel, respModel.Sign, mch.PublicKey); err != nil {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.SubCode != "" {
			return retn, fmt.Errorf("%s", respModel.SubMsg)
		}

		//成功
		retn = respModel.TradeQueryReturn

		return retn, nil
	}
	//失败
	return retn, fmt.Errorf("%s", respModel.Msg)
}

//Refund 统一收单交易退款接口
func (p PayAPI) Refund(mch Merchant, refund Refund) (RefundReturn, error) {
	retn := RefundReturn{}

	if refund.TradeNo == "" && refund.OutTradeNo == "" {
		return retn, fmt.Errorf("商户订单号，支付宝交易号不能同时为空")
	}

	reqModel := req
	reqModel.AppID = mch.AppID
	reqModel.Method = p.RefundMethod
	reqModel.Timestamp = time.Now().Format(timeFormat)
	reqModel.AppAuthToken = refund.AppAuthToken

	reqRefund := ReqRefund{
		OutTradeNo:     refund.OutTradeNo,
		TradeNo:        refund.TradeNo,
		RefundAmount:   refund.RefundAmount,
		RefundCurrency: "CNY",
		RefundReason:   refund.RefundReason,
		OutRequestNo:   refund.OutRequestNo,
	}

	dataBytes, err := json.Marshal(&reqRefund)
	if err != nil {
		return retn, err
	}
	reqModel.BizContent = string(dataBytes)

	//加签
	sign, err := BuildSign(reqModel, mch.PrivateKey)
	if err != nil {
		return retn, err
	}
	req.Sign = sign

	//发起请求
	respModel := RespRefund{}
	if _, err := request(&respModel, http.MethodPost, p.GateWay, &reqModel); err != nil {
		return retn, err
	}

	if respModel.Code == ReturnCodeSuccess {
		//验签
		if err := RespVerify(respModel, respModel.Sign, mch.PublicKey); err != nil {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.SubCode != "" {
			return retn, fmt.Errorf("%s", respModel.SubMsg)
		}

		//成功
		retn = respModel.RefundReturn

		return retn, nil
	}
	//失败
	return retn, fmt.Errorf("%s", respModel.Msg)
}

//RefundQuery 统一收单交易退款查询
func (p PayAPI) RefundQuery(mch Merchant, refundQuery RefundQuery) (RefundQueryReturn, error) {
	retn := RefundQueryReturn{}

	if refundQuery.TradeNo == "" && refundQuery.OutTradeNo == "" {
		return retn, fmt.Errorf("商户订单号，支付宝交易号不能同时为空")
	}

	reqModel := req
	reqModel.AppID = mch.AppID
	reqModel.Method = p.RefundQueryMethod
	reqModel.Timestamp = time.Now().Format(timeFormat)
	reqModel.AppAuthToken = refundQuery.AppAuthToken

	reqQuery := ReqRefundQuery{
		TradeNo:      refundQuery.TradeNo,
		OutTradeNo:   refundQuery.OutTradeNo,
		OutRequestNo: refundQuery.OutRequestNo,
	}

	dataBytes, err := json.Marshal(&reqQuery)
	if err != nil {
		return retn, err
	}
	reqModel.BizContent = string(dataBytes)

	//加签
	sign, err := BuildSign(reqModel, mch.PrivateKey)
	if err != nil {
		return retn, err
	}
	req.Sign = sign

	//发起请求
	respModel := RespRefundQuery{}
	if _, err := request(&respModel, http.MethodPost, p.GateWay, &reqModel); err != nil {
		return retn, err
	}

	if respModel.Code == ReturnCodeSuccess {
		//验签
		if err := RespVerify(respModel, respModel.Sign, mch.PublicKey); err != nil {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.SubCode != "" {
			return retn, fmt.Errorf("%s", respModel.SubMsg)
		}

		//成功
		retn = respModel.RefundQueryReturn

		return retn, nil
	}
	//失败
	return retn, fmt.Errorf("%s", respModel.Msg)
}

//TradeClose 统一收单交易关闭接口
func (p PayAPI) TradeClose(mch Merchant, close TradeClose) (TradeCloseReturn, error) {
	retn := TradeCloseReturn{}

	if close.TradeNo == "" && close.OutTradeNo == "" {
		return retn, fmt.Errorf("商户订单号，支付宝交易号不能同时为空")
	}

	reqModel := req
	reqModel.AppID = mch.AppID
	reqModel.Method = p.TradeCloseMethod
	reqModel.Timestamp = time.Now().Format(timeFormat)
	reqModel.AppAuthToken = close.AppAuthToken

	reqClose := ReqTradeClose{
		TradeNo:    close.TradeNo,
		OutTradeNo: close.OutTradeNo,
		OperatorID: close.OperatorID,
	}

	dataBytes, err := json.Marshal(&reqClose)
	if err != nil {
		return retn, err
	}
	reqModel.BizContent = string(dataBytes)

	//加签
	sign, err := BuildSign(reqModel, mch.PrivateKey)
	if err != nil {
		return retn, err
	}
	req.Sign = sign

	//发起请求
	respModel := RespTradeClose{}
	if _, err := request(&respModel, http.MethodPost, p.GateWay, &reqModel); err != nil {
		return retn, err
	}

	if respModel.Code == ReturnCodeSuccess {
		//验签
		if err := RespVerify(respModel, respModel.Sign, mch.PublicKey); err != nil {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.SubCode != "" {
			return retn, fmt.Errorf("%s", respModel.SubMsg)
		}

		//成功
		retn = respModel.TradeCloseReturn

		return retn, nil
	}
	//失败
	return retn, fmt.Errorf("%s", respModel.Msg)
}

//BillDownloadURLQuery 查询对账单下载地址
func (p PayAPI) BillDownloadURLQuery(mch Merchant, bill BillDownloadURLQuery) (BillDownloadURLQueryReturn, error) {
	retn := BillDownloadURLQueryReturn{}

	reqModel := req
	reqModel.AppID = mch.AppID
	reqModel.Method = p.BillDownloadURLQueryMethod
	reqModel.Timestamp = time.Now().Format(timeFormat)
	reqModel.AppAuthToken = bill.AppAuthToken

	reqBill := ReqBillDownloadURLQuery{
		BillDate: bill.BillDate,
		BillType: bill.BillType,
	}

	dataBytes, err := json.Marshal(&reqBill)
	if err != nil {
		return retn, err
	}
	reqModel.BizContent = string(dataBytes)

	//加签
	sign, err := BuildSign(reqModel, mch.PrivateKey)
	if err != nil {
		return retn, err
	}
	req.Sign = sign

	//发起请求
	respModel := RespBillDownloadURLQuery{}
	if _, err := request(&respModel, http.MethodPost, p.GateWay, &reqModel); err != nil {
		return retn, err
	}

	if respModel.Code == ReturnCodeSuccess {
		//验签
		if err := RespVerify(respModel, respModel.Sign, mch.PublicKey); err != nil {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.SubCode != "" {
			return retn, fmt.Errorf("%s", respModel.SubMsg)
		}

		//成功
		retn = respModel.BillDownloadURLQueryReturn

		return retn, nil
	}
	//失败
	return retn, fmt.Errorf("%s", respModel.Msg)
}
