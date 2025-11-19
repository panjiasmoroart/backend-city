package admin

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
			Id:        post.Id,
			Image:     post.Image,
			Title:     post.Title,
			Slug:      post.Slug,
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
			Category: structs.CategorySimpleResponse{
				ID:   post.Category.Id,
				Name: post.Category.Name,
			},
			User: structs.UserSimpleResponse{
				ID:   post.User.Id,
				Name: post.User.Name,
			},
		})
	}

	// Kirim response
	helpers.PaginateResponse(c, postsResponse, total, page, limit, baseURL, search, "List Data Posts")
}

// Menambahkan data post baru
func CreatePost(c *gin.Context) {
	// Inisialisasi struct
	var req structs.PostCreateRequest

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
		DestinationDir: "public/uploads/posts",
	})
	if uploadResult.Response != nil {
		c.JSON(http.StatusBadRequest, uploadResult.Response)
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

	// Cari data user berdasarkan username
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  map[string]string{"user": "User data not found in database"},
		})
		return
	}

	// buat objek post
	post := models.Post{
		Title:      req.Title,
		Slug:       helpers.Slugify(req.Title),
		Content:    req.Content,
		Image:      uploadResult.FileName,
		CategoryId: req.CategoryID,
		UserId:     user.Id,
	}

	// Simpan post ke database
	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create post",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim response sukses
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Post created successfully",
		Data: structs.PostResponse{
			Id:         post.Id,
			Image:      post.Image,
			Title:      post.Title,
			Slug:       post.Slug,
			Content:    post.Content,
			CategoryID: post.CategoryId,
			UserID:     post.UserId,
			CreatedAt:  post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  post.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})

}
