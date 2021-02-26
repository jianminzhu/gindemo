package util

import (
	"math"
	"strconv"
	"time"
)

const (
	TplYMD    string = "2006-01-02"
	TplYMDHMS string = "2006-01-02 15:04:05"
)

func Ymd(timeDate time.Time) string {
	return timeDate.Format(TplYMD)
}
func Int64ToInt(data int64) int {
	itoa, _ := strconv.Atoi(strconv.FormatInt(data, 10))
	return itoa
}

func GetDataStrAdded(dateStr string, days int64) string {
	a, _ := time.Parse(TplYMD, dateStr)
	return a.AddDate(0, 0, Int64ToInt(days)).Format(TplYMD)
}
func DataStrAddDays(dateStr string, days int) time.Time {
	a, _ := time.Parse(TplYMD, dateStr)
	return TimeAddDays(a, days)
}
func TimeAddDays(timeDate time.Time, days int) time.Time {
	return timeDate.AddDate(0, 0, days)
}

func GetDaysDuration(startDateStr string, endDateStr string) int64 {
	a, _ := time.Parse(TplYMD, startDateStr)
	b, _ := time.Parse(TplYMD, endDateStr)
	d := a.Sub(b)
	return Abs(Wrap(d.Hours()/24, 0))
}

func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

//将float64转成精确的int64
func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

//将int64恢复成正常的float64
func Unwrap(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}

//精准float64
func WrapToFloat64(num float64, retain int) float64 {
	return num * math.Pow10(retain)
}

//精准int64
func UnwrapToInt64(num int64, retain int) int64 {
	return int64(Unwrap(num, retain))
}
