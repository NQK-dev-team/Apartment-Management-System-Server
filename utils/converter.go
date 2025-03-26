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
