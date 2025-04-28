package config

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID genera un ID único de 20 caracteres
func GenerateID() (string, error) {
	bytes := make([]byte, 10)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
