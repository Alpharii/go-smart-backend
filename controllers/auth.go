package controllers

import (
	"backend-go/config"
	"backend-go/models"
	"backend-go/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


func Register(c *gin.Context) {
	type RegisterInput struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role"`

	}

	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("email = ?", input.Email).First(&models.User{}).Error; err == nil {
		c.JSON(400, gin.H{"error": "Email already exists"})
		return
	}

	if err := config.DB.Where("username = ?", input.Username).First(&models.User{}).Error; err == nil {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}


	user := models.User{
		Email:    input.Email,
		Username: input.Username,
		Password: hashedPassword,
		Roles:    input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	return
}

func Login(c *gin.Context) {
	// code here
}