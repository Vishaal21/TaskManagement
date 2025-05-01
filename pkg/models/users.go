package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       int `gorm:"primarykey; type:int;" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"` // Excluded from JSON
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"` // Excluded from JSON
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Excluded from JSON

	Name     string `gorm:"type:text;"`
	Email    string `gorm:"type:text;"`
	Password string `gorm:"type:text" json:"-"`
}
