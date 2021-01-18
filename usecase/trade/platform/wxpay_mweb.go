package platform

import (
	"encoding/json"

	"yumi/usecase/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/externalapi/txapi/wxpay"
)

//NewWxPayMweb ...
func NewWxPayMweb(conf WxPayMwebConfig) WxPayMweb {
	mweb := WxPayMweb {}
	mweb.conf = conf
	bytes, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	mweb.sceneInfo = string(bytes)
	return mweb
}

//WxPayMweb ...
type WxPayMweb struct {
	InternalWxPay
	conf      WxPayMwebConfig
	sceneInfo string
}

//WxPayMwebConfig ...
type WxPayMwebConfig struct {
	H5Info H5Info `json:"h5_info"`
}

//H5Info ...
type H5Info struct {
	Type    string `json:"type"`     //场景类型
	WapURL  string `json:"wap_url"`  //WAP网站URL地址
	WapName string `json:"wap_name"` //WAP网站名
}

//Pay 发起支付
func (wxn1 WxPayMweb) Pay(op trade.OrderPay) (trade.ReturnPay, error) {
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
		NotifyURL:      op.NotifyURL,
		PayExpire:      op.PayExpire,
		SpbillCreateIP: op.SpbillCreateIP,
		SceneInfo:      wxn1.sceneInfo,
	}

	retuo, err := wxpay.GetDefault().UnifiedOrder(wxpay.TradeTypeNative, wxMch, wxorder)
	if err != nil {
		return ret, ecode.ServerErr(err)
	}
	ret.AppID = wxMch.AppID
	ret.MchID = wxMch.MchID
	ret.Data = retuo.MwebURL
	return ret, nil
}