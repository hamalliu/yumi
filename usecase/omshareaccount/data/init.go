package data

import (
	"yumi/pkg/stores/mgoc"
	"yumi/usecase/omshareaccount"
)

// Init ...
func Init() {
	omshareaccount.InitData(New(mgoc.Get().Database("")))
}
