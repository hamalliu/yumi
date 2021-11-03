package media

import (
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/usecase/media/data"
	"yumi/usecase/media/service"
)

// Usecase ...
func Usecase(mysqlC *mysqlx.Client) (*service.Service, error) {
	data := data.New(mysqlC)
	return service.New(data)
}
