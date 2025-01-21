package controllers

import (
	"backend-go/config"
	"backend-go/models"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProfileInput struct {
	FirstName string `form:"first_name" validate:"required"`
	LastName  string `form:"last_name" validate:"required"`
	Phone     string `form:"phone" validate:"required"`
	Address   string `form:"address" validate:"required"`
}


// Create Profile
func CreateProfile(c *gin.Context) {
	var input ProfileInput
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
		uniqueFilename := fmt.Sprintf("%s-%s-%d%s", "profile", input.FirstName, time.Now().Unix(), extension)

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

	// Membuat Profile
	profile := models.Profile{
		UserID:    userID.(uint),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		Address:   input.Address,
		Image:     imageURL,
	}

	// Simpan ke database
	if err := config.DB.Create(&profile).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create profile", "details": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Profile created successfully", "data": profile})
}


// Get Profile by UserID
func GetProfile(c *gin.Context) {
	var profile models.Profile

	// Ambil user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(500, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	// Cari profile berdasarkan UserID
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(404, gin.H{"error": "Profile not found", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": profile})
}

// Update Profile
func UpdateProfile(c *gin.Context) {
	var profile models.Profile
	var input ProfileInput
	var validate = validator.New()

	// Ambil user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(500, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	// Cari profile berdasarkan UserID
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(404, gin.H{"error": "Profile not found", "details": err.Error()})
		return
	}

	// Validasi input JSON atau form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validasi input menggunakan validator
	if err := validate.Struct(&input); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", strings.ToLower(err.Field()), err.Tag()))
		}
		c.JSON(400, gin.H{"error": "Validation failed", "details": errorMessages})
		return
	}

	// Direktori untuk menyimpan file
	publicDir := "./public/uploads"
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		os.MkdirAll(publicDir, os.ModePerm)
	}

	// Proses upload file jika ada file gambar baru
	file, err := c.FormFile("image") // Ambil file dari form-data dengan key "image"
	var imageURL string
	oldImagePath := profile.Image // Simpan URL file lama sebelum diupdate

	if err == nil {
		// Membuat nama file unik
		extension := filepath.Ext(file.Filename)
		uniqueFilename := fmt.Sprintf("%s-%s-%d%s", "profile", profile.FirstName, time.Now().Unix(), extension)

		// Simpan file ke direktori public/uploads
		filePath := filepath.Join(publicDir, uniqueFilename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}

		// Buat URL publik untuk file
		imageURL = fmt.Sprintf("/uploads/%s", uniqueFilename)

		// Hapus file lama jika ada dan tidak kosong
		if oldImagePath != "" {
			oldFilePath := filepath.Join(publicDir, filepath.Base(oldImagePath)) // Path file lama
			if _, err := os.Stat(oldFilePath); err == nil {
				if err := os.Remove(oldFilePath); err != nil {
					fmt.Printf("Failed to delete old file: %s, error: %v\n", oldFilePath, err)
				}
			}
		}
	} else {
		// Jika tidak ada file baru, gunakan URL gambar yang sudah ada
		imageURL = profile.Image
	}

	// Update data profile
	updatedData := models.Profile{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		Address:   input.Address,
		Image:     imageURL, // Perbarui URL gambar
	}

	if err := config.DB.Model(&profile).Updates(updatedData).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update profile", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Profile updated successfully", "data": profile})
}

// Delete Profile
func DeleteProfile(c *gin.Context) {
	var profile models.Profile

	// Ambil user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(500, gin.H{"error": "Failed to retrieve user ID from context"})
		return
	}

	// Cari profile berdasarkan UserID
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(404, gin.H{"error": "Profile not found", "details": err.Error()})
		return
	}

	// Hapus profile dari database
	if err := config.DB.Delete(&profile).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete profile", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Profile deleted successfully"})
}
