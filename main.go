package main

import (
	"backend-go/config"
	"backend-go/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err:= godotenv.Load(); err != nil{
		log.Fatal("Error loading .env file")
	}
	config.ConnectDB()

	r:= gin.Default()
	r = routes.InitRouter()
	r.Run(":8080")
}