package utils

import (
	"api/config"
	"database/sql"
	"time"
)

func StringToNullTime(str string) sql.NullTime {
	if str == "" {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}

	parsedTime, err := ParseTime(str)
	if err != nil {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}

	return sql.NullTime{
		Time:  parsedTime,
		Valid: true,
	}
}

func ParseTime(str string) (time.Time, error) {
	const layout = "2006-01-02"
	return time.Parse(layout, str)
}

func ParseTimeWithZone(str string) (time.Time, error) {
	const layout = "2006-01-02 15:04:05"

	timeZone := config.GetEnv("APP_TIMEZONE")

	if timeZone == "" {
		timeZone = "Asia/Ho_Chi_Minh"
	}

	timeLocation, err := time.LoadLocation(timeZone)
	if err != nil {
		timeLocation = time.Local // Fallback to local time if loading the specified timezone fails
	}

	return time.ParseInLocation(layout, str, timeLocation)
}
