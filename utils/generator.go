package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateString(numberOfChar int) (string, error) {
	bytes := make([]byte, numberOfChar/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
