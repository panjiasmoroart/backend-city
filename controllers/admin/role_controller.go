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

func FindRoleById(c *gin.Context) {
	// Ambil parameter ID dari URL
	id := c.Param("id")

	// Inisialisasi variable
	var role models.Role

	// Ambil data role beserta relasi permissions
	if err := database.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role not found",
		})
		return
	}

	// Mapping relasi permissions ke bentuk response
	permissionResponses := []structs.PermissionResponse{}
	for _, permission := range role.Permissions {
		permissionResponses = append(permissionResponses, structs.PermissionResponse{
			Id:        permission.Id,
			Name:      permission.Name,
			CreatedAt: permission.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: permission.UpdatedAt.Format("2006-01-02 15:04:05"),
		})

	}

	// Siapkan response akhir
	roleResponse := structs.RoleResponse{
		Id:          role.Id,
		Name:        role.Name,
		Permissions: permissionResponses,
		CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Kirim response JSON
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Role found",
		Data:    roleResponse,
	})
}

func UpdateRole(c *gin.Context) {
	// Ambil ID dari parameter
	id := c.Param("id")

	// Inisialisasi variable
	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role not found",
		})
		return
	}

	var req structs.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	role.Name = req.Name

	var permissions []models.Permission
	if len(req.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions)
	}
	database.DB.Model(&role).Association("Permissions").Replace(&permissions)

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update role",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Role updated successfully",
		Data:    role,
	})
}

// Hapus role
func DeleteRole(c *gin.Context) {
	// Ambil ID dari parameter
	id := c.Param("id")

	// Inisialisasi variable
	var role models.Role

	// Ambil data role
	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Hapus seluruh relasi role<->permission di pivot table
	if err := database.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to detach role from permissions",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Hapus role
	if err := database.DB.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete role",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirim response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Role deleted successfully",
	})

}

// Mengambil semua roles
func FindAllRoles(c *gin.Context) {

	// Inisialisasi slice untuk menampung data roles
	var roles []models.Role

	// Ambil data roles dari database
	database.DB.Find(&roles)

	// Kirimkan response sukses dengan data user
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Lists Data Roles",
		Data:    roles,
	})
}
