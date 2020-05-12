package wxpay

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"net/http"
	url2 "net/url"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"

	"yumi/external/pay/internal"
)

/**
 *微信支付
 */

const (
	FieldTagKeyXml  = "xml"
	FieldTagKeyJson = "json"
)

const (
	ReturnCodeSuccess = "SUCCESS"
	ReturnCodeFail    = "FAIL"
)

const (
	TradeStateSuccess     = "SUCCESS"
	TradeStateRefund      = "REFUND"
	TradeStateNotpay      = "NOTPAY"
	TradeStateClosed      = "CLOSED"
	TradeStateRevoked     = "REVOKED"
	TradeStateUserPaying  = "USERPAYING"
	TradeStatePayError    = "PAYERROR"
	TradeStateRefundClose = "REFUNDCLOSE"
	TradeStateProcessing  = "PROCESSING"
	TradeStateChange      = "CHANGE"
)

const (
	timeFormat = "20060102150405"
	dateFormat = "20060102"
)

type TradeType string

const (
	TradeTypeJsapi  TradeType = "JSAPI"
	TradeTypeNative TradeType = "NATIVE"
	TradeTypeMweb   TradeType = "MWEB"
	TradeTypeApp    TradeType = "APP"
)

type PayApi struct {
	//生成支付二维码url
	BizPayUrlUrl string

	//统一下单url
	UnifiedOrderUrl string

	//统一下单url（备用域名）
	UnifiedOrderUrl2 string

	//查询订单url
	OrderQueryUrl string

	//查询订单url（备用域名）
	OrderQueryUrl2 string

	//关闭订单url
	CloseOrderUrl string

	//关闭订单url（备用域名）
	CloseOrderUrl2 string

	//申请退款url
	RefundUrl string

	//查询退款url
	RefundQueryUrl string

	//下载对账单url
	DownloadBillUrl string

	//下载资金账单url
	DownloadFundFlowUrl string

	//交易保障url
	PayitilReportUrl string

	//转换短链接url
	ShortUrlUrl string

	//拉取订单评价数据url
	BatchQueryCommentUrl string

	//
	Oauth2AuthorizeUrl string

	//
	Oauth2AccessTokenUrl string
}

var payapi = PayApi{
	BizPayUrlUrl:         "weixin://wxpay/bizpayurl",
	UnifiedOrderUrl:      "https://api.mch.weixin.qq.com/pay/unifiedorder",
	UnifiedOrderUrl2:     "https://api2.mch.weixin.qq.com/pay/unifiedorder",
	OrderQueryUrl:        "https://api.mch.weixin.qq.com/pay/orderquery",
	OrderQueryUrl2:       "https://api2.mch.weixin.qq.com/pay/orderquery",
	CloseOrderUrl:        "https://api.mch.weixin.qq.com/pay/closeorder",
	CloseOrderUrl2:       "https://api2.mch.weixin.qq.com/pay/closeorder",
	RefundUrl:            "https://api.mch.weixin.qq.com/secapi/pay/refund",
	RefundQueryUrl:       "https://api.mch.weixin.qq.com/pay/refundquery",
	DownloadBillUrl:      "https://api.mch.weixin.qq.com/pay/downloadbill",
	DownloadFundFlowUrl:  "https://api.mch.weixin.qq.com/pay/downloadfundflow",
	PayitilReportUrl:     "https://api.mch.weixin.qq.com/payitil/report",
	ShortUrlUrl:          "https://api.mch.weixin.qq.com/tools/shorturl",
	BatchQueryCommentUrl: "https://api.mch.weixin.qq.com/billcommentsp/batchquerycomment",
	Oauth2AuthorizeUrl:   "https://open.weixin.qq.com/connect/oauth2/authorize",
	Oauth2AccessTokenUrl: "https://api.weixin.qq.com/sns/oauth2/access_token",
}

func GetDefault() PayApi {
	return payapi
}

type Oauth2AuthorizeScope string

