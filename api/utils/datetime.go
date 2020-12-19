package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func GetFirstLastDateCurrentMonth(change string) (string, string, error) {
	now := time.Now()
	if change != "" {
		change, err := strconv.Atoi(change)
		if err != nil {
			return "", "", err
		}
		now = now.AddDate(0, change, 0)
	}
	currentYear, currentMonth, _ := now.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())
	return firstOfMonth.Format("2006-01-02"), firstOfMonth.AddDate(0, 1, -1).Format("2006-01-02"), nil
}

func ParseDate(date string) (string, error) {
	result, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", errors.New("Date format must be: YYYY-MM-DD")
	}
	return result.Format("2006-01-02"), nil
}

func ParseFromToDates(from, to string) (string, string, error) {
	var err error
	from, err = ParseDate(from)
	if err != nil {
		return "", "", errors.New(fmt.Sprint("It could not parse 'from' date: ", err.Error()))
	}
	to, err = ParseDate(to)
	if err != nil {
		return "", "", errors.New(fmt.Sprint("It could not parse 'to' date: ", err.Error()))
	}
	return from, to, nil
}
