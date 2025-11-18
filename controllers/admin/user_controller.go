package admin

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindUsers(c *gin.Context) {

	// Inisialisasi slice untuk menampung data
	var users []models.User
	var usersResponse []structs.UserResponse

	// Inisialisasi total
	var total int64

	// Ambil parameter search, page, limit, offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	// Siapkan query
	query := database.DB.Preload("Roles").Model(&models.User{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Hitung total data dan ambil data sesuai pagination
	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Mapping setiap user ke UserResponse
	for _, user := range users {
		// Inisialisasi slice kosong, defaultnya adalah []
		roleResponses := []structs.RoleResponse{}

		// Mapping roles jika ada
		for _, role := range user.Roles {
			roleResponses = append(roleResponses, structs.RoleResponse{
				Id:        role.Id,
				Name:      role.Name,
				CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: role.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		// Append ke list response
		usersResponse = append(usersResponse, structs.UserResponse{
			Id:        user.Id,
			Username:  user.Username,
			Name:      user.Name,
			Email:     user.Email,
			Roles:     roleResponses, // Akan tetap [] kalau tidak ada role
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Kirim response dengan struktur pagination
	helpers.PaginateResponse(c, usersResponse, total, page, limit, baseURL, search, "List Data Users")
}
