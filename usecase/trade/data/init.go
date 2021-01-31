package data

import (
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/usecase/trade"
)

// Init trade db
func Init(db *mysqlx.Client) {
	trade.InitData(New(db))
}
