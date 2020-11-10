package dbc

import (
	"time"

	"yumi/pkg/conf"
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/pkg/types"
)

var mysqlDb *mysqlx.Client

//Init 初始化数据库连接配置
func Init(conf conf.DB) {
	var err error
	if mysqlDb, err = mysqlx.New(conf); err != nil {
		panic(err)
	}
}

//InitDefault 按默认配置初始化数据库连接
func InitDefault() {
	testConf := conf.DB{
		Dsn:             "",
		DBName:          "",
		MaxOpenConns:    10,
		MaxIdleConns:    10,
		ConnMaxLifetime: types.TimeDuration(2 * time.Hour),
	}
	var err error
	if mysqlDb, err = mysqlx.New(testConf); err != nil {
		panic(err)
	}
}

//Get 返回数据库客户端
func Get() *mysqlx.Client {
	return mysqlDb
}
