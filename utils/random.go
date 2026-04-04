package utils

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"
)

func GenerateCustomID(prefix string, length int) string {

	machine_id := os.Getenv("MACHINE_ID")
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	datePart := time.Now().Format("20060102")

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		randomByte := make([]byte, 1)
		rand.Read(randomByte)
		result[i] = charset[randomByte[0]%uint8(len(charset))]
	}

	return fmt.Sprintf("%s-%s-%s-%s", prefix, datePart, machine_id, string(result))
}
