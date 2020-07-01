package tradeplatform

import (
	"encoding/json"
	"fmt"
	"time"

	"yumi/internal/pay/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/external/pay/wxpay"
	"yumi/utils"
)

const WxPay_APP = trade.Way("wxpay_app")

type WxApp struct {
	InternalWxPay
}

func GetWxApp() WxApp {
	return WxApp{}
}

type WxAppPayRequest struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	PartnerId string `json:"partnerid"`
	PrepayId  string `json:"prepayid"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

func mashalWxAppPayRequest(appId, mchId, privateKey, prePayId string) (string, error) {
	var req WxAppPayRequest

	req.AppId = appId
	req.PartnerId = mchId
	req.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	req.NonceStr = utils.CreateRandomStr(30, utils.ALPHANUM)
	req.Package = fmt.Sprintf("prepay_id=%s", prePayId)
	req.SignType = "MD5"
	req.PaySign = wxpay.Buildsign(&req, wxpay.FieldTagKeyJson, privateKey)

	if reqBytes, err := json.Marshal(&req); err != nil {
		return "", err
	} else {
		return string(reqBytes), nil
	}
}

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
		NotifyUrl:      op.NotifyUrl,
		PayExpire:      op.PayExpire,
		SpbillCreateIp: op.SpbillCreateIp,
	}
	if retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeApp, wxMch, wxorder); err != nil {
		return ret, ecode.ServerErr(err)
	} else {
		ret.AppId = wxMch.AppId
		ret.MchId = wxMch.MchId
		if dataStr, err := mashalWxAppPayRequest(wxMch.AppId, wxMch.MchId, wxMch.PrivateKey, retuo.PrepayId); err != nil {
			return ret, ecode.ServerErr(err)
		} else {
			ret.Data = dataStr
		}
		return ret, nil
	}
}
