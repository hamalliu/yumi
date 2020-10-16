package tradeplatform

import (
	"encoding/json"
	"fmt"
	"time"

	"yumi/internal/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/external/trade/wxpay"
	"yumi/utils"
)

//WxPayAPP ...
const WxPayAPP = trade.Way("wxpay_app")

//WxApp ...
type WxApp struct {
	InternalWxPay
}

//GetWxApp ...
func GetWxApp() WxApp {
	return WxApp{}
}

//WxAppPayRequest ...
type WxAppPayRequest struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	PartnerID string `json:"partnerid"`
	PrepayID  string `json:"prepayid"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

//mashalWxAppPayRequest ...
func mashalWxAppPayRequest(appID, mchID, privateKey, prePayID string) (string, error) {
	var req WxAppPayRequest

	req.AppID = appID
	req.PartnerID = mchID
	req.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	req.NonceStr = utils.CreateRandomStr(30, utils.ALPHANUM)
	req.Package = fmt.Sprintf("prepay_id=%s", prePayID)
	req.SignType = "MD5"
	req.PaySign = wxpay.Buildsign(&req, wxpay.FieldTagKeyJSON, privateKey)

	reqBytes, err := json.Marshal(&req)
	if err != nil {
		return "", err
	}

	return string(reqBytes), nil
}

//Pay 发起支付
func (wxn1 WxApp) Pay(op trade.OrderPay) (trade.ReturnPay, error) {
	ret := trade.ReturnPay{}
	//获取收款商户信息
	wxMch, err := wxn1.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	wxorder := wxpay.UnifiedOrder{
		Body:           op.Body,
		Detail:         op.Detail,
		Attach:         op.Code,
		OutTradeNo:     op.OutTradeNo,
		TotalFee:       op.TotalFee,
		NotifyURL:      op.NotifyURL,
		PayExpire:      op.PayExpire,
		SpbillCreateIP: op.SpbillCreateIP,
	}

	retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeApp, wxMch, wxorder)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}

	ret.AppID = wxMch.AppID
	ret.MchID = wxMch.MchID
	dataStr, err := mashalWxAppPayRequest(wxMch.AppID, wxMch.MchID, wxMch.PrivateKey, retuo.PrepayID)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}

	ret.Data = dataStr
	return ret, nil
}