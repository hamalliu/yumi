package thirdpf

import (
	"encoding/json"
	"fmt"
	"time"

	"yumi/pkg/ecode"
	"yumi/pkg/externalapi/txapi/wxpay"
	"yumi/pkg/random"
	"yumi/usecase/trade/entity"
)

//NewWxPayApp ...
func NewWxPayApp() WxPayApp {
	return WxPayApp{}
}

//WxPayApp ...
type WxPayApp struct {
	InternalWxPay
}

//RequestWxPayApp ...
type RequestWxPayApp struct {
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
func mashalRequestWxPayApp(appID, mchID, privateKey, prePayID string) (string, error) {
	var req RequestWxPayApp

	req.AppID = appID
	req.PartnerID = mchID
	req.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	req.NonceStr = random.Get(30, random.ALPHANUM)
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
func (wxn1 WxPayApp) Pay(op entity.OrderPayAttribute) (entity.ReturnPay, error) {
	ret := entity.ReturnPay{}
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
		PayExpire:      op.PayExpire.Time(),
		SpbillCreateIP: op.SpbillCreateIP,
	}

	retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeApp, wxMch, wxorder)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}

	ret.AppID = wxMch.AppID
	ret.MchID = wxMch.MchID
	dataStr, err := mashalRequestWxPayApp(wxMch.AppID, wxMch.MchID, wxMch.PrivateKey, retuo.PrepayID)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}

	ret.Data = dataStr
	return ret, nil
}
