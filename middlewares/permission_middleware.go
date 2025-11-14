package middlewares

import (
	"backend-city/database"
	"backend-city/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware untuk cek apakah user memiliki permission tertentu
func Permission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil username dari context (disimpan oleh middleware Auth)
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		// Load user lengkap dengan relasi roles dan permissions
		err := database.DB.
			Preload("Roles.Permissions").
			Where("username = ?", username).
			First(&user).Error

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Cek apakah user memiliki permission yang diminta
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				if perm.Name == permissionName {
					c.Next() // User memiliki akses, lanjutkan request
					return
				}
			}
		}

		// Jika tidak ditemukan, tolak akses
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden - permission denied"})
		c.Abort()
	}
}
