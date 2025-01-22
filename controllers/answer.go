package controllers

import (
	"backend-go/config"
	"backend-go/models"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateAnswer - Handler to create a new answer
func CreateAnswer(c *gin.Context) {
	var input struct {
		Content     string `json:"content" binding:"required"`
		QuizID     uint   `json:"quiz_id" binding:"required"`
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

	// Check if the question exists
	var quiz models.Quiz
	if err := config.DB.First(&quiz, input.QuizID).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Question not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch question", "details": err.Error()})
		}
		return
	}

	// Create a new answer
	answer := models.Answer{
		Content:     input.Content,
		QuizID:		 input.QuizID,
	}

	// Save the answer to the database
	if err := config.DB.Create(&answer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create answer", "details": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Answer created successfully", "data": answer})
}

// GetAnswersByQuizID - Handler to fetch all answers for a specific question
func GetAnswersByQuizID(c *gin.Context) {
	QuizID := c.Param("question_id")

	// Fetch all answers belonging to the question
	var answers []models.Answer
	if err := config.DB.Where("question_id = ?", QuizID).Find(&answers).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "No answers found for this question"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch answers", "details": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"data": answers})
}

// GetAnswerByID - Handler to fetch a specific answer by ID
func GetAnswerByID(c *gin.Context) {
	answerID := c.Param("id")

	// Fetch the answer by ID
	var answer models.Answer
	if err := config.DB.First(&answer, answerID).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Answer not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch answer", "details": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"data": answer})
}

// UpdateAnswer - Handler to update an existing answer
func UpdateAnswer(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Content     string `json:"content" binding:"required"`
		QuizID 	   uint   `json:"quiz_id" binding:"required"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Fetch the answer by ID
	var answer models.Answer
	if err := config.DB.First(&answer, id).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Answer not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch answer", "details": err.Error()})
		}
		return
	}

	// Update answer fields
	answer.Content = input.Content
	answer.QuizID = input.QuizID

	// Save the updated answer
	if err := config.DB.Save(&answer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update answer", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Answer updated successfully", "data": answer})
}

// DeleteAnswer - Handler to delete an answer
func DeleteAnswer(c *gin.Context) {
	id := c.Param("id")

	// Delete the answer
	if err := config.DB.Delete(&models.Answer{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete answer", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Answer deleted successfully"})
}
