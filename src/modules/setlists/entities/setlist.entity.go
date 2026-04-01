package entities

import (
	"time"

	"gorm.io/gorm"
)

type Setlist struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string        `json:"title"`
	UserID    uint          `json:"user_id"`
	PublicID  string        `json:"public_id" gorm:"unique;index"`
	IsPublic  bool          `json:"is_public" gorm:"default:false"`
	Songs     []SetlistItem `json:"songs"`
}

type SetlistItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	SetlistID uint           `json:"setlist_id"`
	Title     string         `json:"title"`
	Artist    string         `json:"artist"`
	URL             string         `json:"url"`
	Key             string         `json:"key"`
	Order           int            `json:"order"`
	ChordVariations string         `json:"chord_variations"`
}
