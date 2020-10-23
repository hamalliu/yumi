package tradeplatform

import (
	"fmt"
	"net/http"
	"net/url"

	"yumi/usecase/trade"
	"yumi/usecase/trade/db"
	"yumi/pkg/ecode"
	"yumi/pkg/external/trade/alipay"
)

//AliPayPage ...
const AliPayPage = trade.Way("alipay_page")

//AliPagePay ...
type AliPagePay string

//GetAliPagePay ...
func GetAliPagePay() AliPagePay {
	return ""
}

//Pay 发起支付
func (alipp AliPagePay) Pay(op trade.OrderPay) (trade.ReturnPay, error) {
	ret := trade.ReturnPay{}

	//下单
	pagePay := alipay.PagePay{
		OutTradeNo:     op.OutTradeNo,
		ProductCode:    op.Code,
		TotalAmount:    alipp.toPrice(op.TotalFee),
		Subject:        op.Body,
		Body:           op.Detail,
		GoodsType:      "0",
		NotifyURL:      op.NotifyURL,
		PassbackParams: url.QueryEscape(op.Code), //必须urlencode
		PayExpire:      op.PayExpire,
	}

	//获取收款商户信息
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	resp, err := alipay.GetDefault().UnifiedOrder(aliMch, pagePay)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}
	ret.AppID = aliMch.AppID
	ret.MchID = resp.SellerID
	ret.Data = string(resp.PagePayHTML)
	return ret, nil
}

//PayNotifyReq ...
func (alipp AliPagePay) PayNotifyReq(req *http.Request) (trade.ReturnPayNotify, error) {
	ret := trade.ReturnPayNotify{}

	rawQuery := req.URL.RawQuery
	reqNotify, err := alipay.ParseQuery(rawQuery)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}

	ret.OrderPayCode = reqNotify.PassbackParams
	ret.ReqData = reqNotify
	return ret, nil
}

//PayNotifyCheck ...
func (alipp AliPagePay) PayNotifyCheck(op trade.OrderPay, reqData interface{}) error {
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return err
	}

	reqNotify, ok := reqData.(alipay.ReqNotify)
	if ok {
		err := fmt.Errorf("转换类型失败，alipagepay")
		return ecode.ServerErr(err)
	}
	if err := alipay.CheckPayNotify(aliMch, op.OutTradeNo, alipp.toPrice(op.TotalFee), op.MchID, reqNotify); err != nil {
		return err
	}
	return nil
}

//PayNotifyResp ...
func (alipp AliPagePay) PayNotifyResp(err error, resp http.ResponseWriter) {
	if err == nil {
		_, _ = resp.Write([]byte("success"))
	} else {
		_, _ = resp.Write([]byte(err.Error()))
	}

	return
}

//QueryPayStatus ...
func (alipp AliPagePay) QueryPayStatus(op trade.OrderPay) (trade.ReturnQueryPay, error) {
	ret := trade.ReturnQueryPay{}

	tradeQuery := alipay.TradeQuery{
		TradeNo:    op.TransactionID,
		OutTradeNo: op.OutTradeNo,
	}

	//获取收款商户信息
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	resp, err := alipay.GetDefault().TradeQuery(aliMch, tradeQuery)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}

	if resp.OutTradeNo != op.OutTradeNo {
		err := fmt.Errorf("订单号不一致")
		return ret, ecode.ServerErr(err)
	}
	if resp.TotalAmount != alipp.toPrice(op.TotalFee) {
		err := fmt.Errorf("订单金额不一致")
		return ret, ecode.ServerErr(err)
	}
	ret.TransactionID = resp.TradeNo
	ret.BuyerLogonID = resp.BuyerlogonID

	switch ret.TradeStatus {
	case alipay.TradeStatusSuccess:
		ret.TradeStatus = trade.StatusTradePlatformSuccess
	case alipay.TradeStatusWaitBuyerPay:
		ret.TradeStatus = trade.StatusTradePlatformNotPay
	case alipay.TradeStatusCloseed:
		ret.TradeStatus = trade.StatusTradePlatformClosed
	case alipay.TradeStatusFinished:
		ret.TradeStatus = trade.StatusTradePlatformFinished
	default:
		err := fmt.Errorf("支付宝状态发生变动，请管理员及时更改")
		return ret, ecode.ServerErr(err)
	}
	return ret, nil
}

//TradeClose ...
func (alipp AliPagePay) TradeClose(op trade.OrderPay) error {
	tradeClose := alipay.TradeClose{
		OutTradeNo: op.OutTradeNo,
		TradeNo:    op.TransactionID,
		OperatorID: "sys",
	}

	//获取收款商户信息
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return err
	}

	ret, err := alipay.GetDefault().TradeClose(aliMch, tradeClose)
	if err != nil {
		return ecode.ServerErr(err)
	}

	if ret.TradeNo != op.TransactionID ||
		ret.OutTradeNo != op.OutTradeNo {
		err = fmt.Errorf("订单信息不一致")
		return ecode.ServerErr(err)
	}

	return nil
}

//Refund ...
func (alipp AliPagePay) Refund(op trade.OrderPay, or trade.OrderRefund) error {
	//获取收款商户信息
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return err
	}
	refundFee := alipp.toPrice(or.RefundFee)
	rfd := alipay.Refund{
		OutTradeNo:   op.OutTradeNo,
		TradeNo:      op.TransactionID,
		RefundAmount: refundFee,
		RefundReason: or.RefundDesc,
		OutRequestNo: or.Code,
	}

	ret, err := alipay.GetDefault().Refund(aliMch, rfd)
	if err != nil {
		return ecode.ServerErr(err)
	}

	if ret.TradeNo != op.TransactionID ||
		ret.OutTradeNo != op.OutTradeNo ||
		ret.BuyerLogonID != op.BuyerLogonID ||
		ret.RefundFee != refundFee {
		err = fmt.Errorf("发起退款信息不一致，可能是订单数据被破坏")
		return ecode.ServerErr(err)
	}

	return nil
}

//QueryRefundStatus ...
func (alipp AliPagePay) QueryRefundStatus(op trade.OrderPay, or trade.OrderRefund) {
	//TODO
}

func (alipp AliPagePay) getMch(sellerKey string) (alipay.Merchant, error) {
	ret := alipay.Merchant{}
	//获取收款商户信息
	mch, err := db.GetAliPayMerchantBySellerKey(sellerKey)
	if err != nil {
		return ret, err
	}
	ret.AppID = mch.AppID
	ret.PrivateKey = mch.PrivateKey
	ret.PublicKey = mch.PublicKey

	return ret, nil
}

func (alipp AliPagePay) toPrice(amount int) string {
	return fmt.Sprintf("%d.%02d", amount/100, amount%100)
}
