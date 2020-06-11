package dbc

import (
	"time"
	"yumi/pkg/conf"
	"yumi/pkg/external/dbc/mysqlx"
)

var mysqlDb *mysqlx.Model

func Init(conf conf.DB) {
	var err error
	if mysqlDb, err = mysqlx.New(conf); err != nil {
		panic(err)
	}
}

func InitDefault() {
	testConf := conf.DB{
		Dsn:          "",
		DBName:       "",
		MaxOpenConns: 10,
		MaxIdleConns: 10,
		ConnMaxLifetime: conf.TimeDuration{
			Duration: 2 * time.Hour,
		},
	}
	var err error
	if mysqlDb, err = mysqlx.New(testConf); err != nil {
		panic(err)
	}
}

func Get() *mysqlx.Model {
	return mysqlDb
}
