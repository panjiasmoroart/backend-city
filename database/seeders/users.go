package seeders

import (
	"backend-city/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	// Hash password default
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	// Dapatkan role
	var adminRole, userRole models.Role
	db.Where("name = ?", "admin").First(&adminRole)
	db.Where("name = ?", "user").First(&userRole)

	// Data user awal
	users := []models.User{
		{
			Name:     "Admin",
			Username: "admin",
			Email:    "admin@gmail.com",
			Password: string(password),
			Roles:    []models.Role{adminRole},
		},
		{
			Name:     "User",
			Username: "user",
			Email:    "user@gmail.com",
			Password: string(password),
			Roles:    []models.Role{userRole},
		},
	}

	for _, u := range users {
		var user models.User
		if err := db.Where("username = ?", u.Username).First(&user).Error; err != nil {
			db.Create(&u)
		} else {
			db.Model(&user).Updates(models.User{
				Email:    u.Email,
				Password: string(password),
			})
			db.Model(&user).Association("Roles").Replace(u.Roles)
		}
	}
}
