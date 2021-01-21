package usecase

import (
	trade_db "yumi/usecase/trade/db"
	trade_platform "yumi/usecase/trade/platform"
	user_db "yumi/usecase/user/data"

	"go.mongodb.org/mongo-driver/mongo"
)

// Init trade db
func Init(db *mongo.Database) {
	trade_db.Init()
	trade_platform.Init()

	user_db.Init(db)
}
