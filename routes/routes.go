package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Panggil rute-rute yang sudah dipisah
	AuthRoutes(r)
	BookRoutes(r)

	return r
}
