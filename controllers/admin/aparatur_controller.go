package admin

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

	// Query awal
	query := database.DB.Model(&models.Aparatur{})

	// Filter pencarian berdasarkan nama atau jabatan
	if search != "" {
		query = query.Where("name LIKE ? OR position LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Hitung total data
	query.Count(&total)

	// Ambil data dengan urutan terbaru
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

// CreateAparatur - Membuat data aparatur baru
func CreateAparatur(c *gin.Context) {

	// Inisialisasi struct request
	var req structs.AparaturCreateRequest

	// Validasi input dari form multipart
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Ambil file gambar dari form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  map[string]string{"Image": "Image is required"},
		})
		return
	}

	// Upload file menggunakan helper
	uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
		File:           file,
		AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".gif"},
		MaxSize:        10 << 20, // Maksimal 10MB
		DestinationDir: "public/uploads/aparaturs",
	})
	if uploadResult.Response != nil {
		c.JSON(http.StatusBadRequest, uploadResult.Response)
		return
	}

	// Buat objek aparatur
	aparatur := models.Aparatur{
		Name:        req.Name,
		Position:    req.Position,
		Description: req.Description,
		Image:       uploadResult.FileName,
	}

	// Simpan ke database
	if err := database.DB.Create(&aparatur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create aparatur",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim response sukses
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Aparatur created successfully",
		Data:    aparatur,
	})
}

// FindAparaturById - Ambil data aparatur berdasarkan ID
func FindAparaturById(c *gin.Context) {

	// Ambil parameter ID
	id := c.Param("id")

	// Inisialisasi aparatur
	var aparatur models.Aparatur

	// Cari berdasarkan ID
	if err := database.DB.First(&aparatur, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Aparatur not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim response sukses
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Aparatur found",
		Data:    aparatur,
	})
}
