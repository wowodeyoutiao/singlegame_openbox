package timex

import (
	"fmt"
	"math"
	"time"
)

const (
	// Timestamp format digit lengths
	unixSecondsDigits = 10 // 1683729075
	unixMillisDigits  = 13 // 1683729075000
	unixMicrosDigits  = 16 // 1683729075000000
	unixNanosDigits   = 19 // 1683729075000000000
)

// TimeFormat defines different time formats
type TimeFormat string

const (
	// Using standard formats from the time package
	FormatANSIC       = TimeFormat(time.ANSIC)       // "Mon Jan _2 15:04:05 2006"
	FormatUnixDate    = TimeFormat(time.UnixDate)    // "Mon Jan _2 15:04:05 MST 2006"
	FormatRubyDate    = TimeFormat(time.RubyDate)    // "Mon Jan 02 15:04:05 -0700 2006"
	FormatRFC822      = TimeFormat(time.RFC822)      // "02 Jan 06 15:04 MST"
	FormatRFC822Z     = TimeFormat(time.RFC822Z)     // "02 Jan 06 15:04 -0700"
	FormatRFC850      = TimeFormat(time.RFC850)      // "Monday, 02-Jan-06 15:04:05 MST"
	FormatRFC1123     = TimeFormat(time.RFC1123)     // "Mon, 02 Jan 2006 15:04:05 MST"
	FormatRFC1123Z    = TimeFormat(time.RFC1123Z)    // "Mon, 02 Jan 2006 15:04:05 -0700"
	FormatRFC3339     = TimeFormat(time.RFC3339)     // "2006-01-02T15:04:05Z07:00"
	FormatRFC3339Nano = TimeFormat(time.RFC3339Nano) // "2006-01-02T15:04:05.999999999Z07:00"
	FormatKitchen     = TimeFormat(time.Kitchen)     // "3:04PM"

	// Common timestamp formats
	FormatStamp      = TimeFormat(time.Stamp)      // "Jan _2 15:04:05"
	FormatStampMilli = TimeFormat(time.StampMilli) // "Jan _2 15:04:05.000"
	FormatStampMicro = TimeFormat(time.StampMicro) // "Jan _2 15:04:05.000000"
	FormatStampNano  = TimeFormat(time.StampNano)  // "Jan _2 15:04:05.000000000"

	// Common date and time formats
	FormatDateTime = TimeFormat(time.DateTime) // "2006-01-02 15:04:05"
	FormatDate     = TimeFormat(time.DateOnly) // "2006-01-02"
	FormatTime     = TimeFormat(time.TimeOnly) // "15:04:05"

	// Custom extended formats
	FormatDateSlash      = "2006/01/02"
	FormatDateChinese    = "2006年01月02日"
	FormatUnix           = "unix"                    // 示例: 1712323200 (秒级时间戳)
	FormatUnixMilli      = "unixmilli"               // 示例: 1712323200000 (毫秒级时间戳)
	FormatUnixMicro      = "unixmicro"               // 示例: 1712323200000000 (微秒级时间戳)
	FormatUnixNano       = "unixnano"                // 示例: 1712323200000000000 (纳秒级时间戳)
	FormatDateTimeMillis = "2006-01-02 15:04:05.000" // yyyy-mm-dd hh:mm:ss.sss
)

// DateValue is used to constrain types that can be converted to time
type DateValue interface {
	~string | ~int64 | ~int32 | ~uint32 | ~uint64 | ~float64
}

// ParseTimeWithFormat converts different types of date values to time.Time using the system timezone
func ParseTimeWithFormat[T DateValue](value T, format TimeFormat) (time.Time, error) {
	switch v := any(value).(type) {
	case string:
		return parseStringTime(v, format)
	case int64, int32, uint32, uint64:
		return parseIntTime(v, format)
	case float64:
		return parseFloatTime(v, format)
	default:
		return time.Time{}, fmt.Errorf("unsupported type: %T", value)
	}
}

// parseStringTime parses a string time value using the system timezone
func parseStringTime(value string, format TimeFormat) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("empty time string")
	}

	// If no format is specified, try to auto-detect
	if format == "" {
		return parseAutoDetectFormat(value)
	}

	// Parse using the specified format and system timezone
	t, err := time.ParseInLocation(string(format), value, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse time error: %w", err)
	}
	return t, nil
}

func parseIntTime(value interface{}, format TimeFormat) (time.Time, error) {
	var timestamp int64
	switch v := value.(type) {
	case int64:
		timestamp = v
	case int32:
		timestamp = int64(v)
	case uint32:
		timestamp = int64(v)
	case uint64:
		//nolint:gosec
		timestamp = int64(v)
	default:
		return time.Time{}, fmt.Errorf("unsupported integer type: %T", value)
	}

	// If format is specified, use it directly
	if format != "" {
		var t time.Time
		switch format {
		case FormatUnix:
			t = time.Unix(timestamp, 0)
		case FormatUnixMilli:
			t = time.UnixMilli(timestamp)
		case FormatUnixMicro:
			t = time.UnixMicro(timestamp)
		case FormatUnixNano:
			t = time.Unix(0, timestamp)
		default:
			// Default to seconds
			t = time.Unix(timestamp, 0)
		}
		return t.In(time.Local), nil
	}

	// Auto-detect format based on number of digits
	digits := int(math.Log10(float64(timestamp))) + 1
	var t time.Time

	switch {
	case digits <= unixSecondsDigits:
		t = time.Unix(timestamp, 0)
	case digits <= unixMillisDigits:
		t = time.UnixMilli(timestamp)
	case digits <= unixMicrosDigits:
		t = time.UnixMicro(timestamp)
	case digits <= unixNanosDigits:
		t = time.Unix(0, timestamp)
	default:
		return time.Time{}, fmt.Errorf("timestamp digit count %d exceeds nanosecond precision", digits)
	}

	return t.In(time.Local), nil
}

// parseFloatTime parses a floating-point timestamp and converts it to the system timezone
//
//nolint:exhaustive
func parseFloatTime(value float64, format TimeFormat) (time.Time, error) {
	var t time.Time
	switch format {
	case FormatUnix:
		sec := int64(value)
		nsec := int64((value - float64(sec)) * 1e9)
		t = time.Unix(sec, nsec)
	case FormatUnixMilli:
		t = time.UnixMilli(int64(value))
	case FormatUnixMicro:
		t = time.UnixMicro(int64(value))
	case FormatUnixNano:
		t = time.Unix(0, int64(value))
	default:
		// Default to seconds
		sec := int64(value)
		nsec := int64((value - float64(sec)) * 1e9)
		t = time.Unix(sec, nsec)
	}
	return t.In(time.Local), nil
}

// parseAutoDetectFormat auto-detects the time format using the system timezone
func parseAutoDetectFormat(value string) (time.Time, error) {
	formats := []TimeFormat{
		FormatDateTime, // Most common format first
		FormatDate,
		FormatTime,
		FormatRFC3339,
		FormatRFC3339Nano,
		FormatANSIC,
		FormatUnixDate,
		FormatRubyDate,
		FormatRFC822,
		FormatRFC822Z,
		FormatRFC850,
		FormatRFC1123,
		FormatRFC1123Z,
		FormatStamp,
		FormatStampMilli,
		FormatStampMicro,
		FormatStampNano,
		FormatDateSlash,
		FormatDateChinese,
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(string(format), value, time.Local); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to detect time format for: %s", value)
}
