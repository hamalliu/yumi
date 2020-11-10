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

	"yumi/pkg/trade/internal"
)

/**
 *微信支付
 */

const (
	//FieldTagKeyXML ...
	FieldTagKeyXML = "xml"
	//FieldTagKeyJSON ...
	FieldTagKeyJSON = "json"
)

const (
	//ReturnCodeSuccess ...
	ReturnCodeSuccess = "SUCCESS"
	//ReturnCodeFail ...
	ReturnCodeFail = "FAIL"
)

const (
	//TradeStateSuccess ...
	TradeStateSuccess = "SUCCESS"
	//TradeStateRefund ...
	TradeStateRefund = "REFUND"
	//TradeStateNotpay ...
	TradeStateNotpay = "NOTPAY"
	//TradeStateClosed ...
	TradeStateClosed = "CLOSED"
	//TradeStateRevoked ...
	TradeStateRevoked = "REVOKED"
	//TradeStateUserPaying ...
	TradeStateUserPaying = "USERPAYING"
	//TradeStatePayError ...
	TradeStatePayError = "PAYERROR"
	//TradeStateRefundClose ...
	TradeStateRefundClose = "REFUNDCLOSE"
	//TradeStateProcessing ...
	TradeStateProcessing = "PROCESSING"
	//TradeStateChange ...
	TradeStateChange = "CHANGE"
)

const (
	timeFormat = "20060102150405"
	dateFormat = "20060102"
)

//TradeType ...
type TradeType string

const (
	//TradeTypeJsapi ...
	TradeTypeJsapi TradeType = "JSAPI"
	//TradeTypeNative ...
	TradeTypeNative TradeType = "NATIVE"
	//TradeTypeMweb ...
	TradeTypeMweb TradeType = "MWEB"
	//TradeTypeApp ...
	TradeTypeApp TradeType = "APP"
)

//PayAPI ...
type PayAPI struct {
	//生成支付二维码url
	BizPayURLURL string

	//统一下单url
	UnifiedOrderURL string

	//统一下单url（备用域名）
	UnifiedOrderURL2 string

	//查询订单url
	OrderQueryURL string

	//查询订单url（备用域名）
	OrderQueryURL2 string

	//关闭订单url
	CloseOrderURL string

	//关闭订单url（备用域名）
	CloseOrderURL2 string

	//申请退款url
	RefundURL string

	//查询退款url
	RefundQueryURL string

	//下载对账单url
	DownloadBillURL string

	//下载资金账单url
	DownloadFundFlowURL string

	//交易保障url
	PayitilReportURL string

	//转换短链接url
	ShortURLURL string

	//拉取订单评价数据url
	BatchQueryCommentURL string

	//
	Oauth2AuthorizeURL string

	//
	Oauth2AccessTokenURL string
}

var payapi = PayAPI{
	BizPayURLURL:         "weixin://wxpay/bizpayurl",
	UnifiedOrderURL:      "https://api.mch.weixin.qq.com/pay/unifiedorder",
	UnifiedOrderURL2:     "https://api2.mch.weixin.qq.com/pay/unifiedorder",
	OrderQueryURL:        "https://api.mch.weixin.qq.com/pay/orderquery",
	OrderQueryURL2:       "https://api2.mch.weixin.qq.com/pay/orderquery",
	CloseOrderURL:        "https://api.mch.weixin.qq.com/pay/closeorder",
	CloseOrderURL2:       "https://api2.mch.weixin.qq.com/pay/closeorder",
	RefundURL:            "https://api.mch.weixin.qq.com/secapi/pay/refund",
	RefundQueryURL:       "https://api.mch.weixin.qq.com/pay/refundquery",
	DownloadBillURL:      "https://api.mch.weixin.qq.com/pay/downloadbill",
	DownloadFundFlowURL:  "https://api.mch.weixin.qq.com/pay/downloadfundflow",
	PayitilReportURL:     "https://api.mch.weixin.qq.com/payitil/report",
	ShortURLURL:          "https://api.mch.weixin.qq.com/tools/shorturl",
	BatchQueryCommentURL: "https://api.mch.weixin.qq.com/billcommentsp/batchquerycomment",
	Oauth2AuthorizeURL:   "https://open.weixin.qq.com/connect/oauth2/authorize",
	Oauth2AccessTokenURL: "https://api.weixin.qq.com/sns/oauth2/access_token",
}

