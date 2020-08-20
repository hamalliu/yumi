package tradeplatform

import (
	"encoding/json"
	"fmt"
	"time"
	"yumi/pkg/ecode"

	"yumi/internal/pay/trade"
	"yumi/pkg/external/pay/wxpay"
	"yumi/utils"
)

//WxPayJSAPI ...
const WxPayJSAPI = trade.Way("wxpay_jsapi")

//WxJsapi ...
type WxJsapi struct {
	InternalWxPay
}

//GetWxJsapi ...
func GetWxJsapi() WxJsapi {
	return WxJsapi{}
}

//Wxh5PayRequest ...
type Wxh5PayRequest struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

//mashalWxh5PayRequest ...
func mashalWxh5PayRequest(appID, prePayID, privateKey string) (string, error) {
	var req Wxh5PayRequest

	req.AppID = appID
	req.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	req.NonceStr = utils.CreateRandomStr(30, utils.ALPHANUM)
	req.Package = fmt.Sprintf("prepay_id=%s", prePayID)
	req.SignType = "MD5"
	req.PaySign = wxpay.Buildsign(&req, wxpay.FieldTagKeyJson, privateKey)

	reqBytes, err := json.Marshal(&req)
	if err != nil {
		return "", err
	}
	
	return string(reqBytes), nil
}

//Pay ...
func (wxn1 WxJsapi) Pay(op trade.OrderPay) (trade.ReturnPay, error) {
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
		NotifyUrl:      op.NotifyURL,
		PayExpire:      op.PayExpire,
		SpbillCreateIp: op.SpbillCreateIP,
	}

	retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeJsapi, wxMch, wxorder)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}
	ret.AppID = wxMch.AppId
	ret.MchID = wxMch.MchId
	dataStr, err := mashalWxh5PayRequest(wxMch.AppId, retuo.PrepayId, wxMch.PrivateKey)
	if err != nil {
		return ret, err
	}

	ret.Data = dataStr
	return ret, nil
}
