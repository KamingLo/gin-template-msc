package services

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"template/config"
	"template/models"
	"template/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RequestOTP(email string) error {
	// 1. Cek apakah email sudah punya akun
	var existingUser models.User
	err := config.DB.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		return errors.New("Email sudah terdaftar, silakan langsung login")
	}

	// 2. Generate 6 digit kode random
	code := fmt.Sprintf("%06d", rand.IntN(900000)+100000)

	// 3. Persiapkan data OTP
	otp := models.OTP{
		Email:     email,
		Code:      code,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	// 4. Hapus OTP lama buat email ini (biar gak numpuk di DB)
	config.DB.Where("email = ?", email).Delete(&models.OTP{})

	// 5. Simpan ke database
	if err := config.DB.Create(&otp).Error; err != nil {
		return errors.New("Gagal membuat sesi verifikasi")
	}

	// 6. EKSEKUSI KIRIM EMAIL
	if err := SendRegistrationOTP(email, code); err != nil {
		// Jika email gagal, kita hapus lagi OTP-nya biar konsisten
		config.DB.Delete(&otp)
		return errors.New("Gagal mengirim email, pastikan alamat email benar")
	}

	return nil
}

func RegisterWithOTP(input *models.User, otpCode string) error {
	var otp models.OTP
	err := config.DB.Where("email = ? AND code = ?", input.Email, otpCode).First(&otp).Error
	if err != nil || time.Now().After(otp.ExpiredAt) {
		return errors.New("Kode OTP salah atau kedaluwarsa")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	// GORM akan otomatis menjalankan BeforeCreate untuk GenerateCustomID
	if err := config.DB.Create(input).Error; err != nil {
		return errors.New("Gagal menyimpan akun")
	}

	config.DB.Delete(&otp)
	return nil
}

func LoginUser(input models.UserLogin) (string, error) {
	var user models.User

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return "", errors.New("Email salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", errors.New("Password salah")
	}

	// user.ID sekarang bertipe string
	return utils.GenerateToken(user.ID, user.Email)
}

func HandleGoogleLogin(email string) (string, error) {
	var user models.User

	// Cari user berdasarkan email
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		// Jika tidak ditemukan, kembalikan error agar controller tahu user belum terdaftar
		return "", err
	}

	// Jika user ditemukan, buatkan token
	return utils.GenerateToken(user.ID, user.Email)
}
