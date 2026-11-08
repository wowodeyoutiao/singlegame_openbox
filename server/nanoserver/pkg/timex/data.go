package timex

import (
	"fmt"
	"time"
)

const (
	DaysInLeapYear   = 366 // 闰年天数
	DaysInCommonYear = 365 // 平年天数

	secondsInMinute = 60
	secondsInHour   = 3600
	secondsInDay    = 86400
	hoursInDay      = 24
)

// StartOfWeek 获取本周开始时间（周一）
// Deprecated: use BeginOfWeek instead
func StartOfWeek(t time.Time) time.Time {
	return BeginOfWeek(t)
}

// BeginOfWeek 获取本周开始时间（周一）
func BeginOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	return BeginOfDay(t.AddDate(0, 0, -weekday+1))
}

// EndOfWeek 获取本周结束时间（周日）
func EndOfWeek(t time.Time) time.Time {
	return BeginOfWeek(t).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// BeginOfDay 获取当天开始时间
func BeginOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取当天结束时间
func EndOfDay(t time.Time) time.Time {
	return BeginOfDay(t).AddDate(0, 0, 1).Add(-time.Nanosecond)
}

// BeginOfNextDay 第二天的开始时间
func BeginOfNextDay(t time.Time) time.Time {
	return BeginOfDay(t).AddDate(0, 0, 1)
}

// 获取下周一的开始时间
func BeginOfNextWeek(t time.Time) time.Time {
	return BeginOfWeek(t).AddDate(0, 0, 7)
}

// AddDays 获取指定天数后的时间，并将时间设置为当天的23:59:59
func AddDays(t time.Time, days int) time.Time {
	newDate := t.AddDate(0, 0, days-1)
	return time.Date(newDate.Year(), newDate.Month(), newDate.Day(), 23, 59, 59, 0, newDate.Location())
}

// BeginOfMonth 获取当月开始时间
func BeginOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 获取当月结束时间
func EndOfMonth(t time.Time) time.Time {
	return BeginOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// BeginOfYear 获取当年开始时间
func BeginOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 获取当年结束时间
func EndOfYear(t time.Time) time.Time {
	return BeginOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// IsLeapYear 判断是否为闰年
func IsLeapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}

// DaysInYear 获取指定年份的天数
func DaysInYear(year int) int {
	if IsLeapYear(year) {
		return DaysInLeapYear
	}
	return DaysInCommonYear
}

// DaysInMonth 获取指定年月的天数
func DaysInMonth(year int, month time.Month) int {
	switch month {
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if IsLeapYear(year) {
			return 29
		}
		return 28
	default:
		return 31
	}
}

// IsSameDay 判断两个时间是否为同一天
func IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameWeek 判断两个时间是否为同一周
func IsSameWeek(t1, t2 time.Time) bool {
	return BeginOfWeek(t1).Equal(BeginOfWeek(t2))
}

// IsSameMonth 判断两个时间是否为同一月
func IsSameMonth(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// IsSameYear 判断两个时间是否为同一年
func IsSameYear(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year()
}

// Between 判断时间是否在指定范围内（包含边界）
func Between(t, start, end time.Time) bool {
	return !t.Before(start) && !t.After(end)
}

// SubDays 计算两个时间相差的天数
func SubDays(t1, t2 time.Time) int {
	duration := t1.Sub(t2)
	return int(duration.Hours() / hoursInDay)
}

// Format 格式化时间为字符串，使用预定义的格式
func Format(t time.Time, format TimeFormat) string {
	return t.Format(string(format))
}

// AddWorkDays 添加工作日（跳过周末）
func AddWorkDays(t time.Time, days int) time.Time {
	for days > 0 {
		t = t.AddDate(0, 0, 1)
		if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday {
			days--
		}
	}
	return t
}

// IsWeekend 判断是否为周末
func IsWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

// IsWorkday 判断是否为工作日
func IsWorkday(t time.Time) bool {
	return !IsWeekend(t)
}

// BeginOfHour 获取小时开始时间
func BeginOfHour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

// EndOfHour 获取小时结束时间
func EndOfHour(t time.Time) time.Time {
	return BeginOfHour(t).Add(time.Hour).Add(-time.Nanosecond)
}

// BeginOfMinute 获取分钟开始时间
func BeginOfMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

// EndOfMinute 获取分钟结束时间
func EndOfMinute(t time.Time) time.Time {
	return BeginOfMinute(t).Add(time.Minute).Add(-time.Nanosecond)
}

// RoundToDay 将时间舍入到天
func RoundToDay(t time.Time) time.Time {
	return BeginOfDay(t)
}

// RoundToHour 将时间舍入到小时
func RoundToHour(t time.Time) time.Time {
	return BeginOfHour(t)
}

// RoundToMinute 将时间舍入到分钟
func RoundToMinute(t time.Time) time.Time {
	return BeginOfMinute(t)
}

// DaysFromTime returns the number of days from the given time to now.
// Returns 1 if the given time is today, 2 if yesterday, and so on.
func DaysFromTime(t time.Time) (int, error) {
	now := time.Now()
	targetTime := t.Local()

	if now.Before(targetTime) {
		return 0, fmt.Errorf("input time is in the future")
	}

	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	targetDate := time.Date(targetTime.Year(), targetTime.Month(), targetTime.Day(), 0, 0, 0, 0, time.Local)

	days := int(nowDate.Sub(targetDate).Hours()/24) + 1
	return days, nil
}

// IsToday checks if the given Unix timestamp (in seconds) is from today in Shanghai timezone
func IsToday(target time.Time) bool {
	now := time.Now()
	t := target.Local()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// GetDayZeroTime 获取指定日期偏移的零点时间
// days: 0=今天, 1=明天, -1=昨天
func GetDayZeroTime(days int) time.Time {
	now := time.Now()
	year, month, day := now.Date()
	zeroTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	if days != 0 {
		zeroTime = zeroTime.AddDate(0, 0, days)
	}
	return zeroTime
}
