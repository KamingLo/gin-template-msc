package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateCustomID menghasilkan format: PREFIX-YYYYMMDD-RANDOM
func GenerateCustomID(prefix string, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	datePart := time.Now().Format("20060102")

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// Fallback jika crypto/rand gagal
		return fmt.Sprintf("%s-%s-error", prefix, datePart)
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return fmt.Sprintf("%s-%s-%s", prefix, datePart, string(b))
}
