package controllers

import (
	"fmt"
	"net/http"
	"os"
	"template/models"
	"template/services"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func RequestOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format email tidak valid"})
		return
	}

	if err := services.RequestOTP(input.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kode OTP telah dikirim ke email kamu"})
}

func Register(c *gin.Context) {
	var input struct {
		models.User
		OTPCode string `json:"otp_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.RegisterWithOTP(&input.User, input.OTPCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Berhasil registrasi data baru"})
}

func GoogleLogin(c *gin.Context) {
	// Gothic butuh response writer dan request standar http
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func GoogleCallback(c *gin.Context) {
	// Tetap tambahkan provider secara manual untuk gothic
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	// Ambil URL dari environment
	baseURL := os.Getenv("FRONTEND_URL")
	successURL := os.Getenv("SUCCESS_FRONTEND_URL")

	// 1. JIKA AUTENTIKASI GOOGLE GAGAL
	// Arahkan kembali ke halaman login frontend dengan pesan error
	if err != nil {
		target := fmt.Sprintf("%s/login?error=auth_failed", baseURL)
		c.Redirect(http.StatusTemporaryRedirect, target)
		return
	}

	token, err := services.HandleGoogleLogin(user.Email)

	// 2. JIKA EMAIL BELUM TERDAFTAR (HARUS REGISTRASI)
	// Arahkan ke halaman register frontend dan bawa email user
	if err != nil {
		target := fmt.Sprintf("%s/register?email=%s&message=complete_your_profile", baseURL, user.Email)
		c.Redirect(http.StatusTemporaryRedirect, target)
		return
	}

	// 3. JIKA LOGIN BERHASIL
	// Arahkan ke halaman sukses (Dashboard) dengan membawa token JWT
	target := fmt.Sprintf("%s?token=%s", successURL, token)
	c.Redirect(http.StatusTemporaryRedirect, target)
}

func Login(c *gin.Context) {
	var input models.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.LoginUser(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login berhasil", "token": token})
}
