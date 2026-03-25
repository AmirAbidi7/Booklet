package models

import "gorm.io/gorm"

type Role = int

const (
	Regular Role = iota
	Admin
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"-" validate:"required"` // json:"-" prevents password from being serialized
	Role     Role   `gorm:"default:0" json:"role" validate:"lt=2"`
	Pdfs     []Pdf  `gorm:"many2many:user_pdf"`
}
