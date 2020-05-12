package tradeplatform

import (
	"fmt"
	"net/http"
	"net/url"

	"yumi/external/pay/alipay"
	"yumi/internal/pay/entities/trade"
	"yumi/internal/pay/humble/dbentities"
	"yumi/pkg/ecode"
)

const AliPay_Page = trade.Way("alipay_page")

type AliPagePay string

func GetAliPagePay() AliPagePay {
	return ""
}

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
		NotifyUrl:      op.NotifyUrl,
		PassbackParams: url.QueryEscape(op.Code), //必须urlencode
		PayExpire:      op.PayExpire,
	}

	//获取收款商户信息
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	if resp, err := alipay.GetDefault().UnifiedOrder(aliMch, pagePay); err != nil {
		return ret, ecode.ServerErr(err)
	} else {
		ret.AppId = aliMch.AppId
		ret.MchId = resp.SellerId
		ret.Data = string(resp.PagePayHtml)
		return ret, nil
	}
}

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

func (alipp AliPagePay) PayNotifyCheck(op trade.OrderPay, reqData interface{}) error {
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return err
	}

	if reqNotify, ok := reqData.(alipay.ReqNotify); ok {
		err := fmt.Errorf("转换类型失败，wxnative")
		return ecode.ServerErr(err)
	} else {
		if err := alipay.CheckPayNotify(aliMch, op.OutTradeNo, alipp.toPrice(op.TotalFee), op.MchId, reqNotify); err != nil {
			return err
		}
	}

	return nil
}

func (alipp AliPagePay) PayNotifyResp(err error, resp http.ResponseWriter) {
	if err == nil {
		_, _ = resp.Write([]byte("success"))
	} else {
		_, _ = resp.Write([]byte(err.Error()))
	}

	return
}

func (alipp AliPagePay) QueryPayStatus(op trade.OrderPay) (trade.ReturnQueryPay, error) {
	ret := trade.ReturnQueryPay{}

	tradeQuery := alipay.TradeQuery{
		TradeNo:    op.TransactionId,
		OutTradeNo: op.OutTradeNo,
	}

	//获取收款商户信息
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	if resp, err := alipay.GetDefault().TradeQuery(aliMch, tradeQuery); err != nil {
		return ret, ecode.ServerErr(err)
	} else {
		if resp.OutTradeNo != op.OutTradeNo {
			err := fmt.Errorf("订单号不一致")
			return ret, ecode.ServerErr(err)
		}
		if resp.TotalAmount != alipp.toPrice(op.TotalFee) {
			err := fmt.Errorf("订单金额不一致")
			return ret, ecode.ServerErr(err)
		}
		ret.TransactionId = resp.TradeNo
		ret.BuyerLogonId = resp.BuyerlogonId

		switch ret.TradeStatus {
		case alipay.TradeStatusSuccess:
			ret.TradeStatus = trade.Success
		case alipay.TradeStatusWaitBuyerPay:
			ret.TradeStatus = trade.NotPay
		case alipay.TradeStatusCloseed:
			ret.TradeStatus = trade.Closed
		case alipay.TradeStatusFinished:
			ret.TradeStatus = trade.Finished
		default:
			err := fmt.Errorf("支付宝状态发生变动，请管理员及时更改")
			return ret, ecode.ServerErr(err)
		}
		return ret, nil
	}
}

func (alipp AliPagePay) TradeClose(op trade.OrderPay) error {
	tradeClose := alipay.TradeClose{
		OutTradeNo: op.OutTradeNo,
		TradeNo:    op.TransactionId,
		OperatorId: "sys",
	}

	//获取收款商户信息
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return err
	}

	if ret, err := alipay.GetDefault().TradeClose(aliMch, tradeClose); err != nil {
		return ecode.ServerErr(err)
	} else {
		if ret.TradeNo != op.TransactionId ||
			ret.OutTradeNo != op.OutTradeNo {
			err = fmt.Errorf("订单信息不一致")
			return ecode.ServerErr(err)
		}
	}

	return nil
}

func (alipp AliPagePay) Refund(op trade.OrderPay, or trade.OrderRefund) error {
	//获取收款商户信息
	aliMch, err := alipp.getMch(op.SellerKey)
	if err != nil {
		return err
	}
	refundFee := alipp.toPrice(or.RefundFee)
	rfd := alipay.Refund{
		OutTradeNo:   op.OutTradeNo,
		TradeNo:      op.TransactionId,
		RefundAmount: refundFee,
		RefundReason: or.RefundDesc,
		OutRequestNo: or.Code,
	}

	if ret, err := alipay.GetDefault().Refund(aliMch, rfd); err != nil {
		return ecode.ServerErr(err)
	} else {
		if ret.TradeNo != op.TransactionId ||
			ret.OutTradeNo != op.OutTradeNo ||
			ret.BuyerLogonId != op.BuyerLogonId ||
			ret.RefundFee != refundFee {
			err = fmt.Errorf("发起退款信息不一致，可能是订单数据被破坏")
			return ecode.ServerErr(err)
		}
	}

	return nil
}

func (alipp AliPagePay) QueryRefundStatus(op trade.OrderPay, or trade.OrderRefund) {
	//TODO
}

func (alipp AliPagePay) getMch(sellerKey string) (alipay.Merchant, error) {
	ret := alipay.Merchant{}
	//获取收款商户信息
	mch, err := dbentities.GetAliPayMerchantBySellerKey(sellerKey)
	if err != nil {
		return ret, err
	}
	ret.AppId = mch.AppId
	ret.PrivateKey = mch.PrivateKey
	ret.PublicKey = mch.PublicKey

	return ret, nil
}

func (alipp AliPagePay) toPrice(amount int) string {
	return fmt.Sprintf("%d.%02d", amount/100, amount%100)
}
