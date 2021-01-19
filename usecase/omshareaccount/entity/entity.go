package entity

import (
	"yumi/entities"
	"yumi/pkg/codes"
	"yumi/pkg/status"
)

// ShareAccountAttribute ...
type ShareAccountAttribute struct {
	ShareID string `bson:"share_id"`
	// 账号个数
	Total int `bson:"total"`
	// 还可以分享账号个数
	CanShareNumber int `bson:"can_share_number"`
	// 是否领用
	Received bool `bson:"received"`
	// 领用账号
	ReceivedAccount string `bson:"received_account"`

	Parent *ShareAccountAttribute   `bson:"-"`
	Share  *entities.ShareAttribute `bson:"share"`
}

// ShareAccount ...
type ShareAccount struct {
	share *entities.Share
	attr  *ShareAccountAttribute
}

// NewShareAccount ...
func NewShareAccount(attr *ShareAccountAttribute) *ShareAccount {
	return &ShareAccount{share: entities.NewShare(*attr.Share), attr: attr}
}

// LawEnforcement 执法：检查当前数据是否合乎业务规定
func (sa *ShareAccount) LawEnforcement() error {
	// 1. 分享个数不能超过父级可分享数
	if sa.attr.Parent != nil {
		if sa.attr.Parent.CanShareNumber > sa.attr.CanShareNumber {
			return status.FailedPrecondition().WithDetails("分享个数不能超过父级可分享数")
		}
	}

	return nil
}

// SetCancellationMsg 设置撤销msg
// 如果有子分享或已被领取账号不能撤回
func (sa *ShareAccount) SetCancellationMsg() error {
	if sa.share.ChildrenLen() == 0 {
		return status.New(codes.InvalidArgument, "该分享已有子分享不能撤回")
	}

	if sa.attr.Received {
		return status.New(codes.InvalidArgument, "该分享账号已被领取不能撤回")
	}

	sa.share.SetCancellationMsg(true)
	return nil
}

// SetReceived 设置领取账号
func (sa *ShareAccount) SetReceived(acct string) error {
	sa.attr.ReceivedAccount = acct
	sa.attr.Received = true
	return nil
}
