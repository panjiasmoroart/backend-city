package admin

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats mengambil statistik dashboard
func Dashboard(c *gin.Context) {
	var (
		categoriesCount int64
		postsCount      int64
		productsCount   int64
		aparatursCount  int64
	)

	// Hitung jumlah kategori
	if err := database.DB.Model(&models.Category{}).Count(&categoriesCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to get categories count",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Hitung jumlah postingan
	if err := database.DB.Model(&models.Post{}).Count(&postsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to get posts count",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Hitung jumlah produk
	if err := database.DB.Model(&models.Product{}).Count(&productsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to get products count",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Hitung jumlah aparatur
	if err := database.DB.Model(&models.Aparatur{}).Count(&aparatursCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to get aparaturs count",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Dashboard stats retrieved successfully",
		Data: structs.DashboardResponse{
			CategoriesCount: categoriesCount,
			PostsCount:      postsCount,
			ProductsCount:   productsCount,
			AparatursCount:  aparatursCount,
		},
	})
}
