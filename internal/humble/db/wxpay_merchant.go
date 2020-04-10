package db

import (
	"yumi/external/dbc"
	"yumi/response"
)

type WxPayMerchant struct {
	SeqId      int64  `db:"seq_id" json:"seq_id"`
	SellerKey  string `db:"seller_key" json:"seller_key"`   //卖家key
	AppId      string `db:"app_id" json:"seller_key"`       //开发平台唯一id
	MchId      string `db:"app_id" json:"seller_key"`       //商户id
	PrivateKey string `db:"private_key" json:"private_key"` //私钥
}

func GetWxPayMerchantBySellerKey(sellerKey string) (mch WxPayMerchant, err error) {
	sqlStr := `
		SELECT 
			seq_id, 
			seller_key,
			app_id, 
			mch_id, 
			private_key 
		FROM 
			ali_pay_merchants 
		WHERE 
			seller_key = ?`
	if err = dbc.Get().Get(&mch, sqlStr, sellerKey); err != nil {
		return mch, response.InternalError(err)
	}

	return
}
