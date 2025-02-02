package utils

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func InitCustomValidationRules() {
	Validate.RegisterValidation("password", ValidatePassword)
}

func customPasswordRule(password string) bool {
	var hasUpper, hasLower, hasNumber, hasSpecial = false, false, false, false
	specialChars := "#@$!%*?&"

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Custom validator for password
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return len(password) >= 8 && customPasswordRule(password)
}
