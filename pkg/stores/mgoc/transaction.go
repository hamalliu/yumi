package mgoc

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"yumi/pkg/log"
)

type MongoTX struct {
	Sess mongo.Session
	sctx mongo.SessionContext
}

func (tx *MongoTX) Ctx() context.Context {
	return tx.sctx
}

func (tx *MongoTX) Start() error {
	tx.sctx = mongo.NewSessionContext(context.TODO(), tx.Sess)

	if err := tx.Sess.StartTransaction(options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))); err != nil {
		return err
	}

	return nil
}

func (tx *MongoTX) Rollback() error {
	return tx.Sess.AbortTransaction(context.TODO())
}

func (tx *MongoTX) Commit() error {
	for {
		err := tx.Sess.CommitTransaction(context.TODO())
		switch e := err.(type) {
		case nil:
			return nil
		case mongo.CommandError:
			if e.HasErrorLabel("UnknownTransactionCommitResult") {
				log.Error("UnknownTransactionCommitResult, retrying commit operation...")
				continue
			}
			log.Error("Error during commit...")
			return e
		default:
			log.Error("Error during commit...")
			return e
		}
	}
}

func (tx *MongoTX) End() {
	tx.Sess.EndSession(context.TODO())
}
