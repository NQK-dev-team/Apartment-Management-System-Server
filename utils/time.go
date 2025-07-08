package utils

import "time"

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func GetFirstDayOfMonth(month string) string {
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

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