//GetDefault ...
func GetDefault() PayAPI {
	return payapi
}

//Oauth2AuthorizeScope ...
type Oauth2AuthorizeScope string

const (
	//Oauth2AuthorizeScopeSnsapiBase ...
	Oauth2AuthorizeScopeSnsapiBase Oauth2AuthorizeScope = "snsapi_base"
	//Oauth2AuthorizeScopeSnsapiUserInfo ...
	Oauth2AuthorizeScopeSnsapiUserInfo Oauth2AuthorizeScope = "snsapi_userinfo"
)

//ResponseError ...
type ResponseError struct {
	ErrorCode int    `json:"errorcode"`
	ErrorMsg  string `json:"errormsg"`
}

//AccessToken ...
type AccessToken struct {
	ResponseError
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

//BuildOauth2AuthorizeURL 网页授权（jsapi获取openid）:构建前端浏览器访问微信oath2 url
func (p PayAPI) BuildOauth2AuthorizeURL(mch Merchant, redirectURL string, scope Oauth2AuthorizeScope) (reqURL string, state string) {
	vals := make(url2.Values)
	vals["appid"] = append(vals["appid"], mch.AppID)
	vals["redirect_url"] = append(vals["redirect_url"], mch.AppID)
	vals["response_type"] = append(vals["response_type"], "code")
	vals["scope"] = append(vals["scope"], string(scope))
	state = internal.CreateRandomStr(20, internal.ALPHANUM)
	vals["state"] = append(vals["state"], state)

	reqURL = fmt.Sprintf("%s?%s#wechat_redirect", p.Oauth2AuthorizeURL, vals.Encode())

	return
}

//Oauth2AuthorizeRedirectURLHandler 网页授权（jsapi获取openid）: 通过code获取AccessToken数据
func (p PayAPI) Oauth2AuthorizeRedirectURLHandler(mch Merchant, code string) (AccessToken, error) {
	vals := make(url2.Values)
	vals["appid"] = append(vals["appid"], mch.AppID)
	vals["secret"] = append(vals["secret"], mch.AppSecret)
	vals["code"] = append(vals["code"], code)
	vals["grant_type"] = append(vals["grant_type"], "authorization_code")
	reqURL := fmt.Sprintf("%s?%s", p.Oauth2AccessTokenURL, vals.Encode())

	ret := AccessToken{}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return ret, err
	}
	var cli http.Client
	resp, err := cli.Do(req)
	if err != nil {
		return ret, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return ret, err
	}
	if ret.ErrorCode == 40029 {
		return ret, fmt.Errorf("%s", ret.ErrorMsg)
	}

	return ret, nil
}

//BizPayURL1 模式一：生成二维码url
func (p PayAPI) BizPayURL1(mch Merchant, productID string) (string, error) {
	//验证参数
	if productID == "" {
		return "", fmt.Errorf("商品id不能为空")
	}
	if err := internal.CheckRequire(string(TradeTypeNative), mch); err != nil {
		return "", err
	}

	reqModel := BizPayURL{
		AppID:     mch.AppID,
		MchID:     mch.MchID,
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  internal.CreateRandomStr(30, internal.ALPHANUM),
		ProductID: productID,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXML, mch.PrivateKey)

	codeURL := fmt.Sprintf("%s?%s", p.BizPayURLURL, BuildPrameter(reqModel))

	return codeURL, nil
}

//BizPayURL2 模式二：生成二维码url
func (p PayAPI) BizPayURL2(codeURL string) string {
	codeURL = fmt.Sprintf("%s?sr=%s", p.BizPayURLURL, codeURL)
	return codeURL
}

