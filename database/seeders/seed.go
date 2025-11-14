package seeders

import (
	"backend-city/database"
	"log"
)

func Seed() {
	db := database.DB
	log.Println("Running database seeders...")

	// Jalankan seeder secara berurutan
	SeedPermissions(db)
	SeedRoles(db)
	SeedUsers(db)

	log.Println("Database seeding completed!")
}
