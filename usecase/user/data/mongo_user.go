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

// Get ...
func (db *MongoDB) Get(ids entity.UserAttributeIDs) (ua entity.UserAttribute, err error) {
	coll := db.collUsers()

	ctx := context.Background()
	ret := coll.FindOne(ctx, db.filterIDs(ids))
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

// Exist ...
func (db *MongoDB) Exist(ids entity.UserAttributeIDs) (exist bool, ua entity.UserAttribute, err error) {
	coll := db.collUsers()

	ctx := context.Background()
	ret := coll.FindOne(ctx, db.filterIDs(ids))
	if ret.Err() != nil {
		err = ret.Err()
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return
	}

	err = ret.Decode(&ua)
	if err != nil {
		return
	}

	return true, ua, nil
}

func (db *MongoDB) filterIDs(ids entity.UserAttributeIDs) map[string]interface{} {
	filter := make(map[string]interface{})
	if ids.UserUUID != "" {
		filter["user_uuid"] = ids.UserUUID
		return filter
	}
	if ids.UserID != "" {
		filter["user_id"] = ids.UserID
		return filter
	}
	if ids.PhoneNumber != "" {
		filter["phone_number"] = ids.PhoneNumber
		return filter
	}

	return filter
}

// GetSessionsStore ...
func (db *MongoDB) GetSessionsStore() sessions.Store {
	return nil
}
