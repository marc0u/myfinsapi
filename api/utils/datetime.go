package utils

import (
	"errors"
	"time"
)

func GetFirstLastDateCurrentMonth() (string, string) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return firstOfMonth.Format("2006-01-02"), lastOfMonth.Format("2006-01-02")
}

func ParseDate(date string) (string, error) {
	result, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("Date format must be: YYYYMMDD")
	}
	return result.Format("2006-01-02"), nil
}
