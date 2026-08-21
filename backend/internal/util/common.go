package util

import "time"

func ParseDateTimeLocal(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	t, err := time.ParseInLocation(
		"2006-01-02T15:04",
		value,
		loc,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
