package alipay

import (
	"fmt"
)

func CheckPayNotify(mch Merchant, sellerId, outTradeNo, totalAmount string, req ReqNotify) error {
	if err := NotifyVerify(req, req.Sign, mch.PublicKey); err != nil {
		return fmt.Errorf("签名错误")
	}

	if mch.AppId != req.AppId {
		return fmt.Errorf("开发应用id不一致")
	}

	if outTradeNo != req.OutTradeNo {
		return fmt.Errorf("商户订单号不一致")
	}

	if totalAmount != req.TotalAmount {
		return fmt.Errorf("订单金额不一致")
	}

	if sellerId != req.SellerId {
		return fmt.Errorf("支付宝卖家用户号不一致")
	}

	return nil
}
