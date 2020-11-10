package main

import (
	"yumi/pkg/conf"
	"yumi/pkg/stores/dbc"
)

func main() {
	dbc.Init(conf.DB{Dsn: "", DBName: "", MaxIdleConns: 10, MaxOpenConns: 10, ConnMaxLifetime: 30,})
	//sql.MediaCreateTable()
}
