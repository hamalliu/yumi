package db

import (
	"yumi/pkg/ecode"
	"yumi/pkg/external/dbc"
)

type AliPayMerchant struct {
	SeqId      int64  `db:"seq_id" json:"seq_id"`
	AppId      string `db:"app_id" json:"seller_key"`       //商户id
	SellerKey  string `db:"seller_key" json:"seller_key"`   //卖家key
	PrivateKey string `db:"private_key" json:"private_key"` //私钥
	PublicKey  string `db:"public_key" json:"public_key"`   //公钥
}

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
