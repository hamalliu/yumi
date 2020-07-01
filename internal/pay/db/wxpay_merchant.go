package db

import (
	"yumi/pkg/ecode"
	"yumi/pkg/external/dbc"
)

type WxPayMerchant struct {
	SeqId      int64  `db:"seq_id" json:"seq_id"`
	SellerKey  string `db:"seller_key" json:"seller_key"`   //卖家key
	AppId      string `db:"app_id" json:"app_id"`           //开发平台唯一id
	MchId      string `db:"mch_id" json:"mch_id"`           //开发平台唯一id
	PrivateKey string `db:"private_key" json:"private_key"` //私钥
	Secret     string `db:"secret" json:"secret"`           //AppSecret（Secret）是APPID对应的接口密码，用于获取接口调用凭证access_token时使用
}

func GetWxPayMerchantBySellerKey(sellerKey string) (mch WxPayMerchant, err error) {
	sqlStr := `
		SELECT 
			seq_id, 
			seller_key, 
			app_id, 
			mch_id,
			secret,
			private_key 
		FROM 
			wx_pay_merchants 
		WHERE 
			seller_key = ?`
	if err = dbc.Get().Get(&mch, sqlStr, sellerKey); err != nil {
		return mch, ecode.ServerErr(err)
	}
	return
}

func GetWxPayMerchantByMchId(mchId string) (mch WxPayMerchant, err error) {
	sqlStr := `
		SELECT 
			seq_id, 
			seller_key, 
			app_id, 
			mch_id,
			secret,
			private_key 
		FROM 
			wx_pay_merchants 
		WHERE 
			mch_id = ?`
	if err = dbc.Get().Get(&mch, sqlStr, mchId); err != nil {
		return mch, ecode.ServerErr(err)
	}

	return
}
