package data

import (
	"yumi/pkg/ecode"
	"yumi/pkg/stores/dbc"
	"yumi/usecase/media"
)

//DBTable ...
type DBTable int

var dbt DBTable

//DB ...
func DB() DBTable {
	return dbt
}

//Insert ...
func (DBTable) Insert(suffix, name, realname, path, operator, operatorid string) (int, error) {
	var (
		id int64

		err error
	)

	sql := `
		INSERT 
		INTO 
			back_medias 
			("suffix", "name", "real_name", "path", "operator", "operator_id", "operate_time") 
		VALUES 
			(?, ?, ?, ?, ?, ?, sysdate())`
	if id, err = dbc.Get().Insert(sql, suffix, name, realname, path, operator, operatorid); err != nil {
		return 0, ecode.ServerErr(err)
	}

	return int(id), nil
}

// Media is db table
type Media struct {
	ID int64
	media.Media
}


