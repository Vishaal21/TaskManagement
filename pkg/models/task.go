package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type Task struct {
    // Replace gorm.Model with explicit fields
    ID        uuid.UUID      `gorm:"primarykey; type:uuid; default:uuid_generate_v4()" json:"id"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"-"` // Excluded from JSON
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"-"` // Excluded from JSON
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`          // Excluded from JSON
    
    Title       string `gorm:"type:text;" json:"title"`
    Description string `gorm:"type:text;" json:"description"`
    Status      string `json:"status"`
    
    // Foreign key relationships
    AssignedTo  uint `json:"assigned_to"`
    AssignedBy  uint `json:"assigned_by"`
    
    // Associations
    AssignedToUser User `gorm:"foreignKey:AssignedTo" json:"assigned_to_user"`
    AssignedByUser User `gorm:"foreignKey:AssignedBy" json:"assigned_by_user"`
}