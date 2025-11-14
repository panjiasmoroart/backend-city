package seeders

import (
	"backend-city/models"

	"gorm.io/gorm"
)

func SeedPermissions(db *gorm.DB) {
	permissions := []models.Permission{

		{Name: "dashboard-index"},

		{Name: "users-index"},
		{Name: "users-create"},
		{Name: "users-show"},
		{Name: "users-edit"},
		{Name: "users-update"},
		{Name: "users-delete"},

		{Name: "permissions-index"},
		{Name: "permissions-create"},
		{Name: "permissions-show"},
		{Name: "permissions-edit"},
		{Name: "permissions-update"},
		{Name: "permissions-delete"},

		{Name: "roles-index"},
		{Name: "roles-create"},
		{Name: "roles-show"},
		{Name: "roles-edit"},
		{Name: "roles-update"},
		{Name: "roles-delete"},

		{Name: "categories-index"},
		{Name: "categories-create"},
		{Name: "categories-show"},
		{Name: "categories-edit"},
		{Name: "categories-update"},
		{Name: "categories-delete"},

		{Name: "posts-index"},
		{Name: "posts-create"},
		{Name: "posts-show"},
		{Name: "posts-edit"},
		{Name: "posts-update"},
		{Name: "posts-delete"},

		{Name: "products-index"},
		{Name: "products-create"},
		{Name: "products-show"},
		{Name: "products-edit"},
		{Name: "products-update"},
		{Name: "products-delete"},

		{Name: "pages-index"},
		{Name: "pages-create"},
		{Name: "pages-show"},
		{Name: "pages-edit"},
		{Name: "pages-update"},
		{Name: "pages-delete"},

		{Name: "photos-index"},
		{Name: "photos-create"},
		{Name: "photos-delete"},

		{Name: "aparaturs-index"},
		{Name: "aparaturs-create"},
		{Name: "aparaturs-show"},
		{Name: "aparaturs-edit"},
		{Name: "aparaturs-update"},
		{Name: "aparaturs-delete"},

		{Name: "sliders-index"},
		{Name: "sliders-create"},
		{Name: "sliders-delete"},
	}

	for _, p := range permissions {
		db.FirstOrCreate(&p, models.Permission{Name: p.Name})
	}
}
