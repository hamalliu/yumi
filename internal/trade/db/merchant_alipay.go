package db

import (
	"yumi/pkg/ecode"
	"yumi/pkg/external/dbc"
)

//AliPayMerchant ...
type AliPayMerchant struct {
	SeqID int64 `db:"seq_id" json:"seq_id"`
	//商户id
	AppID string `db:"app_id" json:"app_id"`
	//卖家key
	SellerKey string `db:"seller_key" json:"seller_key"`
	//私钥
	PrivateKey string `db:"private_key" json:"private_key"`
	//公钥
	PublicKey string `db:"public_key" json:"public_key"`
}

//GetAliPayMerchantBySellerKey ...
func GetAliPayMerchantBySellerKey(sellerKey string) (mch AliPayMerchant, err error) {
	sqlStr := `
		SELECT 
			seq_id, 
			app_id,
			seller_key, 
			private_key, 
			public_key 
		FROM 
			ali_pay_merchants 
		WHERE 
			seller_key = ?`
	if err = dbc.Get().Get(&mch, sqlStr, sellerKey); err != nil {
		return mch, ecode.ServerErr(err)
	}

	return
}
