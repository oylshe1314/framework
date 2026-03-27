package util

import "time"

const (
	MinuteSeconds = 60
	HourSeconds   = 3600
	DaySeconds    = 86400

	UTC8TimezoneOffset = HourSeconds * 8
)

func UTC8() *time.Location {
	return time.FixedZone("UTC+8", UTC8TimezoneOffset)
}

func Now() time.Time {
	return time.Now().In(UTC8())
}

func NowUnix() int64 {
	return Now().Unix()
}

func NowMilli() int64 {
	return Now().UnixMilli()
}

func NowMicro() int64 {
	return Now().UnixMicro()
}

func DateOnly() string {
	return Now().Format(time.DateOnly)
}

func TimeOnly() string {
	return Now().Format(time.TimeOnly)
}

func Datetime() string {
	return Now().Format(time.DateTime)
}

func TimeFormat(layout string) string {
	return Now().Format(layout)
}

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

func FromUnix(unix int64) time.Time {
	return time.Unix(unix, 0).In(UTC8())
}

func YMD() (year, month, day int) {
	var now = time.Now().In(UTC8())
	return now.Year(), int(now.Month()), now.Day()
}

func TodayPassedSeconds() int64 {
	return (NowUnix() + UTC8TimezoneOffset) % DaySeconds
}

func TodayBeginTime() int64 {
	var now = NowUnix()
	return now - (now+UTC8TimezoneOffset)%DaySeconds
}

func TomorrowBeginTime() int64 {
	return TodayBeginTime() + DaySeconds
}

func WeekBeginTime(weekBeginDay int64) int64 {
	var now = Now()
	var weekday = int64(now.Weekday())
	weekday -= weekBeginDay
	if weekday < 0 {
		weekday = 6
	}

	var nt = now.Unix()
	return nt - (nt+UTC8TimezoneOffset)%DaySeconds - DaySeconds*weekday
}

func NextWeekBeginTime(weekBeginDay int64) int64 {
	var now = Now()
	var weekday = int64(now.Weekday())
	weekday -= weekBeginDay
	if weekday < 0 {
		weekday = 6
	}

	var nt = now.Unix()
	return nt - (nt+UTC8TimezoneOffset)%DaySeconds + DaySeconds*(7-weekday)
}
