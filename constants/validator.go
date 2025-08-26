package constants

import (
	"api/config"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func InitCustomValidationRules() {
	Validate.RegisterValidation("password", ValidatePassword)
	Validate.RegisterValidation("dob_18", ValidateDoB)
	Validate.RegisterValidation("image_type", ValidateImageType)
	Validate.RegisterValidation("file_type", ValidateFileType)
	Validate.RegisterValidation("upload_type", ValidateUploadType)
	Validate.RegisterValidation("image_size", ValidateImageSize)
	Validate.RegisterValidation("file_size", ValidateFileSize)
	Validate.RegisterValidation("upload_size", ValidateUploadSize)
	Validate.RegisterValidation("check_date_equal_or_after", CheckDateEqualOrAfter)
	Validate.RegisterValidation("check_date_equal_or_before", CheckDateEqualOrBefore)
	Validate.RegisterValidation("contract_type_and_end_date", ValidateContractTypeAndEndDate)
	Validate.RegisterValidation("not_after_current_date", ValidateNotAfterCurrentDate)
	Validate.RegisterValidation("validate_payment_time", ValidatePaymentTime)
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

func ValidateImageType(fl validator.FieldLevel) bool {
	fileType := fl.Field().String()

	// Check if the file's MIME type is in the list of allowed types
	for _, t := range Common.FileUpload.AllowedImageTypes {
		if t == fileType {
			return true
		}
	}

	return false
}

func ValidateFileType(fl validator.FieldLevel) bool {
	fileType := fl.Field().String()

	// Check if the file's MIME type is in the list of allowed types
	for _, t := range Common.FileUpload.AllowedFileTypes {
		if t == fileType {
			return true
		}
	}

	return false
}

func ValidateUploadType(fl validator.FieldLevel) bool {
	fileType := fl.Field().String()

	// Check if the file's MIME type is in the list of allowed types
	for _, t := range Common.FileUpload.AllowedUploadTypes {
		if t == fileType {
			return true
		}
	}

	return false
}

func ValidateImageSize(fl validator.FieldLevel) bool {
	fileSize := fl.Field().Interface().(int64)

	return fileSize < Common.FileUpload.MaxImageSize
}

func ValidateFileSize(fl validator.FieldLevel) bool {
	fileSize := fl.Field().Interface().(int64)

	return fileSize < Common.FileUpload.MaxFileSize
}

func ValidateUploadSize(fl validator.FieldLevel) bool {
	fileSize := fl.Field().Interface().(int64)

	return fileSize < Common.FileUpload.MaxUploadSize
}

func GetValidateErrorMessage(err error) string {
	appEnv := config.GetEnv("APP_ENV")

	if appEnv == "development" {
		return err.Error()
	}

	return ""
}

func CheckDateEqualOrAfter(fl validator.FieldLevel) bool {
	fieldName := fl.Param() // Get the referenced field (StartDate)
	startField := fl.Parent().FieldByName(fieldName)
	endStr := fl.Field().String()
	startStr := startField.String()

	startDate, err1 := time.Parse("2006-01-02", startStr)
	endDate, err2 := time.Parse("2006-01-02", endStr)

	if err1 != nil {
		return true // Other validation rules will catch invalid date formats, not this one
	}

	if err2 != nil {
		return false
	}

	return !startDate.UTC().After(endDate.UTC())
}

func CheckDateEqualOrBefore(fl validator.FieldLevel) bool {
	fieldName := fl.Param() // Get the referenced field (EndDate)
	endField := fl.Parent().FieldByName(fieldName)
	startStr := fl.Field().String()
	endStr := endField.String()

	startDate, err1 := time.Parse("2006-01-02", startStr)
	endDate, err2 := time.Parse("2006-01-02", endStr)

	if err2 != nil {
		return true // Other validation rules will catch invalid date formats, not this one
	}

	if err1 != nil {
		return false
	}

	return !endDate.UTC().Before(startDate.UTC())
}

func ValidateContractTypeAndEndDate(fl validator.FieldLevel) bool {
	fieldName := fl.Param()
	contractType := fl.Parent().FieldByName(fieldName).Int()
	endDate := fl.Field().String()

	if int(contractType) == Common.ContractType.RENT && endDate == "" {
		return false // EndDate is required for Rent contracts
	} else if int(contractType) == Common.ContractType.BUY && endDate != "" {
		return false // EndDate should not be provided for Buy contracts
	}

	return true
}

func ValidateNotAfterCurrentDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		return true // If the field is empty, we consider it valid
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false // Invalid date format
	}

	currentDate := time.Now().UTC()
	return !date.UTC().After(currentDate)
}

func ValidatePaymentTime(fl validator.FieldLevel) bool {
	paymentTimeStr := fl.Field().String()
	if paymentTimeStr != "" {
		paymentTime, err := time.Parse("2006-01-02", paymentTimeStr)
		if err != nil {
			return false // Invalid date format
		}

		currentDate := time.Now().UTC()

		if paymentTime.UTC().After(currentDate) {
			return false // PaymentTime cannot be in the future
		}

		// Get value of Period
		periodStr := fl.Parent().FieldByName("Period").String()
		period, err := time.Parse("2006-01-02", periodStr+"-01")
		if err != nil {
			return true // Do not validate if Period is invalid, the validation rule for Period will handle it
		}

		if paymentTime.UTC().Before(period.UTC()) {
			return false // PaymentTime cannot be before Period
		}
	}

	return true
}
