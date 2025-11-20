package admin

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindSliders - Menampilkan semua data slider dengan fitur pagination
func FindSliders(c *gin.Context) {

	// Inisialisasi slice
	var sliders []models.Slider

	// Inisialisasi total
	var total int64

	// Ambil parameter pencarian dan pagination dari query string
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.Slider{})
	if search != "" {
		query = query.Where("description LIKE ?", "%"+search+"%")
	}

	// Hitung total data untuk pagination
	query.Count(&total)

	// Ambil data dari database dengan sorting dan batasan pagination
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&sliders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch sliders",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kembalikan response dengan format paginasi
	helpers.PaginateResponse(c, sliders, total, page, limit, baseURL, search, "List Data Sliders")
}
