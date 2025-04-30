package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetDbConnection(dbConn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", dbConn)
		c.Next()
	}
}
