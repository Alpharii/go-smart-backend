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
)

// CreateLesson - Handler to create a new lesson
func CreateLesson(c *gin.Context) {
	var input struct {
		Name        string `form:"name" binding:"required"`
		Description string `form:"description" binding:"required"`
		CourseID    uint   `form:"course_id" binding:"required"`
	}

	var validate = validator.New()

	// Parsing data dari multipart/form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validasi input
	if err := validate.Struct(&input); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", strings.ToLower(err.Field()), err.Tag()))
		}
		c.JSON(400, gin.H{"error": "Validation failed", "details": errorMessages})
		return
	}

	// Check if the course exists
	var course models.Course
	if err := config.DB.First(&course, input.CourseID).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Course not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch course", "details": err.Error()})
		}
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
		uniqueFilename := fmt.Sprintf("%s-%s-%d%s", "lesson", input.Name, time.Now().Unix(), extension)

		// Simpan file ke direktori public/uploads
		filePath := filepath.Join(publicDir, uniqueFilename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}

		imageURL = fmt.Sprintf("/uploads/%s", uniqueFilename)
	}

	// Create a new lesson
	lesson := models.Lesson{
		Name:        input.Name,
		Description: input.Description,
		Image:       imageURL,
		CourseID:    input.CourseID,
	}

	// Save the lesson to the database
	if err := config.DB.Create(&lesson).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create lesson", "details": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Lesson created successfully", "data": lesson})
}


// GetLessons - Handler to fetch all lessons
func GetLessons(c *gin.Context) {
	var lessons []models.Lesson

	// Fetch lessons from the database
	if err := config.DB.Find(&lessons).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch lessons", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": lessons})
}

// GetLessonByID - Handler to fetch a specific lesson by ID
func GetLessonByID(c *gin.Context) {
	lessonID := c.Param("id")

	// Fetch lesson by ID from the database
	var lesson models.Lesson
	if err := config.DB.First(&lesson, lessonID).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Lesson not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch lesson", "details": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"data": lesson})
}

func GetLessonsByCourseID(c *gin.Context) {
	courseID := c.Param("course_id")

	// Fetch all lessons belonging to the course
	var lessons []models.Lesson
	if err := config.DB.Where("course_id = ?", courseID).Find(&lessons).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "No lessons found for this course"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch lessons", "details": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"data": lessons})
}


// UpdateLesson - Handler to update an existing lesson (optional)
func UpdateLesson(c *gin.Context) {
	lessonID := c.Param("id")

	var input struct {
		Name        string `form:"name"`
		Description string `form:"description"`
		CourseID    uint   `form:"course_id"`
	}

	// Parsing data dari multipart/form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Check if the lesson exists
	var lesson models.Lesson
	if err := config.DB.First(&lesson, lessonID).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Lesson not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch lesson", "details": err.Error()})
		}
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
		uniqueFilename := fmt.Sprintf("%s-%s-%d%s", "lesson", input.Name, time.Now().Unix(), extension)

		// Simpan file ke direktori public/uploads
		filePath := filepath.Join(publicDir, uniqueFilename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}

		imageURL = fmt.Sprintf("/uploads/%s", uniqueFilename)
	} else {
		// Gunakan gambar sebelumnya jika tidak ada gambar baru
		imageURL = lesson.Image
	}

	// Update lesson details
	if input.Name != "" {
		lesson.Name = input.Name
	}
	if input.Description != "" {
		lesson.Description = input.Description
	}
	if input.CourseID != 0 {
		lesson.CourseID = input.CourseID
	}
	lesson.Image = imageURL

	// Save the updated lesson to the database
	if err := config.DB.Save(&lesson).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update lesson", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Lesson updated successfully", "data": lesson})
}


// DeleteLesson - Handler to delete a lesson (optional)
func DeleteLesson(c *gin.Context) {
	lessonID := c.Param("id")

	// Check if the lesson exists
	var lesson models.Lesson
	if err := config.DB.First(&lesson, lessonID).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Lesson not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch lesson", "details": err.Error()})
		}
		return
	}

	// Delete the lesson from the database
	if err := config.DB.Delete(&lesson).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete lesson", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Lesson deleted successfully"})
}

// GetLessonsInCourse - Handler to fetch all lessons for a specific course
func GetLessonsInCourse(c *gin.Context) {
	// Get the CourseID from the URL parameter
	courseID := c.Param("id")

	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}


	// Fetch all lessons for the specified course
	var lessons []models.Lesson
	if err := config.DB.Where("course_id = ?", courseID).Find(&lessons).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch lessons for the course", "details": err.Error()})
		return
	}

	// If no lessons are found
	if len(lessons) == 0 {
		c.JSON(404, gin.H{"error": "No lessons found for the specified course"})
		return
	}

	// Return the lessons for the course
	c.JSON(200, gin.H{"data": lessons})
}
