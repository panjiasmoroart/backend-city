package public

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindPhotos - Menampilkan semua data foto dengan fitur pencarian dan pagination
func FindPhotos(c *gin.Context) {

	// Inisialisasi slice
	var photos []models.Photo

	// Inisialisasi total
	var total int64

	// Ambil parameter pencarian dan pagination dari query string
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.Photo{})
	if search != "" {
		query = query.Where("caption LIKE ?", "%"+search+"%")
	}

	// Hitung total data untuk pagination
	query.Count(&total)

	// Ambil data dari database dengan sorting dan batasan pagination
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&photos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch photos",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kembalikan response dengan format paginasi
	helpers.PaginateResponse(c, photos, total, page, limit, baseURL, search, "List Data Photos")
}

// FindPhotosHome - Ambil 6 photo terbaru
func FindPhotosHome(c *gin.Context) {

	// Inisialisasi slice
	var photos []models.Photo

	// Ambil maksimal 6 photo terbaru
	err := database.DB.Order("id desc").Limit(6).Find(&photos).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch photos",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List Data Photos Home",
		Data:    photos,
	})
}
