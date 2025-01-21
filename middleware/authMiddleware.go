package middleware

import (
	"backend-go/config"
	"backend-go/models"
	"backend-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IsLogin(c *gin.Context) {
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

func IsEnrolled(c *gin.Context) {
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

	if role != "admin" && role != "user" {
		c.JSON(400, gin.H{"error": "Invalid role"})
		return
	}	

    courseID := c.Param("id")
    if courseID == "" {
        c.JSON(400, gin.H{"error": "Course ID is required"})
        c.Abort()
        return
    }

    // Periksa apakah user terdaftar di kursus ini
    var enrollment models.Enrollment
    if err := config.DB.Where("user_id = ? AND course_id = ?", userID, courseID).First(&enrollment).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(400, gin.H{"error": "Unauthorized", "message": "You are not enrolled in this course"})
        } else {
            c.JSON(500, gin.H{"error": "Failed to check enrollment", "details": err.Error()})
        }
        c.Abort()
        return
    }

    // Jika user terdaftar, lanjutkan ke handler berikutnya
    c.Next()
}
