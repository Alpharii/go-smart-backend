package main

import (
	"backend-go/config"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err:= godotenv.Load(); err != nil{
		log.Fatal("Error loading .env file")
	}

	r:= gin.Default()
	config.ConnectDB()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}