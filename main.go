package main

import (
	"backend-city/config"
	"backend-city/database"

	"github.com/gin-gonic/gin"
)

func main() {
	//load config .env
	config.LoadEnv()

	//inisialisasi database
	database.InitDB()

	//inisialiasai Gin
	router := gin.Default()

	//membuat route dengan method GET
	router.GET("/", func(c *gin.Context) {

		//return response JSON
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	//mulai server
	router.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
