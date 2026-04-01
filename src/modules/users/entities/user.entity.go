package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-" gorm:"column:password_hash"`
	Role      string         `json:"role" gorm:"default:'user'"`
	Avatar    string         `json:"avatar"`
	FontSize  string         `json:"font_size" gorm:"default:'medium'"`
	ChordColor string        `json:"chord_color" gorm:"default:'yellow'"`
}
