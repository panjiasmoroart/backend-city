package admin

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindProducts - Ambil semua produk dengan fitur pencarian dan pagination
func FindProducts(c *gin.Context) {
	// Inisialisasi slice
	var products []models.Product

	// Inisialisasi total
	var total int64

	// Ambil parameter pencarian dan pagination
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	// Query awal dengan preload user
	query := database.DB.Preload("User").Model(&models.Product{})

	// Filter pencarian
	if search != "" {
		query = query.Where("title LIKE ? OR owner LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Hitung total data
	query.Count(&total)

	// Ambil data produk dengan urutan terbaru
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch products",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Mapping manual ke struct response
	productsResponse := []structs.ProductWithRelationResponse{}
	for _, product := range products {
		productsResponse = append(productsResponse, structs.ProductWithRelationResponse{
			Id:      product.Id,
			Title:   product.Title,
			Slug:    product.Slug,
			Image:   product.Image,
			Owner:   product.Owner,
			Price:   product.Price,
			Address: product.Address,
			Phone:   product.Phone,
			User: structs.UserSimpleResponse{
				ID:   product.User.Id,
				Name: product.User.Name,
			},
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Kirim response
	helpers.PaginateResponse(c, productsResponse, total, page, limit, baseURL, search, "List Data Products")
}

// CreateProduct - Buat produk baru
func CreateProduct(c *gin.Context) {
	// Inisialisasi struct request
	var req structs.ProductCreateRequest

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
		DestinationDir: "public/uploads/products",
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

	// Buat objek data product baru
	product := models.Product{
		Title:   req.Title,
		Slug:    helpers.Slugify(req.Title),
		Content: req.Content,
		Owner:   req.Owner,
		Address: req.Address,
		Phone:   req.Phone,
		Price:   req.Price,
		Image:   uploadResult.FileName,
		UserId:  user.Id,
	}

	// Simpan post ke database
	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create product",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim response sukses
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Product created successfully",
		Data: structs.ProductResponse{
			Id:        product.Id,
			Title:     product.Title,
			Slug:      product.Slug,
			Content:   product.Content,
			Image:     product.Image,
			Owner:     product.Owner,
			Price:     product.Price,
			Address:   product.Address,
			Phone:     product.Phone,
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// FindProductById - Ambil 1 produk berdasarkan ID
func FindProductById(c *gin.Context) {
	// Ambil parameter ID
	id := c.Param("id")

	// Inisialisasi produk
	var product models.Product

	// Cari produk
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Product not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim data produk
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Product found",
		Data: structs.ProductResponse{
			Id:        product.Id,
			Title:     product.Title,
			Slug:      product.Slug,
			Content:   product.Content,
			Image:     product.Image,
			Owner:     product.Owner,
			Price:     product.Price,
			Address:   product.Address,
			Phone:     product.Phone,
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
