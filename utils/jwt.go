package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, email string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	isExpires := os.Getenv("JWT_EXPIRES")
	expiresInStr := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"iat":     time.Now().Unix(),
	}

	if isExpires == "enable" {
		hours := 2

		if val, err := strconv.Atoi(expiresInStr); err == nil {
			hours = val
		}

		claims["exp"] = time.Now().Add(time.Hour * time.Duration(hours)).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}
