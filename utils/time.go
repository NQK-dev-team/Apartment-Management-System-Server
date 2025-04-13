package utils

import "time"

func GetCurrentYear() int {
	return time.Now().Year()
}

func GetCurrentQuarter() int {
	month := time.Now().Month()
	switch {
	case month >= 1 && month <= 3:
		return 1
	case month >= 4 && month <= 6:
		return 2
	case month >= 7 && month <= 9:
		return 3
	case month >= 10 && month <= 12:
		return 4
	default:
		return 0 // Invalid month
	}
}
