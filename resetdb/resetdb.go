package main

import (
	"yumi/pkg/conf"
	"yumi/pkg/external/dbc"
	"yumi/resetdb/sql"
)

func main() {
	dbc.Init(conf.DBConfig{Dsn: "", DBName: "", MaxIdleConns: 10, MaxOpenConns: 10, ConnMaxLifetime: 30})
	sql.SysmmngCreateTable()
	//sql.MediaCreateTable()
}
