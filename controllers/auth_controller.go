package controllers

import (
	"net/http"
	"os"
	"template/models"
	"template/services"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// HELPER: Set HttpOnly Cookie
func setAuthCookie(c *gin.Context, token string) {
	isProd := os.Getenv("GIN_MODE") == "release"
	domain := os.Getenv("COOKIE_DOMAIN") // Ambil dari ENV (localhost / domain.com)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, 3600*24, "/", domain, isProd, true)
}

// GOOGLE LOGIN: Ambil URL untuk Frontend
func GoogleLogin(c *gin.Context) {
	platform := c.DefaultQuery("platform", "web")

	// Simpan platform di session gothic agar terbaca di callback
	sess, _ := gothic.Store.Get(c.Request, "auth-session")
	sess.Values["platform"] = platform
	sess.Save(c.Request, c.Writer)

	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	url, err := gothic.GetAuthURL(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal inisiasi Google Auth"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GOOGLE CALLBACK: Handle Redirect dari Google
func GoogleCallback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("OAUTH_FRONTEND_URL")+"?error=failed")
		return
	}

	token, _ := services.HandleGoogleLogin(user.Email)

	sess, _ := gothic.Store.Get(c.Request, "auth-session")
	platform, _ := sess.Values["platform"].(string)

	if platform == "mobile" {
		// Khusus Flutter: Deep Link
		c.Redirect(http.StatusTemporaryRedirect, "myapp://auth?token="+token)
		return
	}

	// Khusus Web: Cookie & Redirect
	setAuthCookie(c, token)
	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("SUCCESS_FRONTEND_URL"))
}

// REQUEST OTP HANDLER
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

// REGISTER HANDLER
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

// LOGIN HANDLER
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

	setAuthCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"message": "Login berhasil", "token": token})
}

// LOGOUT HANDLER
func Logout(c *gin.Context) {
	domain := os.Getenv("COOKIE_DOMAIN")
	c.SetCookie("auth_token", "", -1, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil keluar"})
}

// GET ME HANDLER (Dapatkan profil dari context middleware)
func GetMe(c *gin.Context) {
	id, _ := c.Get("user_id")
	email, _ := c.Get("user_email")

	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"email": email,
	})
}
