package services

import (
	"errors"
	"template/config"
	"template/models"
	"template/utils"
)

func RegisterUser(user *models.User) error {
	var existingUser models.User
	if config.DB.Where("email = ?", user.Email).First(&existingUser).RowsAffected > 0 {
		return errors.New("email ini sudah terpakai")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return config.DB.Create(user).Error
}

func LoginUser(input models.UserLogin) (string, error) {
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return "", errors.New("data tidak ditemukan")
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return "", errors.New("password yang kamu masukkan salah")
	}

	return utils.GenerateToken(user.ID, user.Email)
}
