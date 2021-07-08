package data

import (
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/usecase/media/service"
)

// MysqlDB ...
type MysqlDB struct {
	*mysqlx.Client
}

var _ service.Data = &MysqlDB{}

// New ...
func New(db *mysqlx.Client) *MysqlDB {
	return &MysqlDB{Client: db}
}
