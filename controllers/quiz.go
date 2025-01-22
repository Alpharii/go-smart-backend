package controllers

import (
	"backend-go/config"
	"backend-go/models"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateQuiz - Handler to create a new quiz
func CreateQuiz(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		CourseID    uint   `json:"course_id" binding:"required"`
	}

	validate := validator.New()

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate input
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

	// Create a new quiz
	quiz := models.Quiz{
		Name:        input.Name,
		Description: input.Description,
		CourseID:    input.CourseID,
	}

	// Save the quiz to the database
	if err := config.DB.Create(&quiz).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create quiz", "details": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Quiz created successfully", "data": quiz})
}

// GetQuizzes - Handler to fetch all quizzes
func GetQuizzes(c *gin.Context) {
	var quizzes []models.Quiz

	// Fetch quizzes from the database
	if err := config.DB.Preload("Course").Find(&quizzes).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch quizzes", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": quizzes})
}

// GetQuizByID - Handler to fetch a quiz by ID
func GetQuizByID(c *gin.Context) {
	id := c.Param("id")
	var quiz models.Quiz

	// Fetch the quiz by ID
	if err := config.DB.Preload("Course").First(&quiz, id).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Quiz not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch quiz", "details": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"data": quiz})
}

func GetQuizzesByCourseID(c *gin.Context) {
	courseID := c.Param("id")

	// Fetch all quizzes belonging to the course
	var quizzes []models.Quiz
	if err := config.DB.Where("course_id = ?", courseID).Find(&quizzes).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "No quizzes found for this course"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch quizzes", "details": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"data": quizzes})
}

// UpdateQuiz - Handler to update a quiz
func UpdateQuiz(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		CourseID    uint   `json:"course_id" binding:"required"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Fetch the quiz by ID
	var quiz models.Quiz
	if err := config.DB.First(&quiz, id).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Quiz not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch quiz", "details": err.Error()})
		}
		return
	}

	// Update quiz fields
	quiz.Name = input.Name
	quiz.Description = input.Description
	quiz.CourseID = input.CourseID

	// Save the updated quiz
	if err := config.DB.Save(&quiz).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update quiz", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Quiz updated successfully", "data": quiz})
}

// DeleteQuiz - Handler to delete a quiz
func DeleteQuiz(c *gin.Context) {
	id := c.Param("id")

	// Delete the quiz
	if err := config.DB.Delete(&models.Quiz{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete quiz", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Quiz deleted successfully"})
}
