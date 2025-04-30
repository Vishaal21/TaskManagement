package task

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"task_management/internal/server/ws"
)

type TaskService struct {
	taskRepo *TaskRepo
	hub      *ws.Hub
}

func NewTaskService(repo *TaskRepo, hub *ws.Hub) *TaskService {
	return &TaskService{taskRepo: repo, hub: hub}
}

func (s *TaskService) CreateTask(task *CreateTaskBody, db *gorm.DB, assignedBy int) error {

	// create the task
	err := s.taskRepo.CreateTask(db, task, assignedBy)
	if err != nil {
		return err
	}

	// Create WebSocket message
	message := map[string]interface{}{
		"event":       "task_created",
		"title":       task.Title,
		"description": task.Description,
		"assigned_by": assignedBy,
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
