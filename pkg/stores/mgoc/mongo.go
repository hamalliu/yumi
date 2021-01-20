package mgoc

import (
	"context"
	"errors"
	"fmt"
	"time"
	"yumi/pkg/conf"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client ...
type Client struct {
	*mongo.Client
}

// New ...
func New(uri string, opts ...*options.ClientOptions) (*Client, error) {
	opts = append(opts, options.Client().ApplyURI(uri))
	cli, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = cli.Connect(ctx)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &Client{Client: cli}, nil
}

var _mgoc *Client

// Init ...
func Init(config conf.Mongo) error {
	if config.Dsn == "" {
		panic(errors.New("mongo dsn is empty"))
	}

	opts := []*options.ClientOptions{}
	if config.ConnectTimeout != 0 {
		opts = append(opts, options.Client().SetConnectTimeout(config.ConnectTimeout.Duration()))
	}
	if config.HeartbeatInterval != 0 {
		opts = append(opts, options.Client().SetHeartbeatInterval(config.HeartbeatInterval.Duration()))
	}
	if config.MaxConnIdleTime != 0 {
		opts = append(opts, options.Client().SetMaxConnIdleTime(config.MaxConnIdleTime.Duration()))
	}
	if config.MaxPoolSize != 0 {
		opts = append(opts, options.Client().SetMaxPoolSize(config.MaxPoolSize))
	}
	if config.MinPoolSize != 0 {
		opts = append(opts, options.Client().SetMinPoolSize(config.MinPoolSize))
	}

	cli, err := New(config.Dsn, opts...)
	if err != nil {
		panic(err)
	}

	_mgoc = cli
	return nil
}

// Get ...
func Get() *Client {
	return _mgoc
}

// Connect ...
func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/yumi"))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	type UserInfo struct {
		ID     primitive.ObjectID `bson:"_id"`
		Name   string
		Age    uint16
		Weight uint32
	}
	query := make(primitive.M)
	query["name"] = "xxx"
	cursor, err := client.Database("yumi").Collection("admin").Find(ctx, query)
	if err != nil {
		fmt.Println(err.Error())
	}

	uis := []UserInfo{}
	fmt.Println(cursor.All(ctx, &uis))
	fmt.Println(uis)

	ret := client.Database("yumi").Collection("admin").FindOne(ctx, primitive.M{"_id": uis[0].ID})
	if err != nil {
		fmt.Println(err.Error(), 1)
	}

	uis2 := make(map[string]interface{})
	fmt.Println(ret.Decode(&uis2), 2)
	fmt.Println(uis2)

	// fmt.Println(client.Database("yumi").CreateCollection(ctx, "media"))
	return nil
}
