package public

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindAparaturs - Menampilkan semua data aparatur dengan fitur pencarian dan paginasi
func FindAparaturs(c *gin.Context) {
	// Inisialisasi slice
	var aparaturs []models.Aparatur

	// Inisialisasi total
	var total int64

	// Ambil parameter pencarian dan paginasi dari request
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.Aparatur{})
	if search != "" {
		// Filter berdasarkan nama atau jabatan
		query = query.Where("name LIKE ? OR position LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&aparaturs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch aparaturs",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim response dalam format paginasi
	helpers.PaginateResponse(c, aparaturs, total, page, limit, baseURL, search, "List Data Aparaturs")
}
