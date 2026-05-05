package models

import "gorm.io/gorm"

type JsonSchema struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Schema      string `gorm:"type:text;not null" json:"schema"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"`
}
