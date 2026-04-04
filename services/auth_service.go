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
	var existingUser models.User
	err := config.DB.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		return errors.New("email is already registered, please log in")
	}

	var existingOTP models.OTP
	errOTP := config.DB.Where("email = ?", email).First(&existingOTP).Error

	code := fmt.Sprintf("%06d", rand.IntN(900000)+100000)
	now := time.Now()

	if errOTP == nil {
		var cooldown time.Duration

		switch existingOTP.RequestCount {
		case 1:
			cooldown = 30 * time.Second
		case 2:
			cooldown = 60 * time.Second
		case 3:
			cooldown = 5 * time.Minute
		default:
			cooldown = 1 * time.Hour
		}

		timeSinceLast := time.Since(existingOTP.UpdatedAt)

		if timeSinceLast > 24*time.Hour {
			existingOTP.RequestCount = 0
		} else if timeSinceLast < cooldown {
			timeLeft := int(cooldown.Seconds() - timeSinceLast.Seconds())

			var timeMsg string
			if timeLeft < 60 {
				timeMsg = fmt.Sprintf("%d seconds", timeLeft)
			} else if timeLeft < 3600 {
				timeMsg = fmt.Sprintf("%d minutes %d seconds", timeLeft/60, timeLeft%60)
			} else {
				timeMsg = fmt.Sprintf("%d hours %d minutes", timeLeft/3600, (timeLeft%3600)/60)
			}

			return fmt.Errorf("too many requests; try again in %s", timeMsg)
		}

		existingOTP.Code = code
		existingOTP.ExpiredAt = now.Add(5 * time.Minute)
		existingOTP.RequestCount += 1

		if err := config.DB.Save(&existingOTP).Error; err != nil {
			return errors.New("failed to update verification session")
		}

	} else {
		newOTP := models.OTP{
			Email:        email,
			Code:         code,
			ExpiredAt:    now.Add(5 * time.Minute),
			RequestCount: 1,
		}

		if err := config.DB.Create(&newOTP).Error; err != nil {
			return errors.New("failed to create verification session")
		}
	}

	// Assuming SendRegistrationOTP is defined in your utils or internal helper
	if err := SendRegistrationOTP(email, code); err != nil {
		return errors.New("failed to send email; ensure the address is correct")
	}

	return nil
}

func RegisterWithOTP(input *models.User, otpCode string) error {
	var otp models.OTP
	err := config.DB.Where("email = ? AND code = ?", input.Email, otpCode).First(&otp).Error
	if err != nil || time.Now().After(otp.ExpiredAt) {
		return errors.New("invalid or expired OTP code")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	if err := config.DB.Create(input).Error; err != nil {
		return errors.New("failed to save user account")
	}

	config.DB.Delete(&otp)
	return nil
}

func LoginUser(input models.UserLogin) (string, error) {
	var user models.User

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return "", errors.New("email not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", errors.New("incorrect password")
	}

	return utils.GenerateToken(user.ID, user.Email)
}

func HandleGoogleLogin(email string) (string, error) {
	var user models.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	return utils.GenerateToken(user.ID, user.Email)
}