const (
	Oauth2AuthorizeScope_SnsapiBase     Oauth2AuthorizeScope = "snsapi_base"
	Oauth2AuthorizeScope_SnsapiUserInfo Oauth2AuthorizeScope = "snsapi_userinfo"
)

type ResponseError struct {
	ErrorCode int    `json:"errorcode"`
	ErrorMsg  string `json:"errormsg"`
}

type AccessToken struct {
	ResponseError
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

//网页授权（jsapi获取openid）:构建前端浏览器访问微信oath2 url
func (p PayApi) BuildOauth2AuthorizeUrl(mch Merchant, redirectUrl string, scope Oauth2AuthorizeScope) (reqUrl string, state string) {
	vals := make(url2.Values)
	vals["appid"] = append(vals["appid"], mch.AppId)
	vals["redirect_url"] = append(vals["redirect_url"], mch.AppId)
	vals["response_type"] = append(vals["response_type"], "code")
	vals["scope"] = append(vals["scope"], string(scope))
	state = internal.CreateRandomStr(20, internal.ALPHANUM)
	vals["state"] = append(vals["state"], state)

	reqUrl = fmt.Sprintf("%s?%s#wechat_redirect", p.Oauth2AuthorizeUrl, vals.Encode())

	return
}

//网页授权（jsapi获取openid）: 通过code获取AccessToken数据
func (p PayApi) Oauth2AuthorizeRedirectUrlHandler(mch Merchant, code string) (AccessToken, error) {
	vals := make(url2.Values)
	vals["appid"] = append(vals["appid"], mch.AppId)
	vals["secret"] = append(vals["secret"], mch.AppSecret)
	vals["code"] = append(vals["code"], code)
	vals["grant_type"] = append(vals["grant_type"], "authorization_code")
	reqUrl := fmt.Sprintf("%s?%s", p.Oauth2AccessTokenUrl, vals.Encode())

	ret := AccessToken{}
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return ret, err
	}
	var cli http.Client
	if resp, err := cli.Do(req); err != nil {
		return ret, err
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
			return ret, err
		} else {
			if ret.ErrorCode == 40029 {
				return ret, fmt.Errorf("%s", ret.ErrorMsg)
			} else {
				return ret, nil
			}
		}
	}
}

//模式一：生成二维码url
func (p PayApi) BizPayUrl1(mch Merchant, productId string) (string, error) {
	//验证参数
	if productId == "" {
		return "", fmt.Errorf("商品id不能为空")
	}
	if err := internal.CheckRequire(string(TradeTypeNative), mch); err != nil {
		return "", err
	}

	reqModel := BizPayUrl{
		AppId:     mch.AppId,
		MchId:     mch.MchId,
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  internal.CreateRandomStr(30, internal.ALPHANUM),
		ProductId: productId,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXml, mch.PrivateKey)

	codeUrl := fmt.Sprintf("%s?%s", p.BizPayUrlUrl, BuildPrameter(reqModel))

	return codeUrl, nil
}

//模式二：生成二维码url
func (p PayApi) BizPayUrl2(codeUrl string) string {
	codeUrl = fmt.Sprintf("%s?sr=%s", p.BizPayUrlUrl, codeUrl)
	return codeUrl
}

//转换短链接
func (p PayApi) ShortUrl(mch Merchant, url string) (string, error) {
	//验证参数
	if url == "" {
		return "", fmt.Errorf("url不能为空")
	}
	if err := internal.CheckRequire(string(TradeTypeNative), mch); err != nil {
		return "", err
	}

	reqModel := ReqShortUrl{
		AppId:    mch.AppId,
		MchId:    mch.MchId,
		LongUrl:  url,
		NonceStr: internal.GetNonceStr(),
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXml, mch.PrivateKey)

	//发起请求
	respModel := RespShortUrl{}
	respMap := make(XmlMap)
	if bs, err := request(&respModel, http.MethodPost, p.ShortUrlUrl, &reqModel, nil); err != nil {
		return "", err
	} else {
		if err := xml.Unmarshal(bs, &respMap); err != nil {
			return "", err
		}
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXml, mch.PrivateKey) {
			return "", fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return "", fmt.Errorf("%s", errMap[respModel.ErrCode])
		}

		//成功
		return respModel.ShortUrl, nil
	} else {
		//失败
		return "", fmt.Errorf("%s", respModel.ReturnMsg)
	}
}

