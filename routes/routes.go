package routes

import (
	adminController "backend-city/controllers/admin"
	authController "backend-city/controllers/auth"
	"backend-city/middlewares"

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

	// Protected routes (require authentication)
	protected := router.Group("/api/admin")
	protected.Use(middlewares.AuthMiddleware())
	{
		// Dashboard routes
		protected.GET("/dashboard", middlewares.Permission("dashboard-index"), adminController.Dashboard)
	}

	return router
}
