package thirdpf

import (
	"encoding/json"
	"fmt"
	"time"

	"yumi/pkg/externalapi/txapi/wxpay"
	"yumi/pkg/random"
	"yumi/usecase/trade/entity"
)

//NewWxPayJsapi ...
func NewWxPayJsapi() WxPayJsapi {
	return WxPayJsapi{}
}

//WxPayJsapi ...
type WxPayJsapi struct {
	InternalWxPay
}

//RequestWxPayh5 ...
type RequestWxPayh5 struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

//mashalWxh5PayRequest ...
func mashalRequestWxPayh5(appID, prePayID, privateKey string) (string, error) {
	var req RequestWxPayh5

	req.AppID = appID
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
func (wxn1 WxPayJsapi) Pay(op entity.OrderPayAttribute) (entity.ReturnPay, error) {
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

	retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeJsapi, wxMch, wxorder)
	if err != nil {
		return ret, err
	}
	ret.AppID = wxMch.AppID
	ret.MchID = wxMch.MchID
	dataStr, err := mashalRequestWxPayh5(wxMch.AppID, retuo.PrepayID, wxMch.PrivateKey)
	if err != nil {
		return ret, err
	}

	ret.Data = dataStr
	return ret, nil
}
