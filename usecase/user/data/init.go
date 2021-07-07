package data

import (
	"yumi/pkg/stores/mgoc"
	"yumi/usecase/user"
)

// Init ...
func Init(cli *mgoc.Client) {
	user.InitData(New(cli))
}
