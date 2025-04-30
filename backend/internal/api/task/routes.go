package task

import (
	"github.com/gin-gonic/gin"
	"task_management/internal/util"
	"task_management/internal/server/ws"
	// "task_management/internal/api/user"
)

func SetupRoutes(r *gin.Engine) {

	// Initialize WebSocket hub that will handle WebSocket connections
	hub := ws.NewHub()
	go hub.Run()

	// userRepo := user.NewUserRepo()
	taskRepo := NewTaskRepo()
	taskService := NewTaskService(taskRepo, hub)
	taskController := NewTaskController(taskService)

	r.GET("/ws", ws.ServeWs(hub))

	taskGroup := r.Group("/api/v1/task", util.AuthMiddleware())

	// responsible for creating task
	taskGroup.POST("/create", func(ctx *gin.Context) {
		taskController.CreateTask(ctx, taskService)
	})

	// responsible for getting tasks
	taskGroup.GET("/get-tasks", func(ctx *gin.Context) {
		taskController.GetTasks(ctx, taskService)
	})

	// responsible for updating task status
	taskGroup.POST("/update-status", func(ctx *gin.Context) {
		taskController.UpdateTaskStatus(ctx, taskService)
	})

}
