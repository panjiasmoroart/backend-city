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