//模式一：直接生成二维码短链接
func (p PayApi) BizPayShortUrl(mch Merchant, productId string) (string, error) {
	codeUrl, err := p.BizPayUrl1(mch, productId)
	if err != nil {
		return "", err
	}

	return p.ShortUrl(mch, codeUrl)
}

//统一下单
/**
 *返回参数依次为：预订单号，二维码url，错误
 */
func getNativePayReq(mch Merchant, order UnifiedOrder) ReqUnifiedOrder {
	return ReqUnifiedOrder{
		AppId:          mch.AppId,
		MchId:          mch.MchId,
		DeviceInfo:     "WEB",
		NonceStr:       internal.GetNonceStr(),
		Body:           order.Body,
		Detail:         order.Detail,
		Attach:         order.Attach,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       order.TotalFee,
		SpbillCreateIp: order.SpbillCreateIp,
		TimeStart:      time.Now().Format(timeFormat),
		TimeExpire:     order.PayExpire.Format(timeFormat),
		GoodsTag:       order.GoodsTag,
		NotifyUrl:      order.NotifyUrl,
		TradeType:      string(TradeTypeNative),
		ProductId:      order.ProductId,
		LimitPay:       order.LimitPay,
		Receipt:        "Y",
		SceneInfo:      order.SceneInfo,
	}
}

func getAppPayReq(mch Merchant, order UnifiedOrder) ReqUnifiedOrder {
	return ReqUnifiedOrder{
		AppId:          mch.AppId,
		MchId:          mch.MchId,
		NonceStr:       internal.GetNonceStr(),
		Body:           order.Body,
		Detail:         order.Detail,
		Attach:         order.Attach,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       order.TotalFee,
		SpbillCreateIp: order.SpbillCreateIp,
		TimeStart:      time.Now().Format(timeFormat),
		TimeExpire:     order.PayExpire.Format(timeFormat),
		GoodsTag:       order.GoodsTag,
		NotifyUrl:      order.NotifyUrl,
		TradeType:      string(TradeTypeApp),
		LimitPay:       order.LimitPay,
		Receipt:        "Y",
	}
}

func getMwebPayReq(mch Merchant, order UnifiedOrder) ReqUnifiedOrder {
	return ReqUnifiedOrder{
		AppId:          mch.AppId,
		MchId:          mch.MchId,
		NonceStr:       internal.GetNonceStr(),
		Body:           order.Body,
		Detail:         order.Detail,
		Attach:         order.Attach,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       order.TotalFee,
		SpbillCreateIp: order.SpbillCreateIp,
		TimeStart:      time.Now().Format(timeFormat),
		TimeExpire:     order.PayExpire.Format(timeFormat),
		GoodsTag:       order.GoodsTag,
		NotifyUrl:      order.NotifyUrl,
		TradeType:      string(TradeTypeMweb),
		ProductId:      order.ProductId,
		LimitPay:       order.LimitPay,
		Receipt:        "Y",
		SceneInfo:      order.SceneInfo,
	}
}

func getJsapiPayReq(mch Merchant, order UnifiedOrder) ReqUnifiedOrder {
	return ReqUnifiedOrder{
		AppId:          mch.AppId,
		MchId:          mch.MchId,
		NonceStr:       internal.GetNonceStr(),
		Body:           order.Body,
		Detail:         order.Detail,
		Attach:         order.Attach,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       order.TotalFee,
		SpbillCreateIp: order.SpbillCreateIp,
		TimeStart:      time.Now().Format(timeFormat),
		TimeExpire:     order.PayExpire.Format(timeFormat),
		GoodsTag:       order.GoodsTag,
		NotifyUrl:      order.NotifyUrl,
		TradeType:      string(TradeTypeJsapi),
		LimitPay:       order.LimitPay,
		OpenId:         order.OpendId,
		Receipt:        "Y",
		SceneInfo:      order.SceneInfo,
	}
}

