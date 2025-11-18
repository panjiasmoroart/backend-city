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

func CreateUser(c *gin.Context) {
	// struct user request
	var req = structs.UserCreateRequest{}

	// Bind JSON request ke struct UserCreateRequest + validasi
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Ambil daftar role berdasarkan role_ids (jika dikirim)
	var roles []models.Role
	if len(req.RoleIDs) > 0 {
		database.DB.Where("id IN ?", req.RoleIDs).Find(&roles)
	}

	// Inisialisasi user baru
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
		Roles:    roles,
	}

	// Simpan user ke database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirimkan response sukses (mapping ke UserResponse agar konsisten)
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "User created successfully",
		Data: structs.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func FindUserById(c *gin.Context) {
	// Ambil ID user dari parameter URL
	id := c.Param("id")

	// Inisialisasi user
	var user models.User

	// Cari user berdasarkan ID dan preload relasi Roles
	if err := database.DB.Preload("Roles").First(&user, id).Error; err != nil {
		// Jika user tidak ditemukan, kirim response 404
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Konversi Roles dari model ke struct RoleResponse
	var roleResponses []structs.RoleResponse
	for _, role := range user.Roles {
		roleResponses = append(roleResponses, structs.RoleResponse{
			Id:        role.Id,
			Name:      role.Name,
			CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: role.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Kirim response sukses dengan UserResponse
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User Found",
		Data: structs.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Roles:     roleResponses,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
