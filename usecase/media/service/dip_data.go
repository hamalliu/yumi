package service

import (
	"yumi/usecase/media/entity"
)
// Data ...
type Data interface {
	Create(ma entity.MediaAttribute) error
	Get(fileNo string) (entity.MediaAttribute, error)
	List(page, line int) ([]entity.MediaAttribute, error)
}
