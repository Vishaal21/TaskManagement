package user

import (
	"task_management/internal/util"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{
		userService: service,
	}
}

type LoginUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddUserBody struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *UserController) AddUser(ctx *gin.Context) {

	var user AddUserBody
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// get the db conn from the context
	dbConn, err := util.GetDbConnectionFromContext(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Insert the user in the database
	err = c.userService.AddUser(&user, dbConn)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User added successfully"})

}

func (c *UserController) Login(ctx *gin.Context) {

	// get the db conn from the context
	dbConn, err := util.GetDbConnectionFromContext(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// validate the user details
	var userData LoginUserBody
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// login user
	token, err := c.userService.LoginUser(userData, dbConn)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// set the token in the response
	ctx.JSON(200, gin.H{"accessToken": token})
}

func (c *UserController) GetUsers(ctx *gin.Context) {

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

	// fetch the users
	users, err := c.userService.GetUsers(dbConn, userId.(int))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// return the users
	ctx.JSON(200, gin.H{"users": users})
}