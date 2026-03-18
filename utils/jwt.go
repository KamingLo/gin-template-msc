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

	// Ambil konfigurasi dari env
	isExpires := os.Getenv("JWT_EXPIRES")
	expiresInStr := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"iat":     time.Now().Unix(), // Issued At (Kapan token dibuat)
	}

	// Logika Expired
	if isExpires == "enable" {
		// Default 2 jam jika env tidak valid atau kosong
		hours := 2

		// Coba konversi string env ke int
		if val, err := strconv.Atoi(expiresInStr); err == nil {
			hours = val
		}

		// Tambahkan field "exp" ke claims
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
