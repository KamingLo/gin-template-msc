package routes

import (
	"template/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		// Public Routes
		auth.GET("/google", controllers.GoogleLogin)
		auth.GET("/google/callback", controllers.GoogleCallback)
		auth.POST("/login", controllers.Login)
		auth.POST("/otp", controllers.RequestOTP)
		auth.POST("/register", controllers.Register)

		// Private Routes (login needed)
		private := auth.Group("/")
		private.Use(AuthMiddleware())
		{
			private.GET("/me", controllers.GetMe)
			private.POST("/logout", controllers.Logout)
		}
	}
}
