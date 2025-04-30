package task

import (
	"fmt"
	"log"
	"task_management/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type TaskRepo struct {
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{}
}

func (t *TaskRepo) CreateTask(db *gorm.DB, task *CreateTaskBody, assignedBy int) (uuid.UUID, error) {

	// get the userId from the context
	newTask := models.Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      "pending",
		AssignedTo:  task.AssignedTo,
		AssignedBy:  uint(assignedBy),
	}

	result := db.Create(&newTask)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}

	return newTask.ID, nil
}

func (t *TaskRepo) GetTask(db *gorm.DB, id uuid.UUID) (*models.Task, error) {

	var task models.Task

	result := db.Preload("AssignedToUser").Preload("AssignedByUser").Where("id = ?", id).Find(&task)

	if result.RowsAffected == 0 {
		error := fmt.Errorf("no task found with id %v", id)
		log.Println(error)
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func (t *TaskRepo) GetTasksWithUsers(db *gorm.DB, assignedBy int) ([]models.Task, error) {

	var tasks []models.Task

	result := db.
    // Select specific columns from Task table
    Select("id", "title", "description", "status", "assigned_to", "assigned_by").
    
    // Preload with selected columns for AssignedToUser
    Preload("AssignedToUser", func(db *gorm.DB) *gorm.DB {
        return db.Select("id", "name", "email") // columns from User table
    }).
    
    // Preload with selected columns for AssignedByUser
    Preload("AssignedByUser", func(db *gorm.DB) *gorm.DB {
        return db.Select("id", "name", "email")
    }).
    
    // Your filter condition
    Where("assigned_by = ?", assignedBy).
    
    Find(&tasks)

	if result.RowsAffected == 0 {
		error := fmt.Errorf("no tasks assigned by%v", assignedBy)
		log.Println(error)
		return nil, error
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func (t *TaskRepo) UpdateTaskStatus(db *gorm.DB, id uuid.UUID, status string) error {

	result := db.Model(&models.Task{}).Where("id = ?", id).Update("status", status)

	if result.RowsAffected == 0 {
		log.Println("no task found with id", id)
		return fmt.Errorf("some error occurred while updating task status")
	}

	if result.Error != nil {
		log.Println(result.Error)
		return fmt.Errorf("some error occurred while updating task status")
	}

	return nil
}

	
	
