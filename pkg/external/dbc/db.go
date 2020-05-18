package dbc

import (
	"yumi/pkg/conf"
	"yumi/pkg/external/dbc/mysqlx"
)

var mysqlDb *mysqlx.Model

func Init(conf conf.DBConfig) {
	var err error
	if mysqlDb, err = mysqlx.New(conf); err != nil {
		panic(err)
	}
}

func InitDefault() {
	testConf := conf.DBConfig{
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
