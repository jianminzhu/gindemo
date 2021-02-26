package dateUtil

import (
	"time"
)

var YMD string = "2006-01-02"
var YMDHMS string = "2006-01-02 15:04:05"

func ParseYmd(dateYmdStr string) time.Time {
	parse, _ := time.Parse(YMD, dateYmdStr)
	return parse
}
func Format2Ymd(date time.Time) string {
	return date.Format(YMD)
}

func AddDays(date time.Time, addDays int) time.Time {
	return date.AddDate(0, 0, addDays)
}

func DiffDays(endDate time.Time, startDate time.Time) int64 {
	return int64(endDate.Sub(startDate).Hours() / 24)
}
