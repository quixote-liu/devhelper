package models

import "gorm.io/gorm"

// JsonHistory stores operation undo history using JSON Patch (RFC 6902) diffs.
// The first entry per session (IsBase=true) stores the full content;
// subsequent entries store diffs against the previous entry.
// Each session retains at most 50 entries; older ones are pruned automatically.
type JsonHistory struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	SessionID string `gorm:"size:36;index" json:"session_id"`
	SeqNum    int    `gorm:"not null" json:"seq_num"`
	IsBase    bool   `gorm:"default:false" json:"is_base"`
	Content   string `gorm:"type:text;not null" json:"content"`
	Note      string `gorm:"size:255" json:"note"`
}
