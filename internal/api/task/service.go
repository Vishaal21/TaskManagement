package task

import (
	"encoding/json"
	"fmt"
	"log"
	"task_management/internal/server/ws"
	"task_management/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskService struct {
	taskRepo *TaskRepo
	hub      *ws.Hub
}

func NewTaskService(repo *TaskRepo, hub *ws.Hub) *TaskService {
	return &TaskService{taskRepo: repo, hub: hub}
}

var TaskStatus = map[string]bool{
	"pending":   true,
	"in_progress": true,
	"completed":   true,
}


func (s *TaskService) CreateTask(task *CreateTaskBody, db *gorm.DB, assignedBy int) error {

	// create the task
	taskID, err := s.taskRepo.CreateTask(db, task, assignedBy)
	if err != nil {
		return err
	}

	// fetch the task data
	taskData, err := s.taskRepo.GetTask(db, taskID)
	if err != nil {
		return err
	}

	// Create WebSocket message
	message := map[string]interface{}{
		"event":       "task_created",
		"title":       taskData.Title,
		"description": taskData.Description,
		"assigned_by": taskData.AssignedByUser.Name,
		"assigned_to": taskData.AssignedToUser.Name,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		// Log error but don't fail the task creation
		log.Println("Failed to marshal WebSocket message:", err)
	} else {
		// Broadcast message to all clients
		s.hub.Broadcast <- messageBytes
	}

	return nil

}

func (s *TaskService) GetTasks(db *gorm.DB, assignedBy int) ([]models.Task, error) {
	
	// fetch the tasks with users
	tasks, err := s.taskRepo.GetTasksWithUsers(db, assignedBy)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) UpdateTaskStatus(db *gorm.DB, id uuid.UUID, status string) error {

	if _, exists := TaskStatus[status]; !exists {
		return fmt.Errorf("invalid task status: %s", status)
	}
	return s.taskRepo.UpdateTaskStatus(db, id, status)
}


	
