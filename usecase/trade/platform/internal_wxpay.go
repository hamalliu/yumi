package platform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"yumi/usecase/trade/db"
	"yumi/usecase/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/externalapi/txapi/wxpay"
)

//InternalWxPay ...
type InternalWxPay struct{}

//PayNotifyReq ...
func (iwp InternalWxPay) PayNotifyReq(req *http.Request) (trade.ReturnPayNotify, error) {
	ret := trade.ReturnPayNotify{}

	//获取通知参数
	reqJ := wxpay.ReqPayNotify{}
	if err := json.NewDecoder(req.Body).Decode(&reqJ); err != nil {
		return ret, ecode.ServerErr(err)
	}

	ret.OrderPayCode = reqJ.Attach
	ret.ReqData = reqJ
	return ret, nil
}

//PayNotifyCheck ...
func (iwp InternalWxPay) PayNotifyCheck(op trade.OrderPay, reqData interface{}) error {
	//获取收款商户信息
	wxMch, err := iwp.getMch(op.SellerKey)
	if err != nil {
		return err
	}

	reqJ, ok := reqData.(wxpay.ReqPayNotify)
	if ok {
		err := fmt.Errorf("转换类型失败")
		return ecode.ServerErr(err)
	} 

	if err := wxpay.CheckPayNotify(wxMch, op.TotalFee, op.OutTradeNo, reqJ); err != nil {
		return err
	}
	return nil
}

//PayNotifyResp ...
func (iwp InternalWxPay) PayNotifyResp(err error, resp http.ResponseWriter) {
	respJ := wxpay.RespPayNotify{}
	if err == nil {
		respJ.ReturnCode = "SUCCESS"
		respJ.ReturnMsg = "OK"
	} else {
		respJ.ReturnCode = "FAIL"
		respJ.ReturnMsg = err.Error()
	}

	bytes, _ := json.Marshal(&respJ)

	_, _ = resp.Write(bytes)

	return
}

//QueryPayStatus ...
func (iwp InternalWxPay) QueryPayStatus(op trade.OrderPay) (trade.ReturnQueryPay, error) {
	ret := trade.ReturnQueryPay{}
	//获取收款商户信息
	wxMch, err := iwp.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	resp, err := wxpay.GetDefault().OrderQuery(wxMch, op.TransactionID, op.OutTradeNo)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}
	
	if op.OutTradeNo != resp.OutTradeNo {
		if resp.OutTradeNo != op.OutTradeNo {
			err := fmt.Errorf("订单号不一致")
			return ret, ecode.ServerErr(err)
		}
		if resp.TotalFee != op.TotalFee {
			err := fmt.Errorf("订单金额不一致")
			return ret, ecode.ServerErr(err)
		}
	}
	ret.TransactionID = resp.TransactionID
	switch resp.TradeState {
	case wxpay.TradeStateSuccess:
		ret.TradeStatus = trade.StatusTradePlatformSuccess
	case wxpay.TradeStateNotpay, wxpay.TradeStateUserPaying, wxpay.TradeStatePayError,
		wxpay.TradeStateRefund, wxpay.TradeStateRevoked:
		ret.TradeStatus = trade.StatusTradePlatformNotPay
	case wxpay.TradeStateClosed:
		ret.TradeStatus = trade.StatusTradePlatformClosed
	default:
		err := fmt.Errorf("微信支付状态发生变动，请管理员及时更改")
		return ret, ecode.ServerErr(err)
	}

	return ret, nil
}

