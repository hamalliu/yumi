package data

import (
	"yumi/pkg/stores/mgoc"
	"yumi/usecase/user/service"
)

// MongoCli ...
type MongoCli struct {
	*mgoc.Client
}

var _ service.Data = &MongoCli{}

// New ...
func New(cli *mgoc.Client) *MongoCli {
	return &MongoCli{Client: cli}
}