func (p PayApi) UnifiedOrder(tradeType TradeType, mch Merchant, order UnifiedOrder) (ReturnUnifiedOrder, error) {
	//验证参数
	if err := internal.CheckRequire(string(TradeTypeNative), order); err != nil {
		return ReturnUnifiedOrder{}, err
	}
	if err := internal.CheckRequire(string(TradeTypeNative), mch); err != nil {
		return ReturnUnifiedOrder{}, err
	}

	reqModel := ReqUnifiedOrder{}
	switch tradeType {
	case TradeTypeApp:
		reqModel = getAppPayReq(mch, order)
	case TradeTypeMweb:
		reqModel = getMwebPayReq(mch, order)
	case TradeTypeJsapi:
		reqModel = getJsapiPayReq(mch, order)
	case TradeTypeNative:
		reqModel = getNativePayReq(mch, order)
	default:
		//默认native支付
		reqModel = getNativePayReq(mch, order)
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXml, mch.PrivateKey)

	//发起请求
	respModel := RespUnifiedOrder{}
	respMap := make(XmlMap)
	if bs, err := request(&respModel, http.MethodPost, p.UnifiedOrderUrl, &reqModel, nil); err != nil {
		return ReturnUnifiedOrder{}, err
	} else {
		if err := xml.Unmarshal(bs, &respMap); err != nil {
			fmt.Println(err)
			return ReturnUnifiedOrder{}, err
		}
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXml, mch.PrivateKey) {
			return ReturnUnifiedOrder{}, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return ReturnUnifiedOrder{}, fmt.Errorf("%s", respModel.ErrCodeDes)
		}

		if respModel.TradeType != string(tradeType) {
			return ReturnUnifiedOrder{}, fmt.Errorf("交易类型不匹配")
		}

		//成功
		returnModel := ReturnUnifiedOrder{
			TradeType: respModel.TradeType,
			PrepayId:  respModel.PrepayId,
			CodeUrl:   respModel.CodeUrl,
			MwebUrl:   respModel.MwebUrl,
		}
		return returnModel, nil
	} else {
		//失败
		return ReturnUnifiedOrder{}, fmt.Errorf("%s", respModel.ReturnMsg)
	}
}

//查询订单
func (p PayApi) OrderQuery(mch Merchant, transactionId, outTradeNo string) (OrderQuery, error) {
	order := OrderQuery{}

	//验证参数
	if transactionId == "" && outTradeNo == "" {
		return order, fmt.Errorf("微信订单号，商户订单号不能同时为空")
	}
	if err := internal.CheckRequire("", mch); err != nil {
		return order, err
	}

	reqModel := ReqOrderQuery{
		AppId:         mch.AppId,
		MchId:         mch.MchId,
		TransactionId: transactionId,
		OutTradeNo:    outTradeNo,
		NonceStr:      internal.GetNonceStr(),
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXml, mch.PrivateKey)

	//发起请求
	respModel := RespOrderQuery{}
	respMap := make(XmlMap)
	if bs, err := request(&respModel, http.MethodPost, p.OrderQueryUrl, &reqModel, nil); err != nil {
		return order, err
	} else {
		if err := xml.Unmarshal(bs, &respMap); err != nil {
			return order, err
		}
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXml, mch.PrivateKey) {
			return order, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return order, fmt.Errorf("%s", respModel.ErrCodeDes)
		}

		//成功
		order.DeviceInfo = respModel.DeviceInfo
		order.OpenId = respModel.OpenId
		order.IsSubscribe = respModel.IsSubscribe
		order.TradeType = respModel.TradeType
		order.TradeState = respModel.TradeState
		order.BankType = respModel.BankType
		order.TotalFee = respModel.TotalFee
		order.SettlementTotalFee = respModel.SettlementTotalFee
		order.FeeType = respModel.FeeType
		order.CashFee = respModel.CashFee
		order.CashFeeType = respModel.CashFeeType
		order.CouponFee = respModel.CouponFee
		order.CouponCount = respModel.CouponCount
		order.CouponType = respModel.CouponType
		order.CouponId = respModel.CouponId
		order.CouponFeen = respModel.CouponFeen
		order.TransactionId = respModel.TransactionId
		order.OutTradeNo = respModel.OutTradeNo
		order.Attach = respModel.Attach
		order.TimeEnd = respModel.TimeEnd
		order.TradeStateDesc = respModel.TradeStateDesc
		return order, nil
	} else {
		//失败
		return order, fmt.Errorf("%s", respModel.ReturnMsg)
	}
}

