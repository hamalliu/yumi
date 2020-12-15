package strgmedia

import (
	"yumi/pkg/ecode"
	"yumi/pkg/stores/dbc"
)

//DBTable ...
type DBTable int

var media DBTable

//DB ...
func DB() DBTable {
	return media
}

//Add ...
func (DBTable) Add(suffix, name, realname, path, operator, operatorid string) (int, error) {
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

//Delete ...
func (DBTable) Delete() (int, error) {
	//TODO

	return 0, nil
}

//Update ...
func (DBTable) Update() (int, error) {
	//TODO

	return 0, nil
}

//Search ...
func (DBTable) Search() (int, error) {
	//TODO

	return 0, nil
}

//GetItem ...
func (DBTable) GetItem() (int, error) {
	//TODO

	return 0, nil
}
