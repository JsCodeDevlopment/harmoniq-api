package entities

import (
	"gorm.io/gorm"
)

type Setlist struct {
	gorm.Model
	Title     string        `json:"title"`
	UserID    uint          `json:"user_id"`
	PublicID  string        `json:"public_id" gorm:"unique;index"`
	IsPublic  bool          `json:"is_public" gorm:"default:false"`
	Songs     []SetlistItem `json:"songs"`
}

type SetlistItem struct {
	gorm.Model
	SetlistID uint   `json:"setlist_id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	URL       string `json:"url"`
	Key       string `json:"key"`
	Order     int    `json:"order"`
}
