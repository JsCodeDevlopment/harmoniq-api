package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-" gorm:"column:password_hash"`
	Role     string `json:"role" gorm:"default:'user'"`
	Avatar   string `json:"avatar"`
}
