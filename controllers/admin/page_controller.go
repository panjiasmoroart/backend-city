package admin

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindPages(c *gin.Context) {

	// Inisialisasi slice
	var pages []models.Page

	// Inisialisasi total
	var total int64

	// Ambil parameter pencarian dan pagination
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	// Query awal dengan preload relasi User
	query := database.DB.Preload("User").Model(&models.Page{})

	// Filter pencarian
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Hitung total data
	query.Count(&total)

	// Ambil data halaman dengan urutan terbaru
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&pages).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Gagal mengambil data halaman",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Mapping manual ke struct response
	pagesResponse := []structs.PagetWithRelationResponse{}
	for _, page := range pages {
		pagesResponse = append(pagesResponse, structs.PagetWithRelationResponse{
			Id:        page.Id,
			Title:     page.Title,
			Slug:      page.Slug,
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
			User: structs.UserSimpleResponse{
				ID:   page.User.Id,
				Name: page.User.Name,
			},
		})
	}

	// Kirim response
	helpers.PaginateResponse(c, pagesResponse, total, page, limit, baseURL, search, "List Data Pages")
}

// CreatePage - Menambahkan data halaman baru
func CreatePage(c *gin.Context) {
	// Inisialisasi struct request
	var req structs.PageCreateRequest

	// Validasi input dari form
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Ambil username dari JWT context
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// Cari user berdasarkan username
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  map[string]string{"user": "User data not found in database"},
		})
		return
	}

	// buat objek page
	page := models.Page{
		Title:   req.Title,
		Slug:    helpers.Slugify(req.Title),
		Content: req.Content,
		UserId:  user.Id,
	}

	// Simpan data ke database
	if err := database.DB.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create page",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim respon sukses
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Page created successfully",
		Data: structs.PageResponse{
			Id:        page.Id,
			Title:     page.Title,
			Slug:      page.Slug,
			Content:   page.Content,
			UserID:    page.UserId,
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// FindPageById - Menampilkan detail halaman berdasarkan ID
func FindPageById(c *gin.Context) {

	// Ambil parameter ID
	id := c.Param("id")

	// Inisialisasi page
	var page models.Page

	// Cari data halaman berdasarkan ID
	if err := database.DB.First(&page, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Page not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim respon sukses dengan data
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Page found",
		Data: structs.PageResponse{
			Id:        page.Id,
			Title:     page.Title,
			Slug:      page.Slug,
			Content:   page.Content,
			UserID:    page.UserId,
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
