package entities

import (
	"yumi/pkg/types"
)

// ShareMessage ...
type ShareMessage struct {
	Title   string
	Desc    string
	Link    string
	ImgURL  string
	Type    int
	DataURL string
}

// ShareAttribute ...
type ShareAttribute struct {
	ParentShareID   string
	ShareID         string
	Sender          string
	Receiver        string
	Message         ShareMessage
	CancellationMsg bool
	SendTime        types.Timestamp
	OpenTime        types.Timestamp
}

// Share ...
type Share struct {
	attr     ShareAttribute
	children []Share
}

// ChildrenLen 子分享个数
func (s *Share) ChildrenLen() int {
	return len(s.children)
}

// DataShare ...
type DataShare interface {
	CreateShare(sa ShareAttribute) error
	UpdateShare(sa ShareAttribute) error
	GetShare(shareID string) (ShareAttribute, error)
	GetSubShare(parentShareID string) ([]ShareAttribute, error)
}

// NewShare ...
func NewShare(attr ShareAttribute) *Share {
	return &Share{attr: attr}
}

// NewShareID ...
func NewShareID() string {
	return ""
}
