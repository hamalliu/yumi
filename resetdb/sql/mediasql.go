package sql

import (
	"yumi/pkg/external/dbc"
)

func MediaCreateTable() {

	//管理附件表
	backMediaSql := `CREATE TABLE back_medias
					(
					"id" int primary key,
					"suffix" Nvarchar(32) default(''),
					"name" Nvarchar(512) default(''),
					"real_name" Nvarchar(512) default(''),
					"path" Nvarchar(1024) default(''),
					"operator" Nvarchar(255) default(''),
					"operator_id" Nvarchar(255) default(''),
					"operate_time" datetime default now()
					)`
	if _, err := dbc.Get().Exec(backMediaSql); err != nil {
		panic(err)
	}
}
