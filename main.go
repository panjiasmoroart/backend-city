package main

import (
	"backend-city/config"
	"backend-city/database"
	"backend-city/database/seeders"
	"backend-city/routes"
)

func main() {
	//load config .env
	config.LoadEnv()

	//inisialisasi database
	database.InitDB()

	//run seeders
	seeders.Seed()

	//setup router
	r := routes.SetupRouter()

	//mulai server
	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
