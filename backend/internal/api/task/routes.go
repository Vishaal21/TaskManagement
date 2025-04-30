package task

import (
	"github.com/gin-gonic/gin"
	"task_management/internal/util"
	"task_management/internal/server/ws"
)

func SetupRoutes(r *gin.Engine) {

	// Initialize WebSocket hub that will handle WebSocket connections
	hub := ws.NewHub()
	go hub.Run()

	taskRepo := NewTaskRepo()
	taskService := NewTaskService(taskRepo, hub)
	taskController := NewTaskController(taskService)

	r.GET("/ws", ws.ServeWs(hub))

	taskGroup := r.Group("/api/v1/task", util.AuthMiddleware())

	// responsible for creating task
	taskGroup.POST("/create", func(ctx *gin.Context) {
		taskController.CreateTask(ctx, taskService)
	})

}
