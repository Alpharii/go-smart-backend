package controllers

import (
	"backend-go/config"
	"backend-go/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CourseInput struct {
	Name        string  `form:"name" validate:"required"`
	Description string  `form:"description" validate:"required"`
	Price       float64 `form:"price" validate:"gte=0"`
	Image       string  `form:"image"`
}

// Create Course
func CreateCourse(c *gin.Context) {
	var input CourseInput
	var validate = validator.New()

	// Parsing data dari multipart/form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validasi input (kecuali UserID karena akan diambil dari context)
	if err := validate.Struct(&input); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", strings.ToLower(err.Field()), err.Tag()))
		}
		c.JSON(400, gin.H{"error": "Validation failed", "details": errorMessages})
		return
	}

	// Ambil user_id dari context (diset oleh middleware IsLogin)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(500, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	// Direktori untuk menyimpan file
	publicDir := "./public/uploads"
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		os.MkdirAll(publicDir, os.ModePerm)
	}

	// Proses upload file gambar
	file, err := c.FormFile("image") // Ambil file dari form-data dengan key "image"
	var imageURL string
	if err == nil {
		// Membuat nama file unik
		extension := filepath.Ext(file.Filename)
		uniqueFilename := fmt.Sprintf("%s-%s-%d%s", "course", input.Name, time.Now().Unix(), extension)

		// Simpan file ke direktori public/uploads
		filePath := filepath.Join(publicDir, uniqueFilename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}

		imageURL = fmt.Sprintf("/uploads/%s", uniqueFilename)
	} else {
		imageURL = ""
	}

	// Membuat Course
	course := models.Course{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Image:       imageURL,
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

	// Parsing data dari multipart/form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validasi input
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

	// Direktori untuk menyimpan file
	publicDir := "./public/uploads"
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		os.MkdirAll(publicDir, os.ModePerm)
	}

	// Proses upload file gambar (jika ada)
	file, err := c.FormFile("image")
	var imageURL string
	if err == nil {
		// Membuat nama file unik
		extension := filepath.Ext(file.Filename)
		uniqueFilename := fmt.Sprintf("%s-%s-%d%s", "course", input.Name, time.Now().Unix(), extension)

		// Simpan file ke direktori public/uploads
		filePath := filepath.Join(publicDir, uniqueFilename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}

		imageURL = fmt.Sprintf("/uploads/%s", uniqueFilename)
	} else {
		// Gunakan gambar sebelumnya jika tidak ada gambar baru
		imageURL = course.Image
	}

	// Update data course
	updatedData := models.Course{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Image:       imageURL,
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

