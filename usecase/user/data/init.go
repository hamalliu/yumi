package data

import (
	"go.mongodb.org/mongo-driver/mongo"

	"yumi/usecase/user"
)

// Init ...
func Init(db *mongo.Database) {
	user.InitData(New(db))
}
