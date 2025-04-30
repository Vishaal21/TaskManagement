package user

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	userRepo := NewUserRepo()

	userGroup := r.Group("/api/v1/user")
	service := NewUserService(userRepo)
	userController := NewUserController(service) //responsible for making user

	// responsible for making user
	userGroup.POST("/add", func(ctx *gin.Context) {
		userController.AddUser(ctx)
	})

	// responsible for login
	userGroup.GET("/login", func(ctx *gin.Context) {
		userController.Login(ctx)
	})
}
