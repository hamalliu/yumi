package db

import (
	"yumi/pkg/ecode"
	"yumi/pkg/external/dbc"
)

//WxPayMerchant ...
type WxPayMerchant struct {
	SeqID int64 `db:"seq_id" json:"seq_id"`
	//卖家key
	SellerKey string `db:"seller_key" json:"seller_key"`
	//开发平台唯一id
	AppID string `db:"app_id" json:"app_id"`
	//开发平台唯一id
	MchID string `db:"mch_id" json:"mch_id"`
	//私钥
	PrivateKey string `db:"private_key" json:"private_key"`
	//AppSecret（Secret）是APPID对应的接口密码，用于获取接口调用凭证access_token时使用
	Secret string `db:"secret" json:"secret"`
}

//GetWxPayMerchantBySellerKey ...
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

//GetWxPayMerchantByMchID ...
func GetWxPayMerchantByMchID(mchID string) (mch WxPayMerchant, err error) {
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
	if err = dbc.Get().Get(&mch, sqlStr, mchID); err != nil {
		return mch, ecode.ServerErr(err)
	}

	return
}
