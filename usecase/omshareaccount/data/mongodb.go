package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"yumi/pkg/status"
	"yumi/usecase/omshareaccount"
	"yumi/usecase/omshareaccount/entity"
)

// MongoDB ...
type MongoDB struct {
	*mongo.Database
}

var _ omshareaccount.Data = &MongoDB{}

// New ...
func New(db *mongo.Database) *MongoDB {
	return &MongoDB{Database: db}
}

func (db *MongoDB) collShareAccounts() *mongo.Collection {
	return db.Collection("share_accounts")
}

// ShareAccount ...
func (db *MongoDB) ShareAccount(saa entity.ShareAccountAttribute) omshareaccount.DataShareAccount {
	sa := ShareAccounts{db: db, ShareAccountAttribute: saa}
	return &sa
}

// GetShareAccount ...
func (db *MongoDB) GetShareAccount(shareID string) (omshareaccount.DataShareAccount, error) {
	coll := db.collShareAccounts()
	sa := ShareAccounts{db: db}

	ctx := context.Background()
	ret := coll.FindOne(ctx, primitive.M{"share_id": shareID})
	if ret.Err() != nil {
		return &sa, status.Internal().WithDetails(ret.Err().Error())
	}

	err := ret.Decode(&sa)
	if err != nil {
		return &sa, err
	}

	return &sa, nil
}

// ShareAccounts ...
type ShareAccounts struct {
	db *MongoDB `bson:"-"`

	ID primitive.ObjectID `bson:"_id"`
	entity.ShareAccountAttribute
}

var _ omshareaccount.DataShareAccount = &ShareAccounts{}

// Attribute ...
func (sa *ShareAccounts) Attribute() *entity.ShareAccountAttribute {
	return &sa.ShareAccountAttribute
}

// Create ...
func (sa *ShareAccounts) Create() error {
	coll := sa.db.collShareAccounts()

	ctx := context.Background()
	_, err := coll.InsertOne(ctx, sa)
	if err != nil {
		return status.Internal().WithDetails(err.Error())
	}

	return nil
}

// Update ...
func (sa *ShareAccounts) Update() error {
	coll := sa.db.collShareAccounts()

	ctx := context.Background()
	_, err := coll.UpdateOne(ctx, primitive.M{"_id": sa.ID}, sa)
	if err != nil {
		return status.Internal().WithDetails(err.Error())
	}
	return nil
}
