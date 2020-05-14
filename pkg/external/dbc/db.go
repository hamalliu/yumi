package dbc

import (
	"yumi/pkg/external/dbc/mysqlx"
)

var mysqlDb *mysqlx.Model

func Init(conf mysqlx.Config) {
	var err error
	if mysqlDb, err = mysqlx.New(conf); err != nil {
		panic(err)
	}
}

func InitDefault() {
	testConf := mysqlx.Config{
		Dsn:             "",
		DBName:          "",
		MaxOpenConns:    10,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30,
	}
	var err error
	if mysqlDb, err = mysqlx.New(testConf); err != nil {
		panic(err)
	}
}

func Get() *mysqlx.Model {
	return mysqlDb
}
