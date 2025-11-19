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

// Menambahkan kategori baru
func CreateCategory(c *gin.Context) {
	var req structs.CategoryCreateRequest

	// Validasi input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Buat slug dari nama kategori
	slug := helpers.Slugify(req.Name)

	// Cek apakah slug sudah digunakan
	var existing models.Category
	if err := database.DB.Where("slug = ?", slug).First(&existing).Error; err == nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors: map[string]string{
				"slug": "Slug already exists",
			},
		})
		return
	}

	// Buat objek kategori
	category := models.Category{
		Name: req.Name,
		Slug: slug,
	}

	// Simpan kategori
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create category",
		})
		return
	}

	// Kirim response sukses
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Category created successfully",
		Data:    category,
	})
}

// Mengambil satu kategori berdasarkan ID
func FindCategoryById(c *gin.Context) {
	// Ambil parameter ID
	id := c.Param("id")
	var category models.Category

	// Cari kategori berdasarkan ID
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim data kategori
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Category Found",
		Data:    category,
	})
}
