package public

import (
	"backend-city/database"
	"backend-city/models"
	"backend-city/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Mengambil semua sliders
func FindSliders(c *gin.Context) {

	// Inisialisasi slice untuk menampung data sliders
	var sliders []models.Slider

	// Ambil data sliders dari database
	database.DB.Find(&sliders)

	// Kirimkan response sukses dengan data
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Lists Data Sliders",
		Data:    sliders,
	})
}
