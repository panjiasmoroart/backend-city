package routes

import (
	adminController "backend-city/controllers/admin"
	authController "backend-city/controllers/auth"
	publicController "backend-city/controllers/public"
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

		// Permission routes
		protected.GET("/permissions", middlewares.Permission("permissions-index"), adminController.FindPermissions)
		protected.POST("/permissions", middlewares.Permission("permissions-create"), adminController.CreatePermission)
		protected.GET("/permissions/:id", middlewares.Permission("permissions-show"), adminController.FindPermissionById)
		protected.PUT("/permissions/:id", middlewares.Permission("permissions-update"), adminController.UpdatePermission)
		protected.DELETE("/permissions/:id", middlewares.Permission("permissions-delete"), adminController.DeletePermission)
		protected.GET("/permissions/all", middlewares.Permission("permissions-index"), adminController.FindAllPermissions)

		// Role routes
		protected.GET("/roles", middlewares.Permission("roles-index"), adminController.FindRoles)
		protected.POST("/roles", middlewares.Permission("roles-create"), adminController.CreateRole)
		protected.GET("/roles/:id", middlewares.Permission("roles-show"), adminController.FindRoleById)
		protected.PUT("/roles/:id", middlewares.Permission("roles-update"), adminController.UpdateRole)
		protected.DELETE("/roles/:id", middlewares.Permission("roles-delete"), adminController.DeleteRole)
		protected.GET("/roles/all", middlewares.Permission("roles-index"), adminController.FindAllRoles)

		// User routes
		protected.GET("/users", middlewares.Permission("users-index"), adminController.FindUsers)
		protected.POST("/users", middlewares.Permission("users-create"), adminController.CreateUser)
		protected.GET("/users/:id", middlewares.Permission("users-show"), adminController.FindUserById)
		protected.PUT("/users/:id", middlewares.Permission("users-update"), adminController.UpdateUser)
		protected.DELETE("/users/:id", middlewares.Permission("users-delete"), adminController.DeleteUser)

		// Category routes
		protected.GET("/categories", middlewares.Permission("categories-index"), adminController.FindCategories)
		protected.POST("/categories", middlewares.Permission("categories-create"), adminController.CreateCategory)
		protected.GET("/categories/:id", middlewares.Permission("categories-show"), adminController.FindCategoryById)
		protected.PUT("/categories/:id", middlewares.Permission("categories-update"), adminController.UpdateCategory)
		protected.DELETE("/categories/:id", middlewares.Permission("categories-delete"), adminController.DeleteCategory)
		protected.GET("/categories/all", middlewares.Permission("categories-index"), adminController.FindAllCategories)

		// Post routes
		protected.GET("/posts", middlewares.Permission("posts-index"), adminController.FindPosts)
		protected.POST("/posts", middlewares.Permission("posts-create"), adminController.CreatePost)
		protected.GET("/posts/:id", middlewares.Permission("posts-show"), adminController.FindPostById)
		protected.PUT("/posts/:id", middlewares.Permission("posts-update"), adminController.UpdatePost)
		protected.DELETE("/posts/:id", middlewares.Permission("posts-delete"), adminController.DeletePost)

		// Page routes
		protected.GET("/pages", middlewares.Permission("pages-index"), adminController.FindPages)
		protected.POST("/pages", middlewares.Permission("pages-create"), adminController.CreatePage)
		protected.GET("/pages/:id", middlewares.Permission("pages-show"), adminController.FindPageById)
		protected.PUT("/pages/:id", middlewares.Permission("pages-update"), adminController.UpdatePage)
		protected.DELETE("/pages/:id", middlewares.Permission("pages-delete"), adminController.DeletePage)

		// Product routes
		protected.GET("/products", middlewares.Permission("products-index"), adminController.FindProducts)
		protected.POST("/products", middlewares.Permission("products-create"), adminController.CreateProduct)
		protected.GET("/products/:id", middlewares.Permission("products-show"), adminController.FindProductById)
		protected.PUT("/products/:id", middlewares.Permission("products-update"), adminController.UpdateProduct)
		protected.DELETE("/products/:id", middlewares.Permission("products-delete"), adminController.DeleteProduct)

		// Photo routes
		protected.GET("/photos", middlewares.Permission("photos-index"), adminController.FindPhotos)
		protected.POST("/photos", middlewares.Permission("photos-create"), adminController.CreatePhoto)
		protected.DELETE("/photos/:id", middlewares.Permission("photos-delete"), adminController.DeletePhoto)

		//Slider routes
		protected.GET("/sliders", middlewares.Permission("sliders-index"), adminController.FindSliders)
		protected.POST("/sliders", middlewares.Permission("sliders-create"), adminController.CreateSlider)
		protected.DELETE("/sliders/:id", middlewares.Permission("sliders-delete"), adminController.DeleteSlider)

		// Aparatur routes
		protected.GET("/aparaturs", middlewares.Permission("aparaturs-index"), adminController.FindAparaturs)
		protected.POST("/aparaturs", middlewares.Permission("aparaturs-create"), adminController.CreateAparatur)
		protected.GET("/aparaturs/:id", middlewares.Permission("aparaturs-show"), adminController.FindAparaturById)
		protected.PUT("/aparaturs/:id", middlewares.Permission("aparaturs-update"), adminController.UpdateAparatur)
		protected.DELETE("/aparaturs/:id", middlewares.Permission("aparaturs-delete"), adminController.DeleteAparatur)
	}

	// Public routes
	public := router.Group("/api/public")
	{
		// Post routes
		public.GET("/posts", publicController.FindPosts)
		public.GET("/posts/:slug", publicController.FindPostBySlug)
		public.GET("/posts_home", publicController.FindPostsHome)

		// Page routes
		public.GET("/pages", publicController.FindPages)
		public.GET("/pages/:slug", publicController.FindPageBySlug)

		// Product routes
		public.GET("/products", publicController.FindProducts)
	}

	return router
}
