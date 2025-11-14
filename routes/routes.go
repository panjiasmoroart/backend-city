package routes

import (
	authController "backend-city/controllers/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	// Initialize gin
	router := gin.Default()

	// auth routes (no auth required)
	auth := router.Group("/api")
	{
		auth.POST("/login", authController.Login)
	}

	return router
}
