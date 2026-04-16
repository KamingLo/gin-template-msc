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
	err := godotenv.Load()

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbName, port, sslMode)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database: " + err.Error())
	}

	err = database.AutoMigrate(models.ModelsRegistry...)
	fmt.Println(models.ModelsRegistry...)
	if err != nil {
		fmt.Print("Database Failed to migrate" + err.Error())
	}

	DB = database
	fmt.Println("Database Succesfully Connected")
}
