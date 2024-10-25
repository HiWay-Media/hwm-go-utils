package utils

import "time"

// FormatDate formats the date to "YYYY-MM-DD" by default, or uses a custom layout if provided.
func FormatDate(t time.Time, layout string) string {
	if layout == "" {
		layout = "2006-01-02"
	}
	return t.Format(layout)
}

// ParseDate parses a date string with a specified layout, or "YYYY-MM-DD" by default.
func ParseDate(dateStr, layout string) (time.Time, error) {
	if layout == "" {
		layout = "2006-01-02"
	}
	return time.Parse(layout, dateStr)
}

// DaysBetween calculates the difference in days between two dates.
func DaysBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Hours() / 24)
}

// IsWeekend checks if the given date is a Saturday or Sunday.
func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

// AddDays adds the specified number of days to the date.
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddMonths adds the specified number of months to the date.
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddYears adds the specified number of years to the date.
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// DaysInMonth returns the number of days in the specified month of the specified year.
func DaysInMonth(year int, month time.Month) int {
	// Go to the first day of the next month, then go back one day
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}