package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

func GetDbConnectionFromContext(ctx *gin.Context) (*gorm.DB, error) {
	dbValue, isExists := ctx.Get("db")
	if !isExists {
		return nil, fmt.Errorf("DB does not exist in context")
	}

	dbConn, ok := dbValue.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("DB connection not found in context")
	}

	return dbConn, nil
}

// Verify password at login
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// This function will create a JWT token
func CreateToken(userId int) (string, error) {

	expiryHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil {
		return "", fmt.Errorf("invalid JWT_EXPIRY_HOURS value: %v", err)
	}
	expiry := time.Now().Add(time.Duration(expiryHours) * time.Hour)

	claims := Claims{
		UserId: int(userId),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Sign the token with secret key
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// This function will verify the JWT token
func VerifyToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err // Signature galat hai ya token invalid hai
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil // Token valid hai aur expire nahi hua
	}
	return nil, fmt.Errorf("invalid token")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Token required"})
			c.Abort()
			return
		}
		// Remove "Bearer " prefix
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := VerifyToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
