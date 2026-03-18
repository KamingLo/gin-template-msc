package routes

import (
	"template/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	{
		// Registrasi Manual & OTP
		authGroup.POST("/otp", controllers.RequestOTP)
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)

		// Google OAuth2
		authGroup.GET("/google", controllers.GoogleLogin)
		authGroup.GET("/google/callback", controllers.GoogleCallback)
	}
}
