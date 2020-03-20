package wx_nativepay

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
	"time"

	"yumi/external/pay"
)

/**
 *微信native支付
 */

const (
	ReturnCodeSuccess = "SUCCESS"
	ReturnCodeFail    = "FAIL"
)

const (
	TradeStateSuccess    = "SUCCESS"
	TradeStateRefund     = "REFUND"
	TradeStateNotpay     = "NOTPAY"
	TradeStateClosed     = "CLOSED"
	TradeStateRevoked    = "REVOKED"
	TradeStateUserPaying = "USERPAYING"
	TradeStatePayError   = "PAYERROR"
)

const (
	timeFormat = "20060102150405"
	dateFormat = "20060102"
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

	//订单有效期（分钟）
	OrderValidity int
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
	OrderValidity:        30,
}

func GetDefault() PayApi {
	return payapi
}

//模式一：生成二维码url
func (p PayApi) BizPayUrl1(mch Merchant, productId string) (string, error) {
	//验证参数
	if productId == "" {
		return "", fmt.Errorf("商品id不能为空")
	}
	if err := pay.CheckRequire(mch); err != nil {
		return "", err
	}

	reqModel := BizPayUrl{
		AppId:     mch.AppId,
		MchId:     mch.MchId,
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  pay.CreateRandomStr(30, pay.ALPHANUM),
		ProductId: productId,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, mch.PrivateKey)

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
	if err := pay.CheckRequire(mch); err != nil {
		return "", err
	}

	reqModel := ReqShortUrl{
		AppId:    mch.AppId,
		MchId:    mch.MchId,
		LongUrl:  url,
		NonceStr: pay.GetNonceStr(),
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, mch.PrivateKey)

	//发起请求
	respModel := RespShortUrl{}
	if _, err := request(&respModel, http.MethodPost, p.ShortUrlUrl, &reqModel); err != nil {
		return "", err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respModel, mch.PrivateKey) {
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
func (p PayApi) UnifiedOrder(mch Merchant, order UnifiedOrder) (string, string, error) {
	//验证参数
	if err := pay.CheckRequire(order); err != nil {
		return "", "", err
	}
	if err := pay.CheckRequire(mch); err != nil {
		return "", "", err
	}

	reqModel := ReqUnifiedOrder{
		AppId:          mch.AppId,
		MchId:          mch.MchId,
		DeviceInfo:     "WEB",
		NonceStr:       pay.GetNonceStr(),
		Body:           order.Body,
		Detail:         order.Detail,
		Attach:         order.Attach,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       order.TotalFee,
		SpbillCreateIp: "", //TODO 本机ip
		TimeStart:      time.Now().Format(timeFormat),
		TimeExpire:     time.Now().Add(time.Minute * time.Duration(p.OrderValidity)).Format(timeFormat),
		GoodsTag:       order.GoodsTag,
		NotifyUrl:      order.NotifyUrl,
		TradeType:      "NATIVE",
		ProductId:      order.ProductId,
		LimitPay:       order.LimitPay,
		Receipt:        "Y",
		SceneInfo:      order.SceneInfo,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, mch.PrivateKey)

	//发起请求
	respModel := RespUnifiedOrder{}
	if _, err := request(&respModel, http.MethodPost, p.ShortUrlUrl, &reqModel); err != nil {
		return "", "", err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respModel, mch.PrivateKey) {
			return "", "", fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return "", "", fmt.Errorf("%s", respModel.ErrCodeDesc)
		}

		//成功
		return respModel.PrepayId, respModel.CodeUrl, nil
	} else {
		//失败
		return "", "", fmt.Errorf("%s", respModel.ReturnMsg)
	}
}

//查询订单
func (p PayApi) OrderQuery(mch Merchant, transactionId, outTradeNo string) (OrderQuery, error) {
	order := OrderQuery{}

	//验证参数
	if transactionId == "" && outTradeNo == "" {
		return order, fmt.Errorf("微信订单号，商户订单号不能同时为空")
	}
	if err := pay.CheckRequire(mch); err != nil {
		return order, err
	}

	reqModel := ReqOrderQuery{
		AppId:         mch.AppId,
		MchId:         mch.MchId,
		TransactionId: transactionId,
		OutTradeNo:    outTradeNo,
		NonceStr:      pay.GetNonceStr(),
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, mch.PrivateKey)

	//发起请求
	respModel := RespOrderQuery{}
	if _, err := request(&respModel, http.MethodPost, p.ShortUrlUrl, &reqModel); err != nil {
		return order, err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respModel, mch.PrivateKey) {
			return order, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return order, fmt.Errorf("%s", respModel.ErrCodeDesc)
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
	if err := pay.CheckRequire(mch); err != nil {
		return err
	}

	reqModel := ReqCloseOrder{
		AppId:      mch.AppId,
		MchId:      mch.MchId,
		OutTradeNo: outTradeNo,
		NonceStr:   pay.GetNonceStr(),
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, mch.PrivateKey)

	//发起请求
	respModel := RespCloseOrder{}
	if _, err := request(&respModel, http.MethodPost, p.ShortUrlUrl, &reqModel); err != nil {
		return err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respModel, mch.PrivateKey) {
			return fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return fmt.Errorf("%s", respModel.ErrCodeDesc)
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
	if err := pay.CheckRequire(refund); err != nil {
		return retn, err
	}
	if err := pay.CheckRequire(mch); err != nil {
		return retn, err
	}

	reqModel := ReqRefund{
		AppId:         mch.AppId,
		MchId:         mch.MchId,
		NonceStr:      pay.GetNonceStr(),
		TransactionId: refund.TransactionId,
		OutTradeNo:    refund.OutTradeNo,
		OutRefundNo:   mch.BuildOutRefundNo(),
		TotalFee:      refund.TotalFee,
		RefundFee:     refund.RefundFee,
		RefundDesc:    refund.RefundDesc,
		RefundAccount: refund.RefundAccount,
		NotifyUrl:     refund.NotifyUrl,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, mch.PrivateKey)

	//发起请求
	respModel := RespRefund{}
	if _, err := request(&respModel, http.MethodPost, p.ShortUrlUrl, &reqModel); err != nil {
		return retn, err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respModel, mch.PrivateKey) {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return retn, fmt.Errorf("%s", respModel.ErrCodeDesc)
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
	if err := pay.CheckRequire(mch); err != nil {
		return retn, err
	}

	reqModel := ReqRefundQuery{
		AppId:         mch.AppId,
		MchId:         mch.MchId,
		NonceStr:      pay.GetNonceStr(),
		TransactionId: refundQuery.TransactionId,
		OutTradeNo:    refundQuery.OutTradeNo,
		OutRefundNo:   refundQuery.OutRefundNo,
		RefundId:      refundQuery.RefundId,
		Offset:        refundQuery.Offset,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, mch.PrivateKey)

	//发起请求
	respModel := RespRefundQuery{}
	if _, err := request(&respModel, http.MethodPost, p.ShortUrlUrl, &reqModel); err != nil {
		return retn, err
	}

	if respModel.ReturnCode == ReturnCodeSuccess {
		//验签
		if respModel.Sign != Buildsign(respModel, mch.PrivateKey) {
			return retn, fmt.Errorf("验签失败")
		}

		//验证业务结果
		if respModel.ResultCode != ReturnCodeSuccess {
			return retn, fmt.Errorf("%s", respModel.ErrCodeDesc)
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
func (p PayApi) DownloadBill(mch Merchant, dwnBill DownLoadBill) ([]DownloadBillReturn, DownloadBillStatisticsReturn, error) {
	billStatistics := DownloadBillStatisticsReturn{}

	//验证参数
	if err := pay.CheckRequire(dwnBill); err != nil {
		return nil, billStatistics, err
	}
	if err := pay.CheckRequire(mch); err != nil {
		return nil, billStatistics, err
	}
	if _, err := time.Parse(dateFormat, dwnBill.BillDate); err != nil {
		return nil, billStatistics, err
	}

	reqModel := ReqDownloadBill{
		AppId:    mch.AppId,
		MchId:    mch.MchId,
		NonceStr: pay.GetNonceStr(),
		BillDate: dwnBill.BillDate,
		BillType: dwnBill.BillType,
		TarType:  dwnBill.TarType,
	}
	//生成签名
	reqModel.Sign = Buildsign(reqModel, mch.PrivateKey)

	//发起请求
	respModel := RespDownloadBill{}
	body, err := request(&respModel, http.MethodPost, p.ShortUrlUrl, &reqModel)
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
				if bill, err := unmashalToAllSuccess(items); err != nil {
					return nil, billStatistics, err
				} else {
					bills = append(bills, bill)
				}
			case "REFUND", "RECHARGE_REFUND":
				if bill, err := unmashalToAllRefund(items); err != nil {
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

func unmashalToAllSuccess(items []string) (DownloadBillReturn, error) {
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

func unmashalToAllRefund(items []string) (DownloadBillReturn, error) {
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

