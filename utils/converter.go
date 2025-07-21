package utils

import (
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
	return sql.NullTime{
		Time:  ParseTime(str),
		Valid: true,
	}
}

func ParseTime(str string) time.Time {
	const layout = "2006-01-02"
	parsedTime, _ := time.Parse(layout, str)
	return parsedTime
}

// func ParseTimeWithZone(str string) time.Time {
// 	const layout = "2006-01-02 15:04:05"

// 	timeZone := config.GetEnv("APP_TIMEZONE")

// 	if timeZone == "" {
// 		timeZone = "Asia/Ho_Chi_Minh"
// 	}

// 	timeLocation, err := time.LoadLocation(timeZone)
// 	if err != nil {
// 		timeLocation = time.Local // Fallback to local time if loading the specified timezone fails
// 	}

// 	parsedTime, _ := time.ParseInLocation(layout, str, timeLocation)
// 	return parsedTime
// }
