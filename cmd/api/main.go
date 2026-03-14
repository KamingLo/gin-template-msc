package main

import (
	"os"
	"template/config"
	"template/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Env (Penting dilakukan di awal agar Gin tahu modenya)
	godotenv.Load()

	// 2. Set Gin Mode berdasarkan .env
	// Jika di .env isinya 'release', maka Gin akan masuk ke mode produksi
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Inisialisasi Database
	config.ConnectDatabase()

	// 4. Setup Router
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000" // Default port jika di env kosong
	}
	r.Run(":" + port)
}
