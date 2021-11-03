package alipay

import (
	"fmt"
	"net/http"
	"net/url"

	"yumi/pkg/externalapi/aliapi/alipay"
	"yumi/usecase/trade/entity"
)

//AliPayPage ...

//Page ...
type Page struct {
	data Data
}

//NewPage ...
func NewPage(data Data) Page {
	return Page{data: data}
}

//Pay 发起支付
func (alipp Page) Pay(op entity.OrderPayAttribute) (entity.ReturnPay, error) {
	ret := entity.ReturnPay{}

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
		PayExpire:      op.PayExpire.Time(),
	}

	//获取收款商户信息
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	resp, err := alipay.GetDefault().UnifiedOrder(aliMch, pagePay)
	if err != nil {
		return ret, err
	}
	ret.AppID = aliMch.AppID
	ret.MchID = resp.SellerID
	ret.Data = string(resp.PagePayHTML)
	return ret, nil
}

//PayNotifyReq ...
func (alipp Page) PayNotifyReq(req *http.Request) (entity.ReturnPayNotify, error) {
	ret := entity.ReturnPayNotify{}

	rawQuery := req.URL.RawQuery
	reqNotify, err := alipay.ParseQuery(rawQuery)
	if err != nil {
		return ret, err
	}

	ret.OrderPayCode = reqNotify.PassbackParams
	ret.ReqData = reqNotify
	return ret, nil
}

//PayNotifyCheck ...
func (alipp Page) PayNotifyCheck(op entity.OrderPayAttribute, reqData interface{}) error {
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return err
	}

	reqNotify, ok := reqData.(alipay.ReqNotify)
	if ok {
		err := fmt.Errorf("转换类型失败，alipagepay")
		return err
	}
	if err := alipay.CheckPayNotify(aliMch, op.OutTradeNo, alipp.toPrice(op.TotalFee), op.MchID, reqNotify); err != nil {
		return err
	}
	return nil
}

//PayNotifyResp ...
func (alipp Page) PayNotifyResp(err error, resp http.ResponseWriter) {
	if err == nil {
		_, _ = resp.Write([]byte("success"))
	} else {
		_, _ = resp.Write([]byte(err.Error()))
	}
}

//QueryPayStatus ...
func (alipp Page) QueryPayStatus(op entity.OrderPayAttribute) (entity.ReturnQueryPay, error) {
	ret := entity.ReturnQueryPay{}

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
		return ret, err
	}

	if resp.OutTradeNo != op.OutTradeNo {
		err := fmt.Errorf("订单号不一致")
		return ret, err
	}
	if resp.TotalAmount != alipp.toPrice(op.TotalFee) {
		err := fmt.Errorf("订单金额不一致")
		return ret, err
	}
	ret.TransactionID = resp.TradeNo
	ret.BuyerLogonID = resp.BuyerlogonID

	switch ret.TradeStatus {
	case alipay.TradeStatusFinished, alipay.TradeStatusSuccess:
		ret.TradeStatus = entity.StatusTradePlatformSuccess
	case alipay.TradeStatusWaitBuyerPay:
		ret.TradeStatus = entity.StatusTradePlatformNotPay
	case alipay.TradeStatusCloseed:
		ret.TradeStatus = entity.StatusTradePlatformClosed
	default:
		err := fmt.Errorf("支付宝状态发生变动，请管理员及时更改")
		return ret, err
	}
	return ret, nil
}

//TradeClose ...
func (alipp Page) TradeClose(op entity.OrderPayAttribute) error {
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
		return err
	}

	if ret.TradeNo != op.TransactionID ||
		ret.OutTradeNo != op.OutTradeNo {
		err = fmt.Errorf("订单信息不一致")
		return err
	}

	return nil
}

//Refund ...
func (alipp Page) Refund(op entity.OrderPayAttribute, or entity.OrderRefundAttribute) error {
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
		return err
	}

	if ret.TradeNo != op.TransactionID ||
		ret.OutTradeNo != op.OutTradeNo ||
		ret.BuyerLogonID != op.BuyerLogonID ||
		ret.RefundFee != refundFee {
		err = fmt.Errorf("发起退款信息不一致，可能是订单数据被破坏")
		return err
	}

	return nil
}

//QueryRefundStatus ...
func (alipp Page) QueryRefundStatus(op entity.OrderPayAttribute, or entity.OrderRefundAttribute) (entity.ReturnQueryRefund, error) {
	//TODO
	return entity.ReturnQueryRefund{}, nil
}

// RefundNotifyReq ...
func (alipp Page) RefundNotifyReq(req *http.Request) (entity.ReturnRefundNotify, error) {
	// TODO:
	return entity.ReturnRefundNotify{}, nil
}

// RefundNotifyCheck 检查参数
func (alipp Page) RefundNotifyCheck(op entity.OrderPayAttribute, or entity.OrderRefundAttribute, reqData interface{}) error {
	// TODO:
	return nil
}

// RefundNotifyResp 应答
func (alipp Page) RefundNotifyResp(err error, resp http.ResponseWriter) {

}

func (alipp Page) getMch(sellerKey string) (alipay.Merchant, error) {
	ret := alipay.Merchant{}
	//获取收款商户信息
	mch, err := alipp.data.GetAliPayMerchant(sellerKey)
	if err != nil {
		return ret, err
	}
	ret.AppID = mch.AppID
	ret.PrivateKey = mch.PrivateKey
	ret.PublicKey = mch.PublicKey

	return ret, nil
}

func (alipp Page) toPrice(amount int) string {
	return fmt.Sprintf("%d.%02d", amount/100, amount%100)
}