//关闭订单
func (p PayApi) CloseOrder(mch Merchant, outTradeNo string) error {
	//验证参数
	if outTradeNo == "" {
		return fmt.Errorf("商户订单号")
	}
	if err := internal.CheckRequire("", mch); err != nil {
		return err
	}

	reqModel := ReqCloseOrder{
		AppId:      mch.AppId,
		MchId:      mch.MchId,
		OutTradeNo: outTradeNo,
		NonceStr:   internal.GetNonceStr(),
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXml, mch.PrivateKey)

	//发起请求
	respModel := RespCloseOrder{}
	respMap := make(XmlMap)
	if bs, err := request(&respModel, http.MethodPost, p.CloseOrderUrl, &reqModel, nil); err != nil {
		return err
	} else {
		if err := xml.Unmarshal(bs, &respMap); err != nil {
			return err
		}
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXml, mch.PrivateKey) {
			return fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return fmt.Errorf("%s", respModel.ErrCodeDes)
		}

		//成功
		return nil
	} else {
		//失败
		return fmt.Errorf("%s", respModel.ReturnMsg)
	}
}

//申请退款
func (p PayApi) Refund(mch Merchant, refund Refund) (RefundReturn, error) {
	retn := RefundReturn{}

	//验证参数
	if err := internal.CheckRequire("", refund); err != nil {
		return retn, err
	}
	if err := internal.CheckRequire("", mch); err != nil {
		return retn, err
	}

	reqModel := ReqRefund{
		AppId:         mch.AppId,
		MchId:         mch.MchId,
		NonceStr:      internal.GetNonceStr(),
		TransactionId: refund.TransactionId,
		OutTradeNo:    refund.OutTradeNo,
		OutRefundNo:   refund.OutRefundNo,
		TotalFee:      refund.TotalFee,
		RefundFee:     refund.RefundFee,
		RefundDesc:    refund.RefundDesc,
		RefundAccount: refund.RefundAccount,
		NotifyUrl:     refund.NotifyUrl,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXml, mch.PrivateKey)

	//发起请求
	respModel := RespRefund{}
	respMap := make(XmlMap)
	tr := getPkcs12(mch.MchId, refund.CertP12)
	if bs, err := request(&respModel, http.MethodPost, p.RefundUrl, &reqModel, &tr); err != nil {
		return retn, err
	} else {
		if err := xml.Unmarshal(bs, &respMap); err != nil {
			return retn, err
		}
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXml, mch.PrivateKey) {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return retn, fmt.Errorf("%s", respModel.ErrCodeDes)
		}

		//成功
		retn.TransactionId = respModel.TransactionId
		retn.OutTradeNo = respModel.OutTradeNo
		retn.OutRefundNo = respModel.OutRefundNo
		retn.RefundId = respModel.RefundId
		retn.RefundFee = respModel.RefundFee
		retn.SettlementRefundFee = respModel.SettlementRefundFee
		retn.TotalFee = respModel.TotalFee
		retn.SettlementTotalFee = respModel.SettlementTotalFee
		retn.FeeType = respModel.FeeType
		retn.CashFee = respModel.CashFee
		retn.CashFeeType = respModel.CashFeeType
		retn.CashRefundFee = respModel.CashRefundFee
		retn.CouponType = respModel.CouponType
		retn.CouponRefundFee = respModel.CouponRefundFee
		retn.CouponRefundFeen = respModel.CouponRefundFeen
		retn.CouponRefundCount = respModel.CouponRefundCount
		retn.CouponRefundId = respModel.CouponRefundId
		return retn, nil
	} else {
		//失败
		return retn, fmt.Errorf("%s", respModel.ReturnMsg)
	}
}

