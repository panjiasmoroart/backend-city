package public

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Menampilkan semua post dengan fitur pencarian dan pagination
func FindPosts(c *gin.Context) {

	// Inisialisasi slice
	var posts []models.Post

	// Inisialisasi total
	var total int64

	// Ambil parameter pencarian dan pagination
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	// Query awal dengan preload
	query := database.DB.Preload("Category").Preload("User").Model(&models.Post{})

	// Filter pencarian
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Hitung total
	query.Count(&total)

	// Ambil data post dengan urutan terbaru
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&posts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch posts",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Mapping manual ke struct response
	postsResponse := []structs.PostWithRelationResponse{}
	for _, post := range posts {
		postsResponse = append(postsResponse, structs.PostWithRelationResponse{
			Id:      post.Id,
			Image:   post.Image,
			Title:   post.Title,
			Slug:    post.Slug,
			Content: post.Content,
			Category: structs.CategorySimpleResponse{
				ID:   post.Category.Id,
				Name: post.Category.Name,
			},
			User: structs.UserSimpleResponse{
				ID:   post.User.Id,
				Name: post.User.Name,
			},
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Kirim response
	helpers.PaginateResponse(c, postsResponse, total, page, limit, baseURL, search, "List Data Posts")
}

// Mengambil detail post berdasarkan slug
func FindPostBySlug(c *gin.Context) {

	// Ambil parameter slug
	slug := c.Param("slug")

	// Inisialisasi post
	var post models.Post

	// Cari post dan preload relasi
	if err := database.DB.Preload("Category").Preload("User").Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Post not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim data post
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Post found",
		Data: structs.PostWithRelationResponse{
			Id:      post.Id,
			Image:   post.Image,
			Title:   post.Title,
			Slug:    post.Slug,
			Content: post.Content,
			Category: structs.CategorySimpleResponse{
				ID:   post.Category.Id,
				Name: post.Category.Name,
			},
			User: structs.UserSimpleResponse{
				ID:   post.User.Id,
				Name: post.User.Name,
			},
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// Menampilkan 6 post terbaru
func FindPostsHome(c *gin.Context) {

	// Inisialisasi slice
	var posts []models.Post

	// Ambil maksimal 5 post terbaru dengan preload Category dan User
	err := database.DB.Preload("Category").Preload("User").
		Order("id desc").Limit(6).Find(&posts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch posts",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Mapping ke struct response
	postsResponse := []structs.PostWithRelationResponse{}
	for _, post := range posts {
		postsResponse = append(postsResponse, structs.PostWithRelationResponse{
			Id:      post.Id,
			Image:   post.Image,
			Title:   post.Title,
			Slug:    post.Slug,
			Content: post.Content,
			Category: structs.CategorySimpleResponse{
				ID:   post.Category.Id,
				Name: post.Category.Name,
			},
			User: structs.UserSimpleResponse{
				ID:   post.User.Id,
				Name: post.User.Name,
			},
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Kirim response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List Data Posts Home",
		Data:    postsResponse,
	})
}
