package utils

import "time"

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func GetFirstDayOfMonth() string {
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return firstDay.Format("2006-01-02")
}