//查询退款
func (p PayApi) RefundQuery(mch Merchant, refundQuery RefundQuery) (RefundQueryReturn, error) {
	retn := RefundQueryReturn{}

	//验证参数
	if refundQuery.OutRefundNo == "" &&
		refundQuery.OutTradeNo == "" &&
		refundQuery.TransactionId == "" &&
		refundQuery.RefundId == "" {
		return retn, fmt.Errorf("微信订单号，商户订单号，商户退款单号，微信退款单号不能全为空")
	}
	if err := internal.CheckRequire("", mch); err != nil {
		return retn, err
	}

	reqModel := ReqRefundQuery{
		AppId:         mch.AppId,
		MchId:         mch.MchId,
		NonceStr:      internal.GetNonceStr(),
		TransactionId: refundQuery.TransactionId,
		OutTradeNo:    refundQuery.OutTradeNo,
		OutRefundNo:   refundQuery.OutRefundNo,
		RefundId:      refundQuery.RefundId,
		Offset:        refundQuery.Offset,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXml, mch.PrivateKey)

	//发起请求
	respModel := RespRefundQuery{}
	respMap := make(XmlMap)
	if bs, err := request(&respModel, http.MethodPost, p.RefundQueryUrl, &reqModel, nil); err != nil {
		return retn, err
	} else {
		if err := xml.Unmarshal(bs, &respMap); err != nil {
			return retn, err
		}
		respModel.RefundStatusn = respMap[fmt.Sprintf("refund_status_%d", reqModel.Offset)]
		respModel.RefundAccountn = respMap[fmt.Sprintf("refund_account_%d", reqModel.Offset)]
		respModel.RefundRecvAccoutn = respMap[fmt.Sprintf("refund_recv_accout_%d", reqModel.Offset)]
		respModel.RefundSuccessTimen = respMap[fmt.Sprintf("refund_success_time_%d", reqModel.Offset)]
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXml, mch.PrivateKey) {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return retn, fmt.Errorf("%s", respModel.ErrCodeDes)
		}

		//成功
		retn.TotalRefundCount = respModel.TotalRefundCount
		retn.TransactionId = respModel.TransactionId
		retn.OutTradeNo = respModel.OutTradeNo
		retn.TotalFee = respModel.TotalFee
		retn.SettlementTotalFee = respModel.SettlementTotalFee
		retn.FeeType = respModel.FeeType
		retn.CashFee = respModel.CashFee
		retn.RefundCount = respModel.RefundCount
		retn.OutRefundNon = respModel.OutRefundNon
		retn.RefundIdn = respModel.RefundIdn
		retn.RefundChanneln = respModel.RefundChanneln
		retn.RefundFeen = respModel.RefundFeen
		retn.SettlementRefundFeen = respModel.SettlementRefundFeen
		retn.CouponTypenm = respModel.CouponTypenm
		retn.CouponRefundFeen = respModel.CouponRefundFeen
		retn.CouponRefundCountn = respModel.CouponRefundCountn
		retn.CouponRefundIdnm = respModel.CouponRefundIdnm
		retn.CouponRefundFeenm = respModel.CouponRefundFeenm
		retn.RefundStatusn = respModel.RefundStatusn
		retn.RefundAccountn = respModel.RefundAccountn
		retn.RefundRecvAccoutn = respModel.RefundRecvAccoutn
		retn.RefundSuccessTimen = respModel.RefundSuccessTimen
		return retn, nil
	} else {
		//失败
		return retn, fmt.Errorf("%s", respModel.ReturnMsg)
	}
}

