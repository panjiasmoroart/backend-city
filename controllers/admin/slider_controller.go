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

// CreateSlider - Menambahkan data slider baru
func CreateSlider(c *gin.Context) {

	// Inisialisasi struct
	var req structs.SliderCreateRequest

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
		DestinationDir: "public/uploads/sliders",
	})

	// Jika upload gagal, kembalikan error
	if uploadResult.Response != nil {
		c.JSON(http.StatusBadRequest, uploadResult.Response)
		return
	}

	// buat object slider
	slider := models.Slider{
		Image:       uploadResult.FileName,
		Description: req.Description,
	}

	if err := database.DB.Create(&slider).Error; err != nil {
		// Jika simpan ke DB gagal, hapus gambar yang sudah diupload
		if uploadResult.FileName != "" {
			os.Remove(filepath.Join("public", "uploads", "sliders", uploadResult.FileName))
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create slider",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kembalikan response berhasil
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Slider created successfully",
		Data:    slider,
	})
}

// DeleteSlider - Menghapus data slider berdasarkan ID
func DeleteSlider(c *gin.Context) {
	// Ambil parameter ID
	id := c.Param("id")

	// Inisialisasi struct
	var slider models.Slider

	// Cari data slider berdasarkan ID
	if err := database.DB.First(&slider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Slider not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Simpan path gambar untuk dihapus
	imagePath := ""
	if slider.Image != "" {
		imagePath = filepath.Join("public", "uploads", "sliders", slider.Image)
	}

	// Hapus data slider dari database
	if err := database.DB.Delete(&slider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete slider",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Hapus file gambar jika ada
	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Slider deleted but failed to remove image",
				Errors:  map[string]string{"image": err.Error()},
			})
			return
		}
	}

	// Kirim response sukses
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Slider deleted successfully",
	})
}
