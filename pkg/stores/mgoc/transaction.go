package mgoc

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"yumi/pkg/log"
)

// MongoTX ...
type MongoTX struct {
	Sctx mongo.SessionContext
}

// Start ...
func (tx *MongoTX) Start() error {
	err := tx.Sctx.StartTransaction(options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
	)
	if err != nil {
		return err
	}

	return nil
}

// Rollback ...
func (tx *MongoTX) Rollback() error {
	return tx.Sctx.AbortTransaction(tx.Sctx)
}

// Commit ...
func (tx *MongoTX) Commit() error {
	for {
		err := tx.Sctx.CommitTransaction(tx.Sctx)
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
