package usecase

import (
	"yumi/pkg/stores/dbc/mysqlx"
	tradedata "yumi/usecase/trade/data"
	thirdpf "yumi/usecase/trade/thirdpf"
	userdata "yumi/usecase/user/data"

	"go.mongodb.org/mongo-driver/mongo"
)

// Init trade db
func Init(mongoDB *mongo.Database, mysqlDB *mysqlx.Client) {
	tradedata.Init(mysqlDB)
	userdata.Init(mongoDB)

	thirdpf.Init()
}
