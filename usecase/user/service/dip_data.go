package service

import (
	"yumi/pkg/sessions"
	"yumi/usecase/user/entity"
)

// Data ...
type Data interface {
	Create(entity.UserAttribute) error
	Update(entity.UserAttribute) error
	Get(ids entity.UserAttributeIDs) (entity.UserAttribute, error)
	Exist(ids entity.UserAttributeIDs) (bool, entity.UserAttribute, error)

	GetSessionsStore() sessions.Store
}