//ShortURL 转换短链接
func (p PayAPI) ShortURL(mch Merchant, url string) (string, error) {
	//验证参数
	if url == "" {
		return "", fmt.Errorf("url不能为空")
	}
	if err := internal.CheckRequire(string(TradeTypeNative), mch); err != nil {
		return "", err
	}

	reqModel := ReqShortURL{
		AppID:    mch.AppID,
		MchID:    mch.MchID,
		LongURL:  url,
		NonceStr: internal.GetNonceStr(),
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXML, mch.PrivateKey)

	//发起请求
	respModel := RespShortURL{}
	respMap := make(XMLMap)
	bs, err := request(&respModel, http.MethodPost, p.ShortURLURL, &reqModel, nil)
	if err != nil {
		return "", err
	}
	if err := xml.Unmarshal(bs, &respMap); err != nil {
		return "", err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXML, mch.PrivateKey) {
			return "", fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return "", fmt.Errorf("%s", errMap[respModel.ErrCode])
		}

		//成功
		return respModel.ShortURL, nil
	}

	//失败
	return "", fmt.Errorf("%s", respModel.ReturnMsg)
}

//BizPayShortURL 模式一：直接生成二维码短链接
func (p PayAPI) BizPayShortURL(mch Merchant, productID string) (string, error) {
	codeURL, err := p.BizPayURL1(mch, productID)
	if err != nil {
		return "", err
	}

	return p.ShortURL(mch, codeURL)
}

//统一下单
/**
 *返回参数依次为：预订单号，二维码url，错误
 */
func getNativePayReq(mch Merchant, order UnifiedOrder) ReqUnifiedOrder {
	return ReqUnifiedOrder{
		AppID:          mch.AppID,
		MchID:          mch.MchID,
		DeviceInfo:     "WEB",
		NonceStr:       internal.GetNonceStr(),
		Body:           order.Body,
		Detail:         order.Detail,
		Attach:         order.Attach,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       order.TotalFee,
		SpbillCreateIP: order.SpbillCreateIP,
		TimeStart:      time.Now().Format(timeFormat),
		TimeExpire:     order.PayExpire.Format(timeFormat),
		GoodsTag:       order.GoodsTag,
		NotifyURL:      order.NotifyURL,
		TradeType:      string(TradeTypeNative),
		ProductID:      order.ProductID,
		LimitPay:       order.LimitPay,
		Receipt:        "Y",
		SceneInfo:      CData(order.SceneInfo),
	}
}

func getAppPayReq(mch Merchant, order UnifiedOrder) ReqUnifiedOrder {
	return ReqUnifiedOrder{
		AppID:          mch.AppID,
		MchID:          mch.MchID,
		NonceStr:       internal.GetNonceStr(),
		Body:           order.Body,
		Detail:         order.Detail,
		Attach:         order.Attach,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       order.TotalFee,
		SpbillCreateIP: order.SpbillCreateIP,
		TimeStart:      time.Now().Format(timeFormat),
		TimeExpire:     order.PayExpire.Format(timeFormat),
		GoodsTag:       order.GoodsTag,
		NotifyURL:      order.NotifyURL,
		TradeType:      string(TradeTypeApp),
		LimitPay:       order.LimitPay,
		Receipt:        "Y",
	}
}

func getMwebPayReq(mch Merchant, order UnifiedOrder) ReqUnifiedOrder {
	return ReqUnifiedOrder{
		AppID:          mch.AppID,
		MchID:          mch.MchID,
		NonceStr:       internal.GetNonceStr(),
		Body:           order.Body,
		Detail:         order.Detail,
		Attach:         order.Attach,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       order.TotalFee,
		SpbillCreateIP: order.SpbillCreateIP,
		TimeStart:      time.Now().Format(timeFormat),
		TimeExpire:     order.PayExpire.Format(timeFormat),
		GoodsTag:       order.GoodsTag,
		NotifyURL:      order.NotifyURL,
		TradeType:      string(TradeTypeMweb),
		ProductID:      order.ProductID,
		LimitPay:       order.LimitPay,
		Receipt:        "Y",
		SceneInfo:      CData(order.SceneInfo),
	}
}

