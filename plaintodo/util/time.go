package util

import (
	"strconv"
	"time"
)

// DateTimeFormat is 2006-01-02 15:04
var DateTimeFormat = "2006-01-02 15:04"

// DateFormat is 2006-01-02 not have time format
var DateFormat = "2006-01-02"

// ParseTime convert string which DateTimeFormat or DateFormat convert to time object
func ParseTime(dateString string) (time.Time, bool) {
	var t time.Time
	t, err := time.Parse(DateTimeFormat+"-0700", dateString+"+0900")
	if err != nil {
		t, err = time.Parse(DateFormat+"-0700", dateString+"+0900")
		if err != nil {
			// not date value
			return t, false
		}
	}

	return t, true
}

// AddDuration add time
// num is numeric as string
// unit is time unit string as minutes, hour, day, week...
func AddDuration(base time.Time, num string, unit string) time.Time {
	n, err := strconv.Atoi(num)
	if err != nil {
		return time.Unix(0, 0)
	}
	switch {
	case unit == "minutes":
		return base.Add(time.Duration(n) * time.Minute)
	case unit == "hour":
		return base.Add(time.Duration(n) * time.Hour)
	case unit == "day":
		return base.AddDate(0, 0, n)
	case unit == "week":
		return base.AddDate(0, 0, n*7)
	case unit == "month":
		return base.AddDate(0, n, 0)
	case unit == "year":
		return base.AddDate(n, 0, 0)
	}

	return time.Unix(0, 0)
}
