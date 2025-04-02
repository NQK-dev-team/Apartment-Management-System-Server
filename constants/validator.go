package constants

import (
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func InitCustomValidationRules() {
	Validate.RegisterValidation("password", ValidatePassword)
	Validate.RegisterValidation("dob_18", ValidateDoB)
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

func ValidateDoB(fl validator.FieldLevel) bool {
	const DateFormat = "2006-01-02"
	dobString := fl.Field().String()

	if dobString == "" {
		return true
	}

	// --- 1. Parse the string into a time.Time object ---
	dob, err := time.Parse(DateFormat, dobString)
	if err != nil {
		return false // Invalid date format
	}

	// --- 2. Perform Age and Validity Checks ---
	now := time.Now().UTC() // Use UTC for consistent comparison
	dobUTC := dob.UTC()     // Ensure dob is also in UTC for comparison

	// Check if dob is in the future
	if dobUTC.After(now) {
		return false
	}

	// Calculate the cutoff date (exactly 18 years before now in UTC)
	cutoffDate := now.AddDate(-18, 0, 0)

	// Check if the date of birth is on or before the cutoff date (18+ years old)
	// !dobUTC.After(cutoffDate) is equivalent to dobUTC <= cutoffDate
	isAdult := !dobUTC.After(cutoffDate)

	if !isAdult {
		return false // Age is less than 18
	}

	return true
}