//下载对账单
func unmashalToBillAll(items []string) (DownloadBillReturn, error) {
	bill := DownloadBillReturn{}
	if len(items) != 24 {
		return bill, fmt.Errorf("无法识别账单格式，请重新查看微信native开发文档")
	}
	bill.TradeDate = items[0][1:]
	bill.AppId = items[1][1:]
	bill.MchId = items[2][1:]
	bill.SubMchId = items[3][1:]
	bill.DeviceInfo = items[4][1:]
	bill.TransactionId = items[5][1:]
	bill.OutTradeNo = items[6][1:]
	bill.OpenId = items[7][1:]
	bill.TradeType = items[8][1:]
	bill.TradeState = items[9][1:]
	bill.BankType = items[10][1:]
	bill.FeeType = items[11][1:]
	bill.TotalFee = items[12][1:]
	bill.CouponFee = items[13][1:]
	bill.RefundIdn = items[14][1:]
	bill.OutRefundNon = items[15][1:]
	bill.SettlementRefundFeen = items[16][1:]
	bill.RefundCouponFee = items[17][1:]
	bill.RefundType = items[18][1:]
	bill.RefundState = items[19][1:]
	bill.ProductName = items[20][1:]
	bill.ProductBar = items[21][1:]
	bill.HandlingFee = items[22][1:]
	bill.Rate = items[23][1:]

	return bill, nil
}

func unmashalToBillSuccess(items []string) (DownloadBillReturn, error) {
	bill := DownloadBillReturn{}
	if len(items) != 18 {
		return bill, fmt.Errorf("无法识别账单格式，请重新查看微信native开发文档")
	}
	bill.TradeDate = items[0][1:]
	bill.AppId = items[1][1:]
	bill.MchId = items[2][1:]
	bill.SubMchId = items[3][1:]
	bill.DeviceInfo = items[4][1:]
	bill.TransactionId = items[5][1:]
	bill.OutTradeNo = items[6][1:]
	bill.OpenId = items[7][1:]
	bill.TradeType = items[8][1:]
	bill.TradeState = items[9][1:]
	bill.BankType = items[10][1:]
	bill.FeeType = items[11][1:]
	bill.TotalFee = items[12][1:]
	bill.CouponFee = items[13][1:]
	bill.ProductName = items[14][1:]
	bill.ProductBar = items[15][1:]
	bill.HandlingFee = items[16][1:]
	bill.Rate = items[17][1:]

	return bill, nil
}

func unmashalToBillRefund(items []string) (DownloadBillReturn, error) {
	bill := DownloadBillReturn{}
	if len(items) != 26 {
		return bill, fmt.Errorf("无法识别账单格式，请重新查看微信native开发文档")
	}
	bill.TradeDate = items[0][1:]
	bill.AppId = items[1][1:]
	bill.MchId = items[2][1:]
	bill.SubMchId = items[3][1:]
	bill.DeviceInfo = items[4][1:]
	bill.TransactionId = items[5][1:]
	bill.OutTradeNo = items[6][1:]
	bill.OpenId = items[7][1:]
	bill.TradeType = items[8][1:]
	bill.TradeState = items[9][1:]
	bill.BankType = items[10][1:]
	bill.FeeType = items[11][1:]
	bill.TotalFee = items[12][1:]
	bill.CouponFee = items[13][1:]
	bill.RefundApplyTimen = items[14][1:]
	bill.RefundSuccessTimen = items[15][1:]
	bill.RefundIdn = items[16][1:]
	bill.OutRefundNon = items[17][1:]
	bill.SettlementRefundFeen = items[18][1:]
	bill.RefundCouponFee = items[19][1:]
	bill.RefundType = items[20][1:]
	bill.RefundState = items[21][1:]
	bill.ProductName = items[22][1:]
	bill.ProductBar = items[23][1:]
	bill.HandlingFee = items[24][1:]
	bill.Rate = items[25][1:]

	return bill, nil
}

