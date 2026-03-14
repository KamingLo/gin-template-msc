package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateCustomID menghasilkan format: PREFIX-YYYYMMDD-RANDOM
func GenerateCustomID(prefix string, length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Gunakan Uppercase agar seragam di URL
	datePart := time.Now().Format("20060102")

	result := make([]byte, length)
	// Langsung loop untuk mengisi karakter
	for i := 0; i < length; i++ {
		randomByte := make([]byte, 1)
		rand.Read(randomByte)
		result[i] = charset[randomByte[0]%uint8(len(charset))]
	}

	return fmt.Sprintf("%s-%s-%s", prefix, datePart, string(result))
}
