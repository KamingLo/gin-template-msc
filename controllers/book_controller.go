package controllers

import (
	"belajar-go/models"
	"belajar-go/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	books, _ := services.GetAllBooks()
	c.JSON(http.StatusOK, books)
}

func CreateBook(c *gin.Context) {
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateBook(&input); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Gagal menyimpan data", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Buku berhasil disimpan", "data": input})
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := services.UpdateBook(id, input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan atau gagal update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update berhasil", "data": book})
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteBook(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan atau sudah dihapus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil dihapus"})
}