func getJsapiPayReq(mch Merchant, order UnifiedOrder) ReqUnifiedOrder {
	return ReqUnifiedOrder{
		AppID:          mch.AppID,
		MchID:          mch.MchID,
		NonceStr:       internal.GetNonceStr(),
		Body:           order.Body,
		Detail:         order.Detail,
		Attach:         order.Attach,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       order.TotalFee,
		SpbillCreateIP: order.SpbillCreateIP,
		TimeStart:      time.Now().Format(timeFormat),
		TimeExpire:     order.PayExpire.Format(timeFormat),
		GoodsTag:       order.GoodsTag,
		NotifyURL:      order.NotifyURL,
		TradeType:      string(TradeTypeJsapi),
		LimitPay:       order.LimitPay,
		OpenID:         order.OpendID,
		Receipt:        "Y",
		SceneInfo:      CData(order.SceneInfo),
	}
}

//UnifiedOrder ...
func (p PayAPI) UnifiedOrder(tradeType TradeType, mch Merchant, order UnifiedOrder) (ReturnUnifiedOrder, error) {
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
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXML, mch.PrivateKey)

	//发起请求
	respModel := RespUnifiedOrder{}
	respMap := make(XMLMap)
	bs, err := request(&respModel, http.MethodPost, p.UnifiedOrderURL, &reqModel, nil)
	if err != nil {
		return ReturnUnifiedOrder{}, err
	}
	if err := xml.Unmarshal(bs, &respMap); err != nil {
		fmt.Println(err)
		return ReturnUnifiedOrder{}, err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXML, mch.PrivateKey) {
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
			PrepayID:  respModel.PrepayID,
			CodeURL:   respModel.CodeURL,
			MwebURL:   respModel.MwebURL,
		}
		return returnModel, nil
	}
	//失败
	return ReturnUnifiedOrder{}, fmt.Errorf("%s", respModel.ReturnMsg)
}

//OrderQuery 查询订单
func (p PayAPI) OrderQuery(mch Merchant, transactionID, outTradeNo string) (OrderQuery, error) {
	order := OrderQuery{}

	//验证参数
	if transactionID == "" && outTradeNo == "" {
		return order, fmt.Errorf("微信订单号，商户订单号不能同时为空")
	}
	if err := internal.CheckRequire("", mch); err != nil {
		return order, err
	}

	reqModel := ReqOrderQuery{
		AppID:         mch.AppID,
		MchID:         mch.MchID,
		TransactionID: transactionID,
		OutTradeNo:    outTradeNo,
		NonceStr:      internal.GetNonceStr(),
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXML, mch.PrivateKey)

	//发起请求
	respModel := RespOrderQuery{}
	respMap := make(XMLMap)
	bs, err := request(&respModel, http.MethodPost, p.OrderQueryURL, &reqModel, nil)
	if err != nil {
		return order, err
	}
	if err := xml.Unmarshal(bs, &respMap); err != nil {
		return order, err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXML, mch.PrivateKey) {
			return order, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return order, fmt.Errorf("%s", respModel.ErrCodeDes)
		}

		//成功
		order.DeviceInfo = respModel.DeviceInfo
		order.OpenID = respModel.OpenID
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
		order.CouponID = respModel.CouponID
		order.CouponFeen = respModel.CouponFeen
		order.TransactionID = respModel.TransactionID
		order.OutTradeNo = respModel.OutTradeNo
		order.Attach = respModel.Attach
		order.TimeEnd = respModel.TimeEnd
		order.TradeStateDesc = respModel.TradeStateDesc
		return order, nil
	}
	//失败
	return order, fmt.Errorf("%s", respModel.ReturnMsg)
}

//CloseOrder 关闭订单
func (p PayAPI) CloseOrder(mch Merchant, outTradeNo string) error {
	//验证参数
	if outTradeNo == "" {
		return fmt.Errorf("商户订单号")
	}
	if err := internal.CheckRequire("", mch); err != nil {
		return err
	}

	reqModel := ReqCloseOrder{
		AppID:      mch.AppID,
		MchID:      mch.MchID,
		OutTradeNo: outTradeNo,
		NonceStr:   internal.GetNonceStr(),
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXML, mch.PrivateKey)

	//发起请求
	respModel := RespCloseOrder{}
	respMap := make(XMLMap)
	bs, err := request(&respModel, http.MethodPost, p.CloseOrderURL, &reqModel, nil)
	if err != nil {
		return err
	}
	if err := xml.Unmarshal(bs, &respMap); err != nil {
		return err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXML, mch.PrivateKey) {
			return fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return fmt.Errorf("%s", respModel.ErrCodeDes)
		}

		//成功
		return nil
	}
	//失败
	return fmt.Errorf("%s", respModel.ReturnMsg)
}

