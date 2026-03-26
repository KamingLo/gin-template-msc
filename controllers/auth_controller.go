package controllers

import (
	"net/http"
	"os"
	"template/models"
	"template/services"
	"template/utils"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func GoogleLogin(c *gin.Context) {
	platform := c.DefaultQuery("platform", "web")

	sess, _ := gothic.Store.Get(c.Request, "auth-session")
	sess.Values["platform"] = platform
	sess.Save(c.Request, c.Writer)

	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	url, err := gothic.GetAuthURL(c.Writer, c.Request)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Gagal inisiasi Google Auth", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "URL Auth berhasil dibuat", gin.H{"url": url})
}

func GoogleCallback(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		loginURL := os.Getenv("OAUTH_FRONTEND_URL") + "?error=google_auth_failed"
		c.Redirect(http.StatusTemporaryRedirect, loginURL)
		return
	}

	sess, _ := gothic.Store.Get(c.Request, "auth-session")
	platform, _ := sess.Values["platform"].(string)

	token, err := services.HandleGoogleLogin(user.Email)

	if err != nil {
		errorMessage := "user_not_registered"
		loginURL := os.Getenv("OAUTH_FRONTEND_URL") + "?error=" + errorMessage + "&email=" + user.Email

		if platform == "mobile" {
			c.Redirect(http.StatusTemporaryRedirect, "myapp://login?error="+errorMessage)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, loginURL)
		return
	}

	if platform == "mobile" {
		c.Redirect(http.StatusTemporaryRedirect, "myapp://auth?token="+token)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("SUCCESS_FRONTEND_URL")+"?token="+token)
}

func RequestOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Format email salah", err)
		return
	}

	if err := services.RequestOTP(input.Email); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Gagal mengirim OTP", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Kode OTP telah dikirim ke email kamu", nil)
}

func Register(c *gin.Context) {
	var input struct {
		models.User
		OTPCode string `json:"otp_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Data tidak lengkap", err)
		return
	}

	if err := services.RegisterWithOTP(&input.User, input.OTPCode); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Gagal registrasi", err)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Registrasi berhasil, silakan login", nil)
}

func Login(c *gin.Context) {
	var input models.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Input tidak valid", err)
		return
	}

	token, err := services.LoginUser(input)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "Email atau password salah", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
	})
}

func Logout(c *gin.Context) {
	utils.SendSuccess(c, http.StatusOK, "Berhasil keluar", nil)
}

func GetMe(c *gin.Context) {
	id, _ := c.Get("user_id")
	email, _ := c.Get("user_email")

	utils.SendSuccess(c, http.StatusOK, "Data profil berhasil diambil", gin.H{
		"id":    id,
		"email": email,
	})
}
