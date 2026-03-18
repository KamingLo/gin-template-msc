package controllers

import (
	"net/http"
	"os"
	"template/models"
	"template/services"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// GOOGLE LOGIN: Mendapatkan URL autentikasi untuk dikirim ke Frontend
func GoogleLogin(c *gin.Context) {
	platform := c.DefaultQuery("platform", "web")

	// Simpan platform di session gothic agar bisa dibaca saat callback
	sess, _ := gothic.Store.Get(c.Request, "auth-session")
	sess.Values["platform"] = platform
	sess.Save(c.Request, c.Writer)

	// Tambahkan provider google secara manual ke query request
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	url, err := gothic.GetAuthURL(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal inisiasi Google Auth"})
		return
	}

	// Kirim URL ke frontend agar frontend yang melakukan window.location.href
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GOOGLE CALLBACK: Menangani kembalian dari Google
func GoogleCallback(c *gin.Context) {
	// PENTING: Gothic butuh query 'provider' di URL Callback
	// Jika route kamu adalah /auth/google/callback, tambahkan ini:
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		// Jika error di sini, biasanya karena session/cookie Goth hilang
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("OAUTH_FRONTEND_URL")+"?error=failed_to_complete_auth")
		return
	}

	token, _ := services.HandleGoogleLogin(user.Email)

	sess, _ := gothic.Store.Get(c.Request, "auth-session")
	platform, _ := sess.Values["platform"].(string)

	if platform == "mobile" {
		c.Redirect(http.StatusTemporaryRedirect, "myapp://auth?token="+token)
		return
	}

	// Redirect ke Next.js API Callback
	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("SUCCESS_FRONTEND_URL")+"?token="+token)
}

// REQUEST OTP: Mengirim kode OTP ke email
func RequestOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format email salah"})
		return
	}

	if err := services.RequestOTP(input.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kode OTP telah dikirim ke email kamu"})
}

// REGISTER: Membuat akun baru dengan verifikasi OTP
func Register(c *gin.Context) {
	var input struct {
		models.User
		OTPCode string `json:"otp_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap"})
		return
	}

	if err := services.RegisterWithOTP(&input.User, input.OTPCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil, silakan login"})
}

// LOGIN: Mengembalikan token JWT dalam bentuk JSON
func Login(c *gin.Context) {
	var input models.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	token, err := services.LoginUser(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Output murni JSON, Next.js yang akan simpan ke HttpOnly Cookie
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
	})
}

// LOGOUT: Stateless logout hanya memberikan respon sukses
func Logout(c *gin.Context) {
	// Karena stateless, backend tidak perlu menghapus cookie.
	// Cukup berikan instruksi sukses agar frontend menghapus token di sisinya.
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil keluar"})
}

// GET ME: Mengambil data user yang sedang login
func GetMe(c *gin.Context) {
	// Data ini didapat dari Middleware Auth yang memparsing Token
	id, _ := c.Get("user_id")
	email, _ := c.Get("user_email")

	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"email": email,
	})
}
