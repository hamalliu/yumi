package data

import (
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/usecase/trade"
	"yumi/usecase/trade/thirdpf/alipay"
	"yumi/usecase/trade/thirdpf/wxpay"
)

// Init trade db
func Init(db *mysqlx.Client) {
	trade.InitData(New(db))
	wxpay.InitData(New(db))
	alipay.InitData(New(db))
}
