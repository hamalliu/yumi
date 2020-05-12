package tradeplatform

import (
	"encoding/json"
	"fmt"
	"time"
	"yumi/pkg/ecode"

	"yumi/external/pay/wxpay"
	"yumi/internal/pay/entities/trade"
	"yumi/utils"
)

const WxPay_JSAPI = trade.Way("wxpay_jsapi")

type WxJsapi struct {
	InternalWxPay
}

func GetWxJsapi() WxJsapi {
	return WxJsapi{}
}

type Wxh5PayRequest struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

func mashalWxh5PayRequest(appId, prePayId, privateKey string) (string, error) {
	var req Wxh5PayRequest

	req.AppId = appId
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
		NotifyUrl:      op.NotifyUrl,
		PayExpire:      op.PayExpire,
		SpbillCreateIp: op.SpbillCreateIp,
	}
	if retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeJsapi, wxMch, wxorder); err != nil {
		return ret, ecode.ServerErr(err)
	} else {
		ret.AppId = wxMch.AppId
		ret.MchId = wxMch.MchId
		if dataStr, err := mashalWxh5PayRequest(wxMch.AppId, retuo.PrepayId, wxMch.PrivateKey); err != nil {
			return ret, err
		} else {
			ret.Data = dataStr
		}
		return ret, nil
	}
}
