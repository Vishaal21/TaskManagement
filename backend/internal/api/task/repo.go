package task

import (
	"task_management/pkg/models"

	"gorm.io/gorm"
)

type TaskRepo struct {
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{}
}

func (t *TaskRepo) CreateTask(db *gorm.DB, task *CreateTaskBody, assignedBy int) error {

	// get the userId from the context
	newTask := models.Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		AssignedTo:  task.AssignedTo,
		AssignedBy:  uint(assignedBy),
	}

	result := db.Create(&newTask)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
