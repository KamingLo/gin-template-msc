package controllers

import (
	"net/http"
	"template/models"
	"template/services"
	"template/utils"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	books, err := services.GetAllBooks()
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get books", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Books data get successfully", books)
}

func CreateBook(c *gin.Context) {
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Input is not valid", err)
		return
	}

	if err := services.CreateBook(&input); err != nil {
		utils.SendError(c, http.StatusConflict, "Failed to save data", err)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Books is saved successfully", input)
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Data format is incorrect", err)
		return
	}

	book, err := services.UpdateBook(id, input)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "Data is not found, and failed to update", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Updated Successfully", book)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteBook(id); err != nil {
		utils.SendError(c, http.StatusNotFound, "Data is not found, or already deleted", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Succesfully deleted", nil)
}
