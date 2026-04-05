package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"template/config"
	"template/models"
	"template/services"
	"template/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// GoogleLogin mengarahkan user ke halaman login Google
func GoogleLogin(c *gin.Context) {
	// State sebaiknya di-generate acak dan disimpan di session/cookie untuk keamanan CSRF
	state := os.Getenv("SESSION_SECRET")
	url := config.GoogleOAuthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("prompt", "select_account"),
	)

	log.Printf("[OAuth Init] Google Auth URL generated")
	utils.SendSuccess(c, http.StatusOK, "URL Successfully made", gin.H{"url": url})
}

// GoogleCallback menerima response dari Google
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	// Validasi state di sini jika kamu menyimpannya di session
	if state != os.Getenv("SESSION_SECRET") {
		utils.SendError(c, http.StatusBadRequest, "Invalid state", nil)
		return
	}

	// 1. Tukar Code dengan Token
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("[OAuth Error] Exchange failed: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("OAUTH_FRONTEND_URL")+"?error=token_exchange_failed")
		return
	}

	// 2. Ambil data user dari Google API menggunakan token tersebut
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("[OAuth Error] Failed to get user info: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("OAUTH_FRONTEND_URL")+"?error=fetch_user_failed")
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var googleUser struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(data, &googleUser); err != nil {
		log.Printf("[OAuth Error] Parse JSON failed: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("OAUTH_FRONTEND_URL")+"?error=parse_failed")
		return
	}

	// 3. Logika Backend (Sama seperti sebelumnya)
	jwtToken, err := services.HandleGoogleLogin(googleUser.Email)
	if err != nil {
		loginURL := os.Getenv("OAUTH_FRONTEND_URL") + "?error=user_not_registered&email=" + googleUser.Email
		c.Redirect(http.StatusTemporaryRedirect, loginURL)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("SUCCESS_FRONTEND_URL")+"?token="+jwtToken)
}

func RequestOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Email format is wrong", err)
		return
	}

	if err := services.RequestOTP(input.Email); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Failed to send otp", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "OTP Code is sent to your email", nil)
}

func Register(c *gin.Context) {
	var input struct {
		models.User
		OTPCode string `json:"otp_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Data is not complete", err)
		return
	}

	if err := services.RegisterWithOTP(&input.User, input.OTPCode); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Registration Failed", err)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Registration Success, Please Login", nil)
}

func Login(c *gin.Context) {
	var input models.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Input is not valid", err)
		return
	}

	token, err := services.LoginUser(input)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "Email or password is incorrect", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Login Successfully", gin.H{
		"token": token,
	})
}

func GetMe(c *gin.Context) {
	id, _ := c.Get("user_id")
	email, _ := c.Get("user_email")

	utils.SendSuccess(c, http.StatusOK, "Data profilmu berhasil diambil", gin.H{
		"id":    id,
		"email": email,
	})
}

// Di controller/handler auth kamu
func Logout(c *gin.Context) {
	// Tanpa Goth, kita tidak perlu membersihkan session di sisi server
	// jika kita menggunakan stateless JWT.
	utils.SendSuccess(c, http.StatusOK, "Logged out successfully", nil)
}

func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Format email salah", err)
		return
	}

	if err := services.ForgotPassword(input.Email); err != nil {
		utils.SendError(c, http.StatusTooManyRequests, err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Tautan reset telah dikirim ke email", nil)
}

func ResetPassword(c *gin.Context) {
	var input struct {
		Email       string `json:"email" binding:"required,email"`
		Token       string `json:"token" binding:"required"` // Tambahkan Token di sini
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Data tidak lengkap", err)
		return
	}

	if err := services.ResetPassword(input.Email, input.Token, input.NewPassword); err != nil {
		utils.SendError(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Password berhasil diperbarui", nil)
}
