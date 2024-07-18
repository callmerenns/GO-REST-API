package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string    `gorm:"not null;column:firstname" json:"firstname"`
	LastName  string    `gorm:"not null;column:lastname" json:"lastname"`
	Email     string    `gorm:"not null;unique" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	Role      string    `json:"role"`
	Products  []Product `gorm:"many2many:enrollments;" json:"products"`
}
