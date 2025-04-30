package task

import (
	"task_management/internal/util"

	"github.com/gin-gonic/gin"
)

type CreateTaskBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssignedTo  uint   `json:"assigned_to"`
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
