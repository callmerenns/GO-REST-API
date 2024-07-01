package entity

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Description string  `gorm:"not null" json:"description"`
	Stock       int     `gorm:"not null" json:"stock"`
	Price       float64 `gorm:"not null" json:"price"`
	Users       []User  `gorm:"many2many:enrollments;" json:"users"`
}
