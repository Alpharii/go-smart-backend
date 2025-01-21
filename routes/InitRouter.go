package routes

import (
	"backend-go/controllers"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/uploads", "./public/uploads")

	//auth
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	//profile
	r.POST("/profile", middleware.IsLogin, controllers.CreateProfile)
	r.GET("/profile", middleware.IsLogin, controllers.GetProfile)
	r.PUT("/profile", middleware.IsLogin, controllers.UpdateProfile)
	r.DELETE("/profile", middleware.IsLogin, controllers.DeleteProfile)

	//course
	r.POST("/course", middleware.IsLogin, middleware.IsAdmin, controllers.CreateCourse)
	r.GET("/courses", middleware.IsLogin, controllers.GetCourses)
	r.GET("/course/:id", middleware.IsLogin, controllers.GetCourseByID)
	r.GET("/course/:id/students", controllers.GetStudentsInCourse)
	r.GET("course/:id/lessons", middleware.IsLogin, controllers.GetLessonsInCourse)
	r.PUT("/course/:id", middleware.IsLogin, middleware.IsAdmin, controllers.UpdateCourse)
	r.DELETE("/course/:id", middleware.IsLogin, middleware.IsAdmin, controllers.DeleteCourse)

	//enrollment
	r.POST("/enroll/:id", middleware.IsLogin, controllers.EnrollCourse)
	r.DELETE("/enroll/:id", middleware.IsLogin, controllers.UnenrollCourse)
	r.GET("/enrollments", middleware.IsLogin, controllers.GetEnrollments)

	//lesson
	r.POST("/lesson", middleware.IsLogin, middleware.IsAdmin, controllers.CreateLesson)
	r.GET("/lessons", middleware.IsEnrolled, controllers.GetLessons)
	r.GET("/lesson/:id", middleware.IsEnrolled, controllers.GetLessonByID)
	r.PUT("/lesson/:id", middleware.IsLogin, middleware.IsAdmin, controllers.UpdateLesson)
	r.DELETE("/lesson/:id", middleware.IsLogin, middleware.IsAdmin, controllers.DeleteLesson)

	//test
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