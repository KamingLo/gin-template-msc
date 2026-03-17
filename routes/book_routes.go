package routes

import (
	"template/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoutes(r *gin.Engine) {
	bookGroup := r.Group("/books")
	{
		// Public: Lihat semua buku
		bookGroup.GET("", controllers.GetBooks)

		// Protected: Operasi tulis/hapus
		protected := bookGroup.Group("/")
		protected.Use(AuthMiddleware())
		{
			protected.POST("/", controllers.CreateBook)
			protected.PATCH("/:id", controllers.UpdateBook)
			protected.DELETE("/:id", controllers.DeleteBook)
		}
	}
}
