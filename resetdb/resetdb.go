package main

import (
	"yumi/pkg/external/dbc"
	"yumi/pkg/external/dbc/mysqlx"
	"yumi/resetdb/sql"
)

func main() {
	dbc.Init(mysqlx.Config{Dsn: "", DBName: "", MaxIdleConns: 10, MaxOpenConns: 10, ConnMaxLifetime: 30})
	sql.SysmmngCreateTable()
	//sql.MediaCreateTable()
}
