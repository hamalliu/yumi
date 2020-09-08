package tradeplatform

import (
	"encoding/json"

	"yumi/internal/trade"
	"yumi/pkg/ecode"
	"yumi/pkg/external/pay/wxpay"
)

//WxPayMWEB ...
const WxPayMWEB = trade.Way("wxpay_mweb")

var mweb WxMweb

//WxMweb ...
type WxMweb struct {
	InternalWxPay
	conf      Config
	sceneInfo string
}

//Init ...
func Init(conf Config) {
	mweb.conf = conf
	bytes, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	mweb.sceneInfo = string(bytes)
}

//GetWxMweb ...
func GetWxMweb() WxMweb {
	return mweb
}

//Config ...
type Config struct {
	H5Info H5Info `json:"h5_info"`
}

//H5Info ...
type H5Info struct {
	Type    string `json:"type"`     //场景类型
	WapURL  string `json:"wap_url"`  //WAP网站URL地址
	WapName string `json:"wap_name"` //WAP网站名
}

//Pay 发起支付
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
