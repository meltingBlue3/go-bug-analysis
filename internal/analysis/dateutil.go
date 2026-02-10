package analysis

import (
	"strings"
	"time"
)

// ParseDate attempts to parse a Zentao date string.
// It tries "2006-01-02 15:04:05" first, then "2006-01-02".
// Returns zero time and false for empty, "0000-00-00" prefix, or unparseable strings.
func ParseDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasPrefix(s, "0000-00-00") {
		return time.Time{}, false
	}

	if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		return t, true
	}

	return time.Time{}, false
}

// DateOnly truncates a time.Time to midnight (year, month, day).
func DateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// Today returns today's date at midnight in the local timezone.
func Today() time.Time {
	return DateOnly(time.Now())
}

// Yesterday returns yesterday's date at midnight in the local timezone.
func Yesterday() time.Time {
	return DateOnly(time.Now().AddDate(0, 0, -1))
}
