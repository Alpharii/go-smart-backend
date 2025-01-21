package controllers

import (
	"backend-go/config"
	"backend-go/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CourseInput struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"gte=0"`
	Thumbnail   string  `json:"thumbnail"`
	UserID      uint    `json:"user_id" validate:"required"`
}

// Create Course
func CreateCourse(c *gin.Context) {
	var input CourseInput
	var validate = validator.New()

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validasi input (kecuali UserID karena akan diambil dari header)
	if err := validate.StructExcept(&input, "UserID"); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", strings.ToLower(err.Field()), err.Tag()))
		}
		c.JSON(400, gin.H{"error": "Validation failed", "details": errorMessages})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(500, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	// Membuat Course
	course := models.Course{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Thumbnail:   input.Thumbnail,
		UserID:      userID.(uint),
	}

	// Simpan ke database
	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create course", "details": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Course created successfully", "data": course})
}

// Get All Courses
func GetCourses(c *gin.Context) {
	var courses []models.Course

	// Ambil data dari database
	if err := config.DB.Find(&courses).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch courses", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": courses})
}

// Get Single Course by ID
func GetCourseByID(c *gin.Context) {
	var course models.Course
	id := c.Param("id")

	// Cari course berdasarkan ID
	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Course not found", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": course})
}

// Update Course
func UpdateCourse(c *gin.Context) {
	var course models.Course
	id := c.Param("id")

	// Cari course berdasarkan ID
	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Course not found", "details": err.Error()})
		return
	}

	var input CourseInput
	var validate = validator.New()

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validasi input (kecuali UserID karena akan diambil dari header)
	if err := validate.StructExcept(&input, "UserID"); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", strings.ToLower(err.Field()), err.Tag()))
		}
		c.JSON(400, gin.H{"error": "Validation failed", "details": errorMessages})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(500, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	// Update data course
	updatedData := models.Course{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Thumbnail:   input.Thumbnail,
		UserID:      userID.(uint),
	}

	if err := config.DB.Model(&course).Updates(updatedData).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update course", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Course updated successfully", "data": course})
}

// Delete Course
func DeleteCourse(c *gin.Context) {
	var course models.Course
	id := c.Param("id")

	// Cari course berdasarkan ID
	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Course not found", "details": err.Error()})
		return
	}

	// Hapus course dari database
	if err := config.DB.Delete(&course).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete course", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Course deleted successfully"})
}

func GetStudentsInCourse(c *gin.Context) {
	// Ambil CourseID dari parameter URL
	courseID := c.Param("id")

	// Validasi CourseID
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	// Periksa apakah kursus ada
	var course models.Course
	if err := config.DB.First(&course, courseID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve course", "details": err.Error()})
		}
		return
	}

	// Ambil daftar murid yang terdaftar di kursus ini
	var enrollments []models.Enrollment
	if err := config.DB.Preload("User").Where("course_id = ?", courseID).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollments", "details": err.Error()})
		return
	}

	// Extract students' information
	students := []struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
	}

	for _, enrollment := range enrollments {
		students = append(students, struct {
			UserID   uint   `json:"user_id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			UserID:   enrollment.UserID,
			Username: enrollment.User.Username,
			Email:    enrollment.User.Email,
		})
	}

	c.JSON(http.StatusOK, gin.H{"course": course.Name, "students": students})
}

