package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int     `gorm:"not null;default:0" json:"stock"`
	Picture     string  `gorm:"size:255" json:"picture"`
	Categories  string  `gorm:"size:255" json:"categories"`
}

type Category struct {
	gorm.Model
	Name        string `gorm:"size:255;not null;unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}
