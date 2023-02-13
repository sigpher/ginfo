package router

import (
	"ginfo/database"
	"ginfo/handler"
	"ginfo/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	database.ConnectDB()

	user := router.Group("/user")
	{
		user.GET("/:id", handler.GetUser)
		user.GET("/", middleware.JWTAuthMiddleware(), handler.GetUsers)
		user.POST("", handler.CreateUser)
	}

	router.POST("/login", handler.Login)

	return router
}
