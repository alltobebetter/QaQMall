package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"uniqueIndex;size:64"`
	Password string `json:"password"`
	Role     int    `json:"role" gorm:"default:0"`
}
