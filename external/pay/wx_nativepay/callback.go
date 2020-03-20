package wx_nativepay

import "fmt"


func CheckPrePayNotify(mch Merchant, req ReqPrepayNotify) error {
	if mch.AppId != req.AppId {
		return fmt.Errorf("公众号不一致")
	}
	if mch.MchId != req.MchId {
		return fmt.Errorf("商户号不一致")
	}

	//验签
	if req.Sign != Buildsign(req, mch.PrivateKey) {
		return fmt.Errorf("签名错误")
	}

	return nil
}

func CheckPayNotify(mch Merchant, totalFee int, outTradeNo string, req ReqPayNotify) error {
	if req.ReturnCode == ReturnCodeSuccess {
		if mch.AppId != req.AppId {
			return fmt.Errorf("公众号不一致")
		}
		if mch.MchId != req.MchId {
			return fmt.Errorf("商户号不一致")
		}

		//验签
		if req.Sign != Buildsign(req, mch.PrivateKey) {
			return fmt.Errorf("签名错误")
		}

		//业务结果
		if req.ResultCode != ReturnCodeSuccess {
			return fmt.Errorf("%s", req.ErrCodeDesc)
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
	} else {
		return fmt.Errorf("%s", req.ReturnMsg)
	}
}

func CheckRefundNotify(mch Merchant, req ReqRefundNotify) error {
	if req.ReturnCode == ReturnCodeSuccess {
		if mch.AppId != req.AppId {
			return fmt.Errorf("公众号不一致")
		}
		if mch.MchId != req.MchId {
			return fmt.Errorf("商户号不一致")
		}
		return nil

	} else {
		return fmt.Errorf("%s", req.ReturnMsg)
	}
}
