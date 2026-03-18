package services

import (
	"errors"
	"fmt"
	"math/rand"
	"template/config"
	"template/models"
	"template/utils"
	"time"

	"gorm.io/gorm"
)

func RequestOTP(email string) error {
	// 1. Cek apakah email sudah terdaftar di tabel users
	var existingUser models.User
	err := config.DB.Where("email = ?", email).First(&existingUser).Error

	// Jika tidak ada error (berarti user ditemukan), maka kembalikan error
	if err == nil {
		return errors.New("email sudah terdaftar, silakan gunakan email lain")
	}

	// 2. Generate 6 digit angka
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	otp := models.OTP{
		Email:     email,
		Code:      code,
		ExpiredAt: time.Now().Add(5 * time.Minute), // Valid 5 menit
	}

	// 3. Hapus OTP lama untuk email yang sama agar tidak menumpuk
	config.DB.Where("email = ?", email).Delete(&models.OTP{})

	// 4. Simpan OTP baru
	if err := config.DB.Create(&otp).Error; err != nil {
		return err
	}

	// Simulasi kirim email (Log ke terminal)
	fmt.Printf("OTP untuk %s adalah: %s\n", email, code)
	return nil
}

func RegisterWithOTP(user *models.User, otpCode string) error {
	var otp models.OTP
	// Cek kecocokan email dan kode
	err := config.DB.Where("email = ? AND code = ?", user.Email, otpCode).First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("kode OTP yang kamu masukkan salah")
		}
		return err
	}

	// Cek masa berlaku
	if time.Now().After(otp.ExpiredAt) {
		return errors.New("kode OTP sudah kedaluwarsa, silakan minta kode baru")
	}

	// Lanjut proses hash password dan simpan user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	err = config.DB.Create(user).Error
	if err == nil {
		// Hapus OTP setelah berhasil digunakan
		config.DB.Delete(&otp)
	}
	return err
}

func HandleGoogleLogin(email string) (string, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("akun belum terdaftar, silakan registrasi terlebih dahulu")
	}

	return utils.GenerateToken(user.ID, user.Email)
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
