package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(CORSMiddleware()) // Izinkan akses browser
	r.Use(RateLimitMiddleware())

	// Panggil rute-rute yang sudah dipisah
	AuthRoutes(r)
	BookRoutes(r)

	return r
}
