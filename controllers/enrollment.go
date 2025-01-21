package controllers

import (
	"backend-go/config"
	"backend-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EnrollCourse: Mendaftarkan pengguna ke kursus
func EnrollCourse(c *gin.Context) {
	courseID := c.Param("id")

	// Validasi CourseID
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	// Convert courseID to uint
	courseIDUint, err := strconv.ParseUint(courseID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Course ID format"})
		return
	}

	// Ambil user_id dari context (diset oleh middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	// Periksa apakah kursus ada
	var course models.Course
	if err := config.DB.First(&course, courseIDUint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve course", "details": err.Error()})
		}
		return
	}

	// Periksa apakah pengguna sudah terdaftar
	var existingEnrollment models.Enrollment
	if err := config.DB.Where("user_id = ? AND course_id = ?", userID, courseIDUint).First(&existingEnrollment).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already enrolled in this course"})
		return
	}

	// Buat enrollment baru
	enrollment := models.Enrollment{
		UserID:   userID.(uint),
		CourseID: uint(courseIDUint), // Convert to uint
	}

	// Save the enrollment
	if err := config.DB.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll in course", "details": err.Error()})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"message": "Successfully enrolled in course"})
}

// GetEnrollments: Melihat daftar kursus yang terdaftar
func GetEnrollments(c *gin.Context) {
	// Ambil user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	// Ambil daftar kursus yang terdaftar
	var enrollments []models.Enrollment
	if err := config.DB.Preload("Course").Where("user_id = ?", userID).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollments", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": enrollments})
}

// UnenrollCourse: Membatalkan pendaftaran dari kursus
func UnenrollCourse(c *gin.Context) {
	courseID := c.Param("id")

	// Validasi CourseID
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	// Convert courseID to uint
	courseIDUint, err := strconv.ParseUint(courseID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Course ID format"})
		return
	}

	// Ambil user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	// Periksa apakah pengguna sudah terdaftar
	var enrollment models.Enrollment
	if err := config.DB.Where("user_id = ? AND course_id = ?", userID, courseIDUint).First(&enrollment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollment", "details": err.Error()})
		}
		return
	}

	// Hapus enrollment
	if err := config.DB.Delete(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unenroll", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unenrolled"})
}
