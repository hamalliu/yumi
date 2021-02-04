package data

import (
	"yumi/pkg/ecode"
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
func (db *MysqlDB) Create(suffix, name, realname, path, operator, operatorid string) (int, error) {
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
	if id, err = db.Insert(sql, suffix, name, realname, path, operator, operatorid); err != nil {
		return 0, ecode.ServerErr(err)
	}

	return int(id), nil
}

// Media is db table
type Media struct {
	ID int64
	entity.MediaAttribute
}