//Refund 申请退款
func (p PayAPI) Refund(mch Merchant, refund Refund) (RefundReturn, error) {
	retn := RefundReturn{}

	//验证参数
	if err := internal.CheckRequire("", refund); err != nil {
		return retn, err
	}
	if err := internal.CheckRequire("", mch); err != nil {
		return retn, err
	}

	reqModel := ReqRefund{
		AppID:         mch.AppID,
		MchID:         mch.MchID,
		NonceStr:      internal.GetNonceStr(),
		TransactionID: refund.TransactionID,
		OutTradeNo:    refund.OutTradeNo,
		OutRefundNo:   refund.OutRefundNo,
		TotalFee:      refund.TotalFee,
		RefundFee:     refund.RefundFee,
		RefundDesc:    refund.RefundDesc,
		RefundAccount: refund.RefundAccount,
		NotifyURL:     refund.NotifyURL,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXML, mch.PrivateKey)

	//发起请求
	respModel := RespRefund{}
	respMap := make(XMLMap)
	tr := getPkcs12(mch.MchID, refund.CertP12)
	bs, err := request(&respModel, http.MethodPost, p.RefundURL, &reqModel, &tr)
	if err != nil {
		return retn, err
	}
	if err := xml.Unmarshal(bs, &respMap); err != nil {
		return retn, err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXML, mch.PrivateKey) {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return retn, fmt.Errorf("%s", respModel.ErrCodeDes)
		}

		//成功
		retn.TransactionID = respModel.TransactionID
		retn.OutTradeNo = respModel.OutTradeNo
		retn.OutRefundNo = respModel.OutRefundNo
		retn.RefundID = respModel.RefundID
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
		retn.CouponRefundID = respModel.CouponRefundID
		return retn, nil
	}
	//失败
	return retn, fmt.Errorf("%s", respModel.ReturnMsg)
}

//RefundQuery 查询退款
func (p PayAPI) RefundQuery(mch Merchant, refundQuery RefundQuery) (RefundQueryReturn, error) {
	retn := RefundQueryReturn{}

	//验证参数
	if refundQuery.OutRefundNo == "" &&
		refundQuery.OutTradeNo == "" &&
		refundQuery.TransactionID == "" &&
		refundQuery.RefundID == "" {
		return retn, fmt.Errorf("微信订单号，商户订单号，商户退款单号，微信退款单号不能全为空")
	}
	if err := internal.CheckRequire("", mch); err != nil {
		return retn, err
	}

	reqModel := ReqRefundQuery{
		AppID:         mch.AppID,
		MchID:         mch.MchID,
		NonceStr:      internal.GetNonceStr(),
		TransactionID: refundQuery.TransactionID,
		OutTradeNo:    refundQuery.OutTradeNo,
		OutRefundNo:   refundQuery.OutRefundNo,
		RefundID:      refundQuery.RefundID,
		Offset:        refundQuery.Offset,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXML, mch.PrivateKey)

	//发起请求
	respModel := RespRefundQuery{}
	respMap := make(XMLMap)
	bs, err := request(&respModel, http.MethodPost, p.RefundQueryURL, &reqModel, nil)
	if err != nil {
		return retn, err
	}
	if err := xml.Unmarshal(bs, &respMap); err != nil {
		return retn, err
	}
	respModel.RefundStatusn = respMap[fmt.Sprintf("refund_status_%d", reqModel.Offset)]
	respModel.RefundAccountn = respMap[fmt.Sprintf("refund_account_%d", reqModel.Offset)]
	respModel.RefundRecvAccoutn = respMap[fmt.Sprintf("refund_recv_accout_%d", reqModel.Offset)]
	respModel.RefundSuccessTimen = respMap[fmt.Sprintf("refund_success_time_%d", reqModel.Offset)]

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respMap, FieldTagKeyXML, mch.PrivateKey) {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return retn, fmt.Errorf("%s", respModel.ErrCodeDes)
		}

		//成功
		retn.TotalRefundCount = respModel.TotalRefundCount
		retn.TransactionID = respModel.TransactionID
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
	}
	//失败
	return retn, fmt.Errorf("%s", respModel.ReturnMsg)
}

