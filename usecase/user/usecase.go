package user

import (
	"yumi/pkg/stores/mgoc"
	"yumi/usecase/user/data"
	"yumi/usecase/user/service"
)

func Usecase(mongoC *mgoc.Client) (*service.Service, error) {
	data := data.New(mongoC)
	return service.New(data)
}
