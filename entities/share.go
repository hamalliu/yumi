package entities

import (
	"yumi/pkg/types"
)

// ShareMessage ...
type ShareMessage struct {
	Title   string `bson:"title"`
	Desc    string `bson:"desc"`
	Link    string `bson:"link"`
	ImgURL  string `bson:"img_url"`
	Type    int    `bson:"type"`
	DataURL string `bson:"data_url"`
}

// SetLink ...
func (sm *ShareMessage) SetLink(shareID string) {
	sm.Link += "?share_id=" + shareID
}

// ShareAttribute ...
type ShareAttribute struct {
	ParentShareID   string          `bson:"parent_share_id"`
	ShareID         string          `bson:"share_idd"`
	Sender          string          `bson:"sender"`
	Receiver        string          `bson:"receiver"`
	Message         ShareMessage    `bson:"message"`
	CancellationMsg bool            `bson:"cancellation_msg"`
	SendTime        types.Timestamp `bson:"send_time"`
	OpenTime        types.Timestamp `bson:"opend_time"`
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
