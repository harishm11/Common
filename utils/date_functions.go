package utils

import "time"

func GetCurrentDateTime() time.Time {
	return time.Now()
}

func FormatDate(date time.Time, layout string) string {
	return date.Format(layout)
}

func ParseDate(dateStr, layout string) (time.Time, error) {
	return time.Parse(layout, dateStr)
}

func AddDays(date time.Time, days int) time.Time {
	return date.AddDate(0, 0, days)
}

func DaysBetween(startDate, endDate time.Time) int {
	duration := endDate.Sub(startDate)
	return int(duration.Hours() / 24)
}

func IsWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func StartOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func EndOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, int(time.Second-1), date.Location())
}
