package database

import (
	"fmt"
	"time"
)

const TimeFormat = time.RFC3339Nano

func FormatTime(value time.Time) string {
	return value.UTC().Format(TimeFormat)
}

func ParseTime(value string) (time.Time, error) {
	parsed, err := time.Parse(TimeFormat, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse database time %q: %w", value, err)
	}
	return parsed, nil
}
