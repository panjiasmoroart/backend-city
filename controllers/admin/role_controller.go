package admin

import (
	"backend-city/database"
	"backend-city/helpers"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindRoles(c *gin.Context) {
	var roles []models.Role
	var rolesResponse []structs.RoleResponse
	var total int64

	// Ambil parameter search, page, limit, offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	// Siapkan query
	query := database.DB.Preload("Permissions").Model(&models.Role{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Hitung total data dan ambil data sesuai pagination
	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&roles).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch roles",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Mapping setiap role ke RoleResponse
	for _, role := range roles {
		permissionResponses := []structs.PermissionResponse{} // selalu slice kosong

		for _, permission := range role.Permissions {
			permissionResponses = append(permissionResponses, structs.PermissionResponse{
				Id:        permission.Id,
				Name:      permission.Name,
				CreatedAt: permission.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: permission.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		rolesResponse = append(rolesResponse, structs.RoleResponse{
			Id:          role.Id,
			Name:        role.Name,
			Permissions: permissionResponses,
			CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Kirim response dengan struktur pagination
	helpers.PaginateResponse(c, rolesResponse, total, page, limit, baseURL, search, "List Data Roles")
}

func CreateRole(c *gin.Context) {
	var req structs.RoleCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var permissions []models.Permission
	if len(req.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions)
	}

	role := models.Role{
		Name:        req.Name,
		Permissions: permissions,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create role",
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Role created successfully",
		Data:    role,
	})
}
