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

// UpdateAparatur - Perbarui data aparatur berdasarkan ID
func UpdateAparatur(c *gin.Context) {

	// Ambil parameter ID
	id := c.Param("id")

	// Inisialisasi aparatur
	var aparatur models.Aparatur

	// Cari data yang akan diperbarui
	if err := database.DB.First(&aparatur, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Aparatur not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.AparaturUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Simpan path gambar lama jika ada
	oldImagePath := ""
	if aparatur.Image != "" {
		oldImagePath = filepath.Join("public", "uploads", "aparaturs", aparatur.Image)
	}

	// Cek apakah ada gambar baru yang diupload
	file, err := c.FormFile("image")
	if err == nil {
		uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
			File:           file,
			AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".gif"},
			MaxSize:        10 << 20,
			DestinationDir: "public/uploads/aparaturs",
		})

		if uploadResult.Response != nil {
			c.JSON(http.StatusBadRequest, uploadResult.Response)
			return
		}

		// Set gambar baru
		aparatur.Image = uploadResult.FileName
	}

	// Perbarui data lainnya
	aparatur.Name = req.Name
	aparatur.Position = req.Position
	aparatur.Description = req.Description

	// Simpan ke database
	if err := database.DB.Save(&aparatur).Error; err != nil {
		// Hapus gambar baru jika penyimpanan gagal
		if file != nil && aparatur.Image != "" {
			os.Remove(filepath.Join("public", "uploads", "aparaturs", aparatur.Image))
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update aparatur",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Jika ada file baru dan file lama masih ada, hapus gambar lama
	if file != nil && oldImagePath != "" {
		_ = os.Remove(oldImagePath)
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Aparatur updated successfully",
		Data:    aparatur,
	})
}

// DeleteAparatur - Hapus data aparatur
func DeleteAparatur(c *gin.Context) {

	// Ambil parameter ID
	id := c.Param("id")

	// Inisialisasi aparatur
	var aparatur models.Aparatur

	// Cari aparatur berdasarkan ID
	if err := database.DB.First(&aparatur, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Aparatur not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Simpan path gambar untuk dihapus nanti
	imagePath := ""
	if aparatur.Image != "" {
		imagePath = filepath.Join("public", "uploads", "aparaturs", aparatur.Image)
	}

	// Hapus data dari database
	if err := database.DB.Delete(&aparatur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete aparatur",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Hapus file gambar jika ada
	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Aparatur deleted but failed to remove image",
				Errors:  map[string]string{"image": err.Error()},
			})
			return
		}
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Aparatur deleted successfully",
	})
}
