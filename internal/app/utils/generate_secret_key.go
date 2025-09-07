package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

const secretKeyLength = 32

// GenerateRandomSecretKey generates a cryptographically secure random secret key.
// The key is 32 bytes long and encoded as base64 URL-safe string.
// This function is typically used for JWT signing keys.
func GenerateRandomSecretKey() string {
	b := make([]byte, secretKeyLength)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("Failed to generate random secret key:", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
