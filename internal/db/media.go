package db

import (
	"yumi/pkg/ecode"
	"yumi/pkg/external/dbc"
)

//DataMedia ...
type DataMedia int

var media DataMedia

//Media ...
func Media() DataMedia {
	return media
}

//Add ...
func (dm DataMedia) Add(suffix, name, realname, path, operator, operatorid string) (int, error) {
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
func (dm DataMedia) Delete() (int, error) {
	//TODO

	return 0, nil
}

//Update ...
func (dm DataMedia) Update() (int, error) {
	//TODO

	return 0, nil
}

//Search ...
func (dm DataMedia) Search() (int, error) {
	//TODO

	return 0, nil
}

//GetItem ...
func (dm DataMedia) GetItem() (int, error) {
	//TODO

	return 0, nil
}
