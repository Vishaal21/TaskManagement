package task

import (
	"task_management/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTaskBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AssignedTo  uint   `json:"assigned_to"`
}

type UpdateTaskStatusBody struct {
	TaskId uuid.UUID `json:"task_id"`
	Status string `json:"status"`
}

type TaskController struct {
	taskService *TaskService
}

func NewTaskController(service *TaskService) *TaskController {
	return &TaskController{taskService: service}
}

func (c *TaskController) CreateTask(ctx *gin.Context, taskService *TaskService) {

	// get the db conn from the context
	dbConn, err := util.GetDbConnectionFromContext(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// get the userId from the context
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(500, gin.H{"error": "userId not found in context"})
		return
	}

	// validate the task details
	var task CreateTaskBody
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// create the task
	err = c.taskService.CreateTask(&task, dbConn, userId.(int))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Task created successfully"})
}

func (c *TaskController) GetTasks(ctx *gin.Context, taskService *TaskService) {

	// get the db conn from the context
	dbConn, err := util.GetDbConnectionFromContext(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// get the userId from the context
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(500, gin.H{"error": "userId not found in context"})
		return
	}

	// get the tasks
	tasks, err := taskService.GetTasks(dbConn, userId.(int))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"tasks": tasks})
}

func (c *TaskController) UpdateTaskStatus(ctx *gin.Context, taskService *TaskService) {

	// get the db conn from the context
	dbConn, err := util.GetDbConnectionFromContext(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// validate the task status
	var taskStatus UpdateTaskStatusBody
	if err := ctx.ShouldBindJSON(&taskStatus); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// update the task status
	err = taskService.UpdateTaskStatus(dbConn, taskStatus.TaskId, taskStatus.Status)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Task status updated successfully"})
}