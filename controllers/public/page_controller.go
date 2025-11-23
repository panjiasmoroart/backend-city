package public

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
			Message: "Failed to fetch pages",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Mapping manual ke struct response
	pagesResponse := []structs.PagetWithRelationResponse{}
	for _, page := range pages {
		pagesResponse = append(pagesResponse, structs.PagetWithRelationResponse{
			Id:    page.Id,
			Title: page.Title,
			Slug:  page.Slug,
			User: structs.UserSimpleResponse{
				ID:   page.User.Id,
				Name: page.User.Name,
			},
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Kirim response
	helpers.PaginateResponse(c, pagesResponse, total, page, limit, baseURL, search, "List Data Pages")
}

// FindPageBySlug - Menampilkan detail halaman berdasarkan Slug
func FindPageBySlug(c *gin.Context) {
	// Ambil parameter slug
	slug := c.Param("slug")

	// Inisialisasi page
	var page models.Page

	// Cari data halaman berdasarkan ID
	if err := database.DB.Preload("User").First(&page, "slug = ?", slug).Error; err != nil {
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
		Data: structs.PagetWithRelationResponse{
			Id:      page.Id,
			Title:   page.Title,
			Slug:    page.Slug,
			Content: page.Content,
			User: structs.UserSimpleResponse{
				ID:   page.User.Id,
				Name: page.User.Name,
			},
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
