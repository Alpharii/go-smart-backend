package config

import (
	"backend-go/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var DB*gorm.DB

func ConnectDB(){
	dsn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Connected to database")

	// DeleteMigration()

	DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Course{},
		&models.Enrollment{},
		&models.Lesson{},
		&models.Quiz{},
		&models.Question{},
		&models.Answer{},
		&models.UserQuiz{},
		&models.UserAnswer{},
	)
	fmt.Println("Database Migrated")
}

func DeleteMigration(){
	DB.Migrator().DropTable(
		&models.User{},
		&models.Profile{},
		&models.Course{},
		&models.Enrollment{},
		&models.Lesson{},
		&models.Quiz{},
		&models.Question{},
		&models.Answer{},
		&models.UserQuiz{},
		&models.UserAnswer{},
	)
	fmt.Println("Table deleted")
}