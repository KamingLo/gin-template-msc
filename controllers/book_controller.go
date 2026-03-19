package controllers

import (
	"net/http"
	"template/models"
	"template/services"
	"template/utils" // Pastikan import helper ini ada

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	books, err := services.GetAllBooks()
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Gagal mengambil data buku", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Berhasil mengambil semua data buku", books)
}

func CreateBook(c *gin.Context) {
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Input tidak valid", err)
		return
	}

	if err := services.CreateBook(&input); err != nil {
		utils.SendError(c, http.StatusConflict, "Gagal menyimpan data", err)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Buku berhasil disimpan", input)
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Format data salah", err)
		return
	}

	book, err := services.UpdateBook(id, input)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "Data tidak ditemukan atau gagal update", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Update berhasil", book)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteBook(id); err != nil {
		utils.SendError(c, http.StatusNotFound, "Data tidak ditemukan atau sudah dihapus", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Berhasil dihapus", nil)
}
