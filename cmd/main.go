package main

import (
	"fmt"
	"log"
	"os"
	"task_management/configs"
	"task_management/internal/api/task"
	"task_management/internal/api/user"
	"time"

	// "task_management/internal/ws"
	"task_management/pkg/middleware"
	"task_management/pkg/models"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Initialize() (*gorm.DB, error) {

	// load environment variables
	if err := godotenv.Load(); err != nil {
		
		fmt.Println("No .env file found, using environment variables")
	}

	// initialize database
	dbConn, err := configs.ConnectToDB()
	if err != nil {
		return nil, err
	}

	// Check if migration already done
	if _, err := os.Stat("migrated.txt"); os.IsNotExist(err) {

		dbConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
		// Run migration
		dbConn.AutoMigrate(&models.User{}, &models.Task{})
		// Mark migration as done
		os.Create("migrated.txt")
	}

	return dbConn, nil

}

func main() {

	dbConn, err := Initialize()
	if err != nil {
		log.Println(err)
		return
	}

	r := gin.Default()

	// setup cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Be explicit
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// execute the global middleware
	r.Use(middleware.SetDbConnection(dbConn))

	user.SetupRoutes(r)
	task.SetupRoutes(r)
	// ws.SetupRoutes(r)


    

	r.Run(":8080")
}
