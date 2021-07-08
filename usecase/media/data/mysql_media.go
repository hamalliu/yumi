package data

import (
	"yumi/usecase/media/entity"
)

//Create ...
func (db *MysqlDB) Create(ma entity.MediaAttribute) error {
	return nil
}

//Get ...
func (db *MysqlDB) Get(fileNo string) (ma entity.MediaAttribute, err error) {
	return
}

//List ...
func (db *MysqlDB) List(page, line int) (mas []entity.MediaAttribute, err error) {
	return
}
