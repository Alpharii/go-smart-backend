package routes

import (
	"backend-go/controllers"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/protected", middleware.IsLogin, middleware.IsAdmin, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello World",
			"data": gin.H{ // Gunakan gin.H untuk nested map
				"user_id": ctx.GetInt("userId"),
				"role":    ctx.GetString("role"),
			},
		})
	})
	

	return r
}