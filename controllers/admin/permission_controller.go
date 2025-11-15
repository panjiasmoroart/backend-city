package admin

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ambil semua permission dengan search dan pagination
func FindPermissions(c *gin.Context) {
	var permissions []models.Permission
	var total int64

	// Ambil parameter search, page, limit, dan offset dari helper
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	// Query awal dari tabel permissions
	query := database.DB.Model(&models.Permission{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Hitung total data
	query.Count(&total)

	// Ambil data sesuai limit dan offset
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&permissions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch permissions",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Response JSON dengan pagination
	helpers.PaginateResponse(c, permissions, total, page, limit, baseURL, search, "List Data Permissions")
}

// Buat permission baru
func CreatePermission(c *gin.Context) {
	var req structs.PermissionCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	permission := models.Permission{
		Name: req.Name,
	}

	if err := database.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create permission",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Permission created successfully",
		Data:    permission,
	})

}