//TradeClose ...
func (iwp InternalWxPay) TradeClose(op trade.OrderPay) error {
	//获取收款商户信息
	wxMch, err := iwp.getMch(op.SellerKey)
	if err != nil {
		return err
	}

	if err := wxpay.GetDefault().CloseOrder(wxMch, op.OutTradeNo); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

//Refund ...
func (iwp InternalWxPay) Refund(op trade.OrderPay, or trade.OrderRefund) error {
	//获取收款商户信息
	wxMch, err := iwp.getMch(op.SellerKey)
	if err != nil {
		return err
	}

	rfd := wxpay.Refund{
		TransactionID: op.TransactionID,
		OutTradeNo:    op.OutTradeNo,
		OutRefundNo:   or.Code, //商户退款单号就是退款订单的唯一编码
		TotalFee:      op.TotalFee,
		RefundFee:     or.RefundFee,
		RefundDesc:    or.RefundDesc,
		NotifyURL:     or.NotifyURL,
		CertP12:       getP12(),
	}

	resp, err := wxpay.GetDefault().Refund(wxMch, rfd)
	if err != nil {
		return ecode.ServerErr(err)
	}
	
	if resp.TransactionID != op.TransactionID ||
		resp.OutTradeNo != op.OutTradeNo ||
		resp.OutRefundNo != or.OutRefundNo {
		err = fmt.Errorf("订单错误")
		return ecode.ServerErr(err)
	}
	if resp.RefundFee != or.RefundFee {
		err = fmt.Errorf("退款金额不一致")
		return ecode.ServerErr(err)
	}
	if resp.TotalFee != op.TotalFee {
		err = fmt.Errorf("订单总金额不一致")
		return ecode.ServerErr(err)
	}

	return nil
}

//QueryRefundStatus ...
func (iwp InternalWxPay) QueryRefundStatus(op trade.OrderPay, or trade.OrderRefund) (trade.ReturnQueryRefund, error) {
	ret := trade.ReturnQueryRefund{}
	//获取收款商户信息
	wxMch, err := iwp.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	rq := wxpay.RefundQuery{
		TransactionID: op.TransactionID,
		OutTradeNo:    op.OutTradeNo,
		OutRefundNo:   or.OutRefundNo,
		Offset:        or.SerialNum,
	}

	resp, err := wxpay.GetDefault().RefundQuery(wxMch, rq)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}

	ret.RefundID = resp.RefundIdn
	switch resp.RefundStatusn {
	case wxpay.TradeStateSuccess:
		ret.TradeStatus = trade.StatusTradePlatformSuccess
	case wxpay.TradeStateRefundClose:
		ret.TradeStatus = trade.StatusTradePlatformClosed
	case wxpay.TradeStateProcessing:
		ret.TradeStatus = trade.StatusTradePlatformRefundProcessing
	case wxpay.TradeStateChange:
		ret.TradeStatus = trade.StatusTradePlatformError
	default:
		err := fmt.Errorf("微信退款状态未识别 %s", resp.RefundStatusn)
		return ret, ecode.ServerErr(err)
	}
	return ret, nil
}

//RefundNotifyReq ...
func (iwp InternalWxPay) RefundNotifyReq(req *http.Request) (trade.ReturnRefundNotify, error) {
	ret := trade.ReturnRefundNotify{}

	n, err := wxpay.GetRefundNotify(req)
	if err != nil {
		return ret, err
	}

	//获取收款商户信息
	wxMch, err := iwp.getMch2(n.MchID)
	if err != nil {
		return ret, err
	}

	n.DecryptReqInfo, err = wxpay.DecryptoRefundNotify(wxMch, n.ReqInfo)
	if err != nil {
		return ret, err
	}

	ret.ReqData = n
	ret.OrderRefundCode = n.DecryptReqInfo.OutRefundNo
	return ret, nil
}

//RefundNotifyCheck ...
func (iwp InternalWxPay) RefundNotifyCheck(op trade.OrderPay, or trade.OrderRefund, reqData interface{}) error {
	//获取收款商户信息
	wxMch, err := iwp.getMch(op.SellerKey)
	if err != nil {
		return err
	}

	if err := wxpay.CheckRefundNotify(wxMch, reqData.(wxpay.ReqRefundNotify)); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

//RefundNotifyResp ...
func (iwp InternalWxPay) RefundNotifyResp(err error, resp http.ResponseWriter) {
	respJ := wxpay.RespRefundNotify{}
	if err == nil {
		respJ.ReturnCode = "SUCCESS"
		respJ.ReturnMsg = "OK"
	} else {
		respJ.ReturnCode = "FAIL"
		respJ.ReturnMsg = err.Error()
	}

	bytes, _ := json.Marshal(&respJ)

	_, _ = resp.Write(bytes)

	return
}

func (iwp InternalWxPay) getMch(sellerKey string) (wxpay.Merchant, error) {
	ret := wxpay.Merchant{}

	//获取收款商户信息
	mch, err := db.GetWxPayMerchantBySellerKey(sellerKey)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}
	ret.AppID = mch.AppID
	ret.MchID = mch.MchID
	ret.PrivateKey = mch.PrivateKey

	return ret, nil
}

func (iwp InternalWxPay) getMch2(mchID string) (wxpay.Merchant, error) {
	ret := wxpay.Merchant{}

	//获取收款商户信息
	mch, err := db.GetWxPayMerchantByMchID(mchID)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}
	ret.AppID = mch.AppID
	ret.MchID = mch.MchID
	ret.PrivateKey = mch.PrivateKey

	return ret, nil
}

func getP12() []byte {
	f, err := os.Open("E:\\commonsoft\\wxcertuils\\WXCertUtil\\cert\\apiclient_cert.p12")
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return bs
}
