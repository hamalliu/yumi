package entity

import (
	"yumi/entities"
	"yumi/pkg/codes"
	"yumi/pkg/status"
)

// ShareAccountAttribute ...
type ShareAccountAttribute struct {
	ShareID string
	// 账号个数
	Total int
	// 还可以分享账号个数
	CanShareNumber int
	// 是否领用
	Received bool

	Parent *ShareAccountAttribute
	Share  *entities.ShareAttribute
}

// ShareAccount ...
type ShareAccount struct {
	share *entities.Share
	attr  ShareAccountAttribute
}

// NewShareAccount ...
func NewShareAccount(attr ShareAccountAttribute) *ShareAccount {
	return &ShareAccount{share: entities.NewShare(*attr.Share), attr: attr}
}

// Create ...
func (sa *ShareAccount) Create() error {
	if sa.attr.Total < 1 || sa.attr.CanShareNumber < 1 {
		return status.New(codes.FailedPrecondition, codes.FailedPrecondition.String())
	}

	if sa.attr.Parent != nil {
		if sa.attr.Parent.CanShareNumber > sa.attr.CanShareNumber {
			return status.New(codes.FailedPrecondition, codes.FailedPrecondition.String())
		}
	}

	return nil
}

// CloseSubShare 取消子分享
func (sa *ShareAccount) CloseSubShare(shareID string) error {
	return nil
}