func parseBill(billByte []byte, billType, tarType string) ([]DownloadBillReturn, DownloadBillStatisticsReturn, error) {
	billStr := ""
	bills := []DownloadBillReturn{}
	billStatistics := DownloadBillStatisticsReturn{}

	if tarType == "GZIP" {
		buf := bytes.NewBuffer(billByte)
		zr, err := gzip.NewReader(buf)
		if err != nil {
			return bills, billStatistics, err
		}
		var uncompressBytes []byte
		if _, err := zr.Read(uncompressBytes); err != nil {
			return bills, billStatistics, err
		}

		billStr = string(uncompressBytes)
	} else {
		billStr = string(billByte)
	}

	entrys := strings.Split(billStr, "\n")
	isBillData := true
	for i := range entrys {
		if i == 0 {
			continue
		}

		if strings.Index(entrys[i], "`") == -1 {
			isBillData = false
		}

		if isBillData {
			//账单数据
			items := strings.Split(entrys[i], ",")
			switch billType {
			case "ALL":
				if bill, err := unmashalToBillAll(items); err != nil {
					return nil, billStatistics, err
				} else {
					bills = append(bills, bill)
				}
			case "SUCCESS":
				if bill, err := unmashalToBillSuccess(items); err != nil {
					return nil, billStatistics, err
				} else {
					bills = append(bills, bill)
				}
			case "REFUND", "RECHARGE_REFUND":
				if bill, err := unmashalToBillRefund(items); err != nil {
					return nil, billStatistics, err
				} else {
					bills = append(bills, bill)
				}
			default:
				return nil, billStatistics, fmt.Errorf("未识别的参数账单类型")
			}
		} else {
			//统计数据
			items := strings.Split(entrys[i], ",")
			billStatistics.TotalTransactions = items[0][1:]
			billStatistics.TotalTransactionValue = items[1][1:]
			billStatistics.TotalRefundValue = items[2][1:]
			billStatistics.TotalCouponRefundValue = items[3][1:]
			billStatistics.TotalHandlingFeeValue = items[4][1:]
		}

	}

	return bills, billStatistics, nil
}

func (p PayApi) DownloadBill(mch Merchant, dwnBill DownLoadBill) ([]DownloadBillReturn, DownloadBillStatisticsReturn, error) {
	billStatistics := DownloadBillStatisticsReturn{}

	//验证参数
	if err := internal.CheckRequire("", dwnBill); err != nil {
		return nil, billStatistics, err
	}
	if err := internal.CheckRequire("", mch); err != nil {
		return nil, billStatistics, err
	}
	if _, err := time.Parse(dateFormat, dwnBill.BillDate); err != nil {
		return nil, billStatistics, err
	}

	reqModel := ReqDownloadBill{
		AppId:    mch.AppId,
		MchId:    mch.MchId,
		NonceStr: internal.GetNonceStr(),
		BillDate: dwnBill.BillDate,
		BillType: dwnBill.BillType,
		TarType:  dwnBill.TarType,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXml, mch.PrivateKey)

	//发起请求
	respModel := RespDownloadBill{}
	body, err := request(&respModel, http.MethodPost, p.DownloadBillUrl, &reqModel, nil)
	if err != nil {
		return nil, billStatistics, err
	}

	if respModel.ReturnCode == ReturnCodeFail {
		//失败
		return nil, billStatistics, fmt.Errorf("%s", respModel.ReturnMsg)

	} else {
		//成功
		//解析数据流
		return parseBill(body, dwnBill.BillType, dwnBill.TarType)
	}
}

//交易保障
func (p PayApi) DownloadFundFlow(mch Merchant) error {
	//TODO
	return nil
}

//下载资金账单
func (p PayApi) PayitilReport(mch Merchant) error {
	//TODO
	return nil
}

//拉取订单评价数据
func (p PayApi) BatchQueryComment(mch Merchant) error {
	//TODO
	return nil
}

func getPkcs12(mchId string, p12 []byte) http.Transport {
	//rootCAs := x509.NewCertPool()
	//_, xcert, err := pkcs12.Decode(p12, mchId)
	//rootCAs.AddCert(xcert)

	pemBlock, err := pkcs12.ToPEM(p12, mchId)
	if err != nil {
		panic(err)
	}
	var pemData []byte
	for _, b := range pemBlock {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		panic(err)
	}
	return http.Transport{
		TLSClientConfig: &tls.Config{
			//RootCAs: rootCAs,
			Certificates: []tls.Certificate{cert},
		},
	}
}