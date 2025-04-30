package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primarykey; default:uuid_generate_v4()"`
	Title       string    `gorm:"type:text;" json:"title"`
	Description string    `gorm:"type:text;" json:"description"`
	Status      string    `json:"status"`
	AssignedTo  uint      `json:"assigned_to"`
	AssignedBy  uint      `json:"assigned_by"`
}
