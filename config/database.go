package config

import (
	"fmt"
	"os"
	"template/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load .env
	err := godotenv.Load()

	// Ambil variabel individu
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")

	// Susun DSN (Data Source Name) untuk Postgres
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbName, port, sslMode)

	// Koneksi ke GORM
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi ke Supabase: " + err.Error())
	}

	// Auto Migrate tabel
	database.AutoMigrate(&models.Book{}, &models.User{}, &models.OTP{})

	DB = database
	fmt.Println("Berhasil terkoneksi ke Supabase!")
}
