package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"yumi/pkg/sessions"
	"yumi/pkg/status"
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

// User ...
func (db *MongoDB) User(saa entity.UserAttribute) user.DataUser {
	sa := Users{db: db, UserAttribute: saa}
	return &sa
}

// GetUser ...
func (db *MongoDB) GetUser(userID string) (user.DataUser, error) {
	coll := db.collUsers()
	sa := Users{db: db}

	ctx := context.Background()
	ret := coll.FindOne(ctx, primitive.M{"user_id": userID})
	if ret.Err() != nil {
		return &sa, status.Internal().WithDetails(ret.Err().Error())
	}

	err := ret.Decode(&sa)
	if err != nil {
		return &sa, err
	}

	return &sa, nil
}

// GetSessionsStore ...
func (db *MongoDB) GetSessionsStore() sessions.Store {
	return nil
}

// Users ...
type Users struct {
	db *MongoDB `bson:"-"`

	ID primitive.ObjectID `bson:"_id"`
	entity.UserAttribute
}

var _ user.DataUser = &Users{}

// Attribute ...
func (sa *Users) Attribute() *entity.UserAttribute {
	return &sa.UserAttribute
}

// Create ...
func (sa *Users) Create() error {
	coll := sa.db.collUsers()

	ctx := context.Background()
	_, err := coll.InsertOne(ctx, sa)
	if err != nil {
		return status.Internal().WithDetails(err.Error())
	}

	return nil
}

// Update ...
func (sa *Users) Update() error {
	coll := sa.db.collUsers()

	ctx := context.Background()
	_, err := coll.UpdateOne(ctx, primitive.M{"_id": sa.ID}, sa)
	if err != nil {
		return status.Internal().WithDetails(err.Error())
	}
	return nil
}
