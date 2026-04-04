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
		utils.SendError(c, http.StatusInternalServerError, "Failed to initiate google auth", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "URL Successfully made", gin.H{"url": url})
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
