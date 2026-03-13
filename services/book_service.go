package services

import (
	"belajar-go/config"
	"belajar-go/models"

	"gorm.io/gorm"
)

func GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	err := config.DB.Find(&books).Error
	return books, err
}

func CreateBook(book *models.Book) error {
	return config.DB.Create(book).Error
}

func UpdateBook(id string, input map[string]interface{}) (models.Book, error) {
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		return book, err
	}
	err := config.DB.Model(&book).Updates(input).Error
	return book, err
}

func DeleteBook(id string) error {
	result := config.DB.Delete(&models.Book{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
