package tradeplatform

import (
	"encoding/json"
	"yumi/pkg/ecode"

	"yumi/external/pay/wxpay"
	"yumi/internal/pay/entities/trade"
)

const WxPay_MWEB = trade.Way("wxpay_mweb")

var mweb WxMweb

type WxMweb struct {
	InternalWxPay
	conf      Config
	sceneInfo string
}

func Init(conf Config) {
	mweb.conf = conf
	bytes, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	mweb.sceneInfo = string(bytes)
}

func GetWxMweb() WxMweb {
	return mweb
}

type Config struct {
	H5Info H5Info `json:"h5_info"`
}

type H5Info struct {
	Type    string `json:"type"`     //场景类型
	WapUrl  string `json:"wap_url"`  //WAP网站URL地址
	WapName string `json:"wap_name"` //WAP网站名
}

func (wxn1 WxMweb) Pay(op trade.OrderPay) (trade.ReturnPay, error) {
	ret := trade.ReturnPay{}
	//获取收款商户信息
	wxMch, err := wxn1.getMch(op.SellerKey)
	if err != nil {
		return ret, err
	}

	//

	wxorder := wxpay.UnifiedOrder{
		Body:           op.Body,
		Detail:         op.Detail,
		Attach:         op.Code,
		OutTradeNo:     op.OutTradeNo,
		TotalFee:       op.TotalFee,
		NotifyUrl:      op.NotifyUrl,
		PayExpire:      op.PayExpire,
		SpbillCreateIp: op.SpbillCreateIp,
		SceneInfo:      wxn1.sceneInfo,
	}
	if retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeNative, wxMch, wxorder); err != nil {
		return ret, ecode.ServerErr(err)
	} else {
		ret.AppId = wxMch.AppId
		ret.MchId = wxMch.MchId
		ret.Data = retuo.MwebUrl
		return ret, nil
	}
}
