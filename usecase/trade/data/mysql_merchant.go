package data

import (
	"yumi/usecase/trade/entity"
	"yumi/usecase/trade/thirdpf/wxpay"

	"github.com/pkg/errors"
)

// CreateWxPayMerchant ...
func (db *MysqlDB) CreateWxPayMerchant(entity.WxPayMerchant) error {
	return nil
}

// GetWxPayMerchant ...
func (db *MysqlDB) GetWxPayMerchant(req wxpay.FilterMerchantIDs) (mch entity.WxPayMerchant, err error) {
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
	if err = db.Get(&mch, sqlStr, req.SellerKey); err != nil {
		return mch, errors.WithStack(err)
	}

	return
}

// CreateAliPayMerchant ...
func (db *MysqlDB) CreateAliPayMerchant(entity.AliPayMerchant) error {
	return nil
}

// GetAliPayMerchant ...
func (db *MysqlDB) GetAliPayMerchant(sellerKey string) (mch entity.AliPayMerchant, err error) {
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
	if err = db.Get(&mch, sqlStr, sellerKey); err != nil {
		return mch, errors.WithStack(err)
	}
	return
}
