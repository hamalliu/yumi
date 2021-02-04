package data

import (
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/usecase/trade"
	"yumi/usecase/trade/thirdpf"
)

// Init trade db
func Init(db *mysqlx.Client) {
	trade.InitData(New(db))
	thirdpf.InitData(New(db))
}
