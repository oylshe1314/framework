package util

import (
	"time"
)

const DayTotalSeconds int64 = 86400
const WeekTotalSeconds int64 = DayTotalSeconds * 7
const UTC8TimezoneOffset int64 = 3600 * 8

func UTC8() *time.Location {
	return time.FixedZone("UTC+8", int(UTC8TimezoneOffset))
}

// Now 取当前时间时不要用time包的Now()函数，用这个，防止服务器时区导致一些问题，取时间戳用下面的Unix()函数
func Now() time.Time {
	return time.Now().In(UTC8())
}

// Unix 秒时间戳没有时区
func Unix() int64 {
	return time.Now().Unix()
}

// UnixMilli 毫秒时间戳没有时区
func UnixMilli() int64 {
	return time.Now().UnixMilli()
}

// UnixMicro 微秒时间戳没有时区
func UnixMicro() int64 {
	return time.Now().UnixMicro()
}

// UnixNano 纳秒时间戳没有时区
func UnixNano() int64 {
	return time.Now().UnixNano()
}

// DateOnly  仅返回日期，按UTC8时区处理
func DateOnly() string {
	return Now().Format(time.DateOnly)
}

// TimeOnly  仅返回时间，按UTC8时区处理
func TimeOnly() string {
	return Now().Format(time.TimeOnly)
}

// Datetime  返回日期和时间，按UTC8时区处理
func Datetime() string {
	return Now().Format(time.DateTime)
}

// TimeFormat  按给定格式返回时间，按UTC8时区处理
func TimeFormat(layout string) string {
	return Now().Format(layout)
}

// ParseUnix 解析为秒时间戳，按UTC8时区处理
func ParseUnix(layout, value string) (int64, error) {
	if value == "" {
		return 0, nil
	}

	var t, err = time.ParseInLocation(layout, value, UTC8())
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

// YMD 返回整数的年、月、日，按UTC8时区处理
func YMD() (year, month, day int) {
	var now = time.Now().In(UTC8())
	return now.Year(), int(now.Month()), now.Day()
}

// TodayPassedSeconds 返回当天过了多少秒，按UTC8时区处理
func TodayPassedSeconds() int64 {
	return (Unix() + UTC8TimezoneOffset) % DayTotalSeconds
}

// TodayBeginTime 返回当天零点时的时间戳，按UTC8时区处理
func TodayBeginTime() int64 {
	var now = Unix()
	return now - (now+UTC8TimezoneOffset)%DayTotalSeconds
}

// WeekBeginTime 返回本周一零点时的时间戳，按UTC8时区处理
func WeekBeginTime() int64 {
	var now = Now()
	var weekday = int64(now.Weekday())
	weekday -= 1
	if weekday < 0 {
		weekday = 6
	}

	var nt = now.Unix()
	return nt - (nt+UTC8TimezoneOffset)%DayTotalSeconds - DayTotalSeconds*weekday
}
