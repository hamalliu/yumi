package usecase

import (
	"yumi/pkg/stores/dbc/mysqlx"
	trade "yumi/usecase/trade/data"
	thirdpf "yumi/usecase/trade/thirdpf"
	user "yumi/usecase/user/data"

	"go.mongodb.org/mongo-driver/mongo"
)

// Init trade db
func Init(db *mongo.Database, mysqlDB *mysqlx.Client) {
	trade.Init(mysqlDB)
	user.Init(db)

	thirdpf.Init()
}
