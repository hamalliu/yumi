package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"yumi/pkg/sessions"
	"yumi/pkg/stores/mgoc"
	"yumi/usecase/user"
	"yumi/usecase/user/entity"
)

// MongoCli ...
type MongoCli struct {
	*mgoc.Client
}

var _ user.Data = &MongoCli{}

// New ...
func New(cli *mgoc.Client) *MongoCli {
	return &MongoCli{Client: cli}
}

type MongoTX struct {
	*MongoCli
	*mgoc.MongoTX
}

func (tx MongoTX) Ctx() context.Context {
	return tx.Sctx
}

func (cli *MongoCli) NewTx() (*MongoTX, error) {
	sess, err := cli.StartSession()
	if err != nil {
		return nil, err
	}

	sctx := mongo.NewSessionContext(cli.Ctx(), sess)

	return &MongoTX{MongoCli: cli, MongoTX: &mgoc.MongoTX{Sctx: sctx}}, nil
}

func (cli *MongoCli) collUsers() *mongo.Collection {
	return cli.Database("yumi").Collection("users")
}

func (cli *MongoCli) Ctx() context.Context {
	return context.Background()
}

// Create ...
func (cli *MongoCli) Create(ua entity.UserAttribute) error {
	coll := cli.collUsers()

	_, err := coll.InsertOne(cli.Ctx(), ua)
	if err != nil {
		return err
	}

	return nil
}

// Update ...
func (cli *MongoCli) Update(ua entity.UserAttribute) error {
	coll := cli.collUsers()

	_, err := coll.ReplaceOne(cli.Ctx(), primitive.M{"user_uuid": ua.UserUUID}, ua)
	if err != nil {
		return err
	}
	return nil
}

// Get ...
func (cli *MongoCli) Get(ids entity.UserAttributeIDs) (ua entity.UserAttribute, err error) {
	coll := cli.collUsers()

	ret := coll.FindOne(cli.Ctx(), cli.filterIDs(ids))
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
func (cli *MongoCli) Exist(ids entity.UserAttributeIDs) (exist bool, ua entity.UserAttribute, err error) {
	coll := cli.collUsers()

	ret := coll.FindOne(cli.Ctx(), cli.filterIDs(ids))
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

func (cli *MongoCli) filterIDs(ids entity.UserAttributeIDs) map[string]interface{} {
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
func (cli *MongoCli) GetSessionsStore() sessions.Store {
	return nil
}
