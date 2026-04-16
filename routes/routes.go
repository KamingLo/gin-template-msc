package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// r.Use(CORSMiddleware())
	// r.Use(RateLimitMiddleware())

	AuthRoutes(r)
	BookRoutes(r)

	return r
}
