package utils

import "time"

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func GetFirstDayOfMonth(month string) string {
	now := time.Now()
	var firstDay time.Time

	if month == "" {
		firstDay = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	} else {
		m, err := time.Parse("2006-01", month)
		if err != nil {
			return ""
		}
		firstDay = time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	return firstDay.Format("2006-01-02")
}

func GetLastDayOfMonth(month string) string {
	now := time.Now()
	var lastDay time.Time

	if month == "" {
		lastDay = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())
	} else {
		m, err := time.Parse("2006-01", month)
		if err != nil {
			return ""
		}
		lastDay = time.Date(m.Year(), m.Month()+1, 0, 0, 0, 0, 0, now.Location())
	}

	return lastDay.Format("2006-01-02")
}

func GetFirstDayOfQuarter() string {
	now := time.Now()
	month := now.Month()
	var firstMonth time.Month

	switch month {
	case time.January, time.February, time.March:
		firstMonth = time.January
	case time.April, time.May, time.June:
		firstMonth = time.April
	case time.July, time.August, time.September:
		firstMonth = time.July
	case time.October, time.November, time.December:
		firstMonth = time.October
	}

	firstDay := time.Date(now.Year(), firstMonth, 1, 0, 0, 0, 0, now.Location())
	return firstDay.Format("2006-01-02")
}

func CompareDates(date1, date2 string) (int, error) {
	t1, err := time.Parse("2006-01-02", date1)
	if err != nil {
		return 2, err
	}
	t2, err := time.Parse("2006-01-02", date2)
	if err != nil {
		return 2, err
	}

	if t1.Before(t2) {
		return -1, nil
	} else if t1.After(t2) {
		return 1, nil
	}
	return 0, nil
}
