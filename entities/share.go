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

// SetLink ...
func (sm *ShareMessage) SetLink(shareID string) {
	sm.Link += "?share_id="+shareID
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

// SetCancellationMsg 设置CancellationMsg
func (s *Share) SetCancellationMsg(cm bool) {
	s.attr.CancellationMsg = cm
	return 
}

// NewShare ...
func NewShare(attr ShareAttribute) *Share {
	return &Share{attr: attr}
}

// NewShareID ...
func NewShareID() string {
	return ""
}
