package seeders

import (
	"backend-city/models"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {

	// Data role dasar dengan nama
	roles := []models.Role{
		{Name: "admin"},
		{Name: "user"},
	}

	// Loop dan assign permissions sesuai role
	for _, role := range roles {
		// Cek/insert role
		db.FirstOrCreate(&role, models.Role{Name: role.Name})

		// Ambil semua permission dari database
		var allPermissions []models.Permission
		db.Find(&allPermissions)

		switch role.Name {
		case "admin":
			// Admin: assign semua permission
			db.Model(&role).Association("Permissions").Replace(allPermissions)
		case "user":
			// User: assign sebagian permission (misalnya hanya index)
			var viewOnly []models.Permission
			db.Where("name IN ?", []string{"posts-index", "photos-index", "sliders-index", "pages-index"}).Find(&viewOnly)
			db.Model(&role).Association("Permissions").Replace(viewOnly)
		}
	}
}
