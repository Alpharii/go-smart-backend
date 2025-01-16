package middleware

import (
	"backend-go/utils"

	"github.com/gin-gonic/gin"
)

func IsLogin(c*gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(400, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	token := authHeader[len("Bearer "):]
	userID, role, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(400, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Set("user_id", userID)
	c.Set("role", role)
	c.Next()
}

func IsAdmin(c*gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(400, gin.H{"error": "Unauthorized", "message": "You are not an admin"})
		c.Abort()
		return
	}
	c.Next()
}