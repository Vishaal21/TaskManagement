package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:text;"`
	Email    string `gorm:"type:text;"`
	Password string `gorm:"type:text;"`
}
