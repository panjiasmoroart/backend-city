package admin

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"
	"os"
	"path/filepath"

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

// CreatePhoto - Menambahkan data foto baru
func CreatePhoto(c *gin.Context) {

	// Inisialisasi struct
	var req structs.PhotoCreateRequest

	// Validasi data input dari form
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Validasi dan upload file gambar
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  map[string]string{"Image": "Image is required"},
		})
		return
	}

	// Proses upload gambar
	uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
		File:           file,
		AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".gif"},
		MaxSize:        10 << 20, // Maksimal 10MB
		DestinationDir: "public/uploads/photos",
	})

	// Jika upload gagal, kembalikan error
	if uploadResult.Response != nil {
		c.JSON(http.StatusBadRequest, uploadResult.Response)
		return
	}

	// buat objek photo
	photo := models.Photo{
		Image:       uploadResult.FileName,
		Caption:     req.Caption,
		Description: req.Description,
	}

	if err := database.DB.Create(&photo).Error; err != nil {
		// Jika simpan ke DB gagal, hapus gambar yang sudah diupload
		if uploadResult.FileName != "" {
			os.Remove(filepath.Join("public", "uploads", "photos", uploadResult.FileName))
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create photo",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kembalikan response berhasil
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Photo created successfully",
		Data:    photo,
	})
}
