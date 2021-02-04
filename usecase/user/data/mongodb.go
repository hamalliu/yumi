package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"yumi/pkg/sessions"
	"yumi/usecase/user"
	"yumi/usecase/user/entity"
)

// MongoDB ...
type MongoDB struct {
	*mongo.Database
}

var _ user.Data = &MongoDB{}

// New ...
func New(db *mongo.Database) *MongoDB {
	return &MongoDB{Database: db}
}

func (db *MongoDB) collUsers() *mongo.Collection {
	return db.Collection("users")
}

// Create ...
func (db *MongoDB) Create(ua entity.UserAttribute) error {
	coll := db.collUsers()

	ctx := context.Background()
	_, err := coll.InsertOne(ctx, ua)
	if err != nil {
		return err
	}

	return nil
}

// Update ...
func (db *MongoDB) Update(ua entity.UserAttribute) error {
	coll := db.collUsers()

	ctx := context.Background()
	_, err := coll.ReplaceOne(ctx, primitive.M{"user_uuid": ua.UserUUID}, ua)
	if err != nil {
		return err
	}
	return nil
}

// GetUser ...
func (db *MongoDB) GetUser(userID string) (ua entity.UserAttribute, err error) {
	coll := db.collUsers()

	ctx := context.Background()
	ret := coll.FindOne(ctx, primitive.M{"user_id": userID})
	if ret.Err() != nil {
		err = ret.Err()
		return
	}

	err = ret.Decode(&ua)
	if err != nil {
		return
	}

	return
}

// GetSessionsStore ...
func (db *MongoDB) GetSessionsStore() sessions.Store {
	return nil
}
