package utils

import (
	"api/config"
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(inputStr string) (string, error) {
	// Use bcrypt algorithm to hash the input string
	roundConfig := config.GetEnv("BCRYPT_ROUNDS")
	if roundConfig == "" {
		roundConfig = "12"
	}
	rounds, err := strconv.Atoi(roundConfig)
	if err != nil {
		return "", err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(inputStr), rounds)
	return string(bytes), err
}

func CompareHashPassword(hashedStr string, inputStr string) bool {
	// Compare the hashed string (by bcrypt) with the input string
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(inputStr))
	return err == nil
}

func HashString(inputStr string) (string, error) {
	// Use SHA-256 algorithm to hash the input string
	h := sha256.New()
	_, err := h.Write([]byte(inputStr))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func CompareHashString(hashedStr string, inputStr string) bool {
	// Compare the hashed string (by SHA-256) with the input string
	inputStrAfterHash, err := HashString(inputStr)

	if err != nil {
		return false
	}

	return hashedStr == inputStrAfterHash
}
