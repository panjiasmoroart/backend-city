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
