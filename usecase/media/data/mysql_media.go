package data

import (
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/usecase/media/entity"
)

// MysqlDB ...
type MysqlDB struct {
	*mysqlx.Client
}

// New ...
func New(db *mysqlx.Client) *MysqlDB {
	return &MysqlDB{Client: db}
}

//Create ...
func (db *MysqlDB) Create(ma entity.MediaAttribute) error {
	return nil
}

//Get ...
func (db *MysqlDB) Get(fileNo string) (ma entity.MediaAttribute, err error) {
	return
}

//List ...
func (db *MysqlDB) List(page, line int) (mas []entity.MediaAttribute, err error) {
	return
}
