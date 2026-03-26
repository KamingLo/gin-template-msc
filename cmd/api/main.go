package main

import (
	"log"
	"os"
	"template/config"
	"template/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	config.ConnectDatabase()
	config.InitOAuth()

	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	err := r.Run(":" + port)
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
