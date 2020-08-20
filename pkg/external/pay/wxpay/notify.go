package wxpay

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"yumi/pkg/external/pay/internal"
)

//CheckPrePayNotify ...
func CheckPrePayNotify(mch Merchant, req ReqPrepayNotify) error {
	if mch.AppID != req.AppID {
		return fmt.Errorf("公众号不一致")
	}
	if mch.MchID != req.MchID {
		return fmt.Errorf("商户号不一致")
	}

	//验签
	if req.Sign != Buildsign(req, FieldTagKeyXML, mch.PrivateKey) {
		return fmt.Errorf("签名错误")
	}

	return nil
}

//CheckPayNotify ...
func CheckPayNotify(mch Merchant, totalFee int, outTradeNo string, req ReqPayNotify) error {
	if req.ReturnCode == ReturnCodeSuccess {
		if mch.AppID != req.AppID {
			return fmt.Errorf("公众号不一致")
		}
		if mch.MchID != req.MchID {
			return fmt.Errorf("商户号不一致")
		}

		//验签
		if req.Sign != Buildsign(req, FieldTagKeyXML, mch.PrivateKey) {
			return fmt.Errorf("签名错误")
		}

		//业务结果
		if req.ResultCode != ReturnCodeSuccess {
			return fmt.Errorf("%s", req.ErrCodeDes)
		}

		//验证商户订单号
		if outTradeNo != req.OutTradeNo {
			return fmt.Errorf("商户订单号不一致")
		}

		//验证订单
		if totalFee != req.TotalFee {
			return fmt.Errorf("订单金额不一致")
		}

		if req.TradeType != "NATIVE" {
			return fmt.Errorf("交易类型错误")
		}

		return nil
	}
	
	return fmt.Errorf("%s", req.ReturnMsg)
}

//CheckRefundNotify ...
func CheckRefundNotify(mch Merchant, req ReqRefundNotify) error {
	if req.ReturnCode == ReturnCodeSuccess {
		if mch.AppID != req.AppID {
			return fmt.Errorf("公众号不一致")
		}
		if mch.MchID != req.MchID {
			return fmt.Errorf("商户号不一致")
		}
		return nil
	}

	return fmt.Errorf("%s", req.ReturnMsg)
}

//GetRefundNotify ...
func GetRefundNotify(req *http.Request) (ReqRefundNotify, error) {
	ret := ReqRefundNotify{}
	err := xml.NewDecoder(req.Body).Decode(&ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

//DecryptoRefundNotify ...
func DecryptoRefundNotify(mch Merchant, info string) (ReqRefundNotifyEncryptInfo, error) {
	ret := ReqRefundNotifyEncryptInfo{}
	infoBytes, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(info))))
	if err != nil {
		return ret, err
	}
	md5ctx := md5.New()
	md5ctx.Write([]byte(mch.PrivateKey))
	key := strings.ToLower(hex.EncodeToString(md5ctx.Sum(nil)))

	info = internal.AesDecrypt(string(infoBytes), []byte(key))
	if err := xml.Unmarshal([]byte(info), &ret); err != nil {
		return ret, err
	}

	return ret, nil
}
