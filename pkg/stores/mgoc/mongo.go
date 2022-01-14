package mgoc

import (
	"context"
	"time"

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