//下载对账单
func unmashalToBillAll(items []string) (DownloadBillReturn, error) {
	bill := DownloadBillReturn{}
	if len(items) != 24 {
		return bill, fmt.Errorf("无法识别账单格式，请重新查看微信native开发文档")
	}
	bill.TradeDate = items[0][1:]
	bill.AppID = items[1][1:]
	bill.MchID = items[2][1:]
	bill.SubMchID = items[3][1:]
	bill.DeviceInfo = items[4][1:]
	bill.TransactionID = items[5][1:]
	bill.OutTradeNo = items[6][1:]
	bill.OpenID = items[7][1:]
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
	bill.AppID = items[1][1:]
	bill.MchID = items[2][1:]
	bill.SubMchID = items[3][1:]
	bill.DeviceInfo = items[4][1:]
	bill.TransactionID = items[5][1:]
	bill.OutTradeNo = items[6][1:]
	bill.OpenID = items[7][1:]
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
	bill.AppID = items[1][1:]
	bill.MchID = items[2][1:]
	bill.SubMchID = items[3][1:]
	bill.DeviceInfo = items[4][1:]
	bill.TransactionID = items[5][1:]
	bill.OutTradeNo = items[6][1:]
	bill.OpenID = items[7][1:]
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
				bill, err := unmashalToBillAll(items)
				if err != nil {
					return nil, billStatistics, err
				}
				bills = append(bills, bill)
			case "SUCCESS":
				bill, err := unmashalToBillSuccess(items)
				if err != nil {
					return nil, billStatistics, err
				}
				bills = append(bills, bill)
			case "REFUND", "RECHARGE_REFUND":
				bill, err := unmashalToBillRefund(items)
				if err != nil {
					return nil, billStatistics, err
				}
				bills = append(bills, bill)
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

//DownloadBill ...
func (p PayAPI) DownloadBill(mch Merchant, dwnBill DownLoadBill) ([]DownloadBillReturn, DownloadBillStatisticsReturn, error) {
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
		AppID:    mch.AppID,
		MchID:    mch.MchID,
		NonceStr: internal.GetNonceStr(),
		BillDate: dwnBill.BillDate,
		BillType: dwnBill.BillType,
		TarType:  dwnBill.TarType,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, FieldTagKeyXML, mch.PrivateKey)

	//发起请求
	respModel := RespDownloadBill{}
	body, err := request(&respModel, http.MethodPost, p.DownloadBillURL, &reqModel, nil)
	if err != nil {
		return nil, billStatistics, err
	}

	if respModel.ReturnCode == ReturnCodeFail {
		//失败
		return nil, billStatistics, fmt.Errorf("%s", respModel.ReturnMsg)

	}
	//成功
	//解析数据流
	return parseBill(body, dwnBill.BillType, dwnBill.TarType)
}

//DownloadFundFlow 交易保障
func (p PayAPI) DownloadFundFlow(mch Merchant) error {
	//TODO
	return nil
}

//PayitilReport 下载资金账单
func (p PayAPI) PayitilReport(mch Merchant) error {
	//TODO
	return nil
}

//BatchQueryComment 拉取订单评价数据
func (p PayAPI) BatchQueryComment(mch Merchant) error {
	//TODO
	return nil
}

func getPkcs12(mchID string, p12 []byte) http.Transport {
	//rootCAs := x509.NewCertPool()
	//_, xcert, err := pkcs12.Decode(p12, mchID)
	//rootCAs.AddCert(xcert)

	pemBlock, err := pkcs12.ToPEM(p12, mchID)
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
