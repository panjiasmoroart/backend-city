package admin

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Menampilkan semua kategori dengan fitur pencarian dan pagination
func FindCategories(c *gin.Context) {
	// Inisialisasi slice untuk menampung data
	var categories []models.Category
	var total int64

	// Ambil parameter search, page, limit, offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	// Siapkan query
	query := database.DB.Model(&models.Category{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Hitung total data dan ambil data sesuai pagination
	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&categories).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch categories",
			Errors:  helpers.TranslateErrorMessage(err),
		})
	}

	// Kirim response dengan struktur pagination
	helpers.PaginateResponse(c, categories, total, page, limit, baseURL, search, "List Data Categories")
}
