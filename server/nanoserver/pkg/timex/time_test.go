package timex

import (
	"testing"
	"time"
)

func TestParseTimeWithFormat_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		value    string
		format   TimeFormat
		expected time.Time
		hasError bool
	}{
		{"2023-10-01 12:34:56", FormatDateTime, time.Date(2023, 10, 1, 12, 34, 56, 0, time.Local), false},
		{"2023-10-01", FormatDate, time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local), false},
		{"12:34:56", FormatTime, time.Date(0, 1, 1, 12, 34, 56, 0, time.Local), false},
		{"", FormatDateTime, time.Time{}, true},
		{"invalid", FormatDateTime, time.Time{}, true},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.value, func(t *testing.T) {
			t.Parallel()
			result, err := ParseTimeWithFormat(tt.value, tt.format)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v, got: %v", tt.hasError, err)
			}
			if !tt.hasError && !result.Equal(tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestParseTimeWithFormat_Int(t *testing.T) {
	t.Parallel()

	// Create test timestamps
	unixTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local)
	unixTimestamp := unixTime.Unix()
	unixMilliTimestamp := unixTime.UnixMilli()
	unixMicroTimestamp := unixTime.UnixMicro()
	unixNanoTimestamp := unixTime.UnixNano()

	tests := []struct {
		value    int64
		format   TimeFormat
		expected time.Time
		hasError bool
	}{
		{unixTimestamp, FormatUnix, unixTime, false},
		{unixMilliTimestamp, FormatUnixMilli, unixTime, false},
		{unixMicroTimestamp, FormatUnixMicro, unixTime, false},
		{unixNanoTimestamp, FormatUnixNano, unixTime, false},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()
			result, err := ParseTimeWithFormat(tt.value, tt.format)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v, got: %v", tt.hasError, err)
			}
			if !tt.hasError && !result.Equal(tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestParseTimeWithFormat_Float(t *testing.T) {
	t.Parallel()

	// Create test timestamps
	baseTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local)

	tests := []struct {
		value    float64
		format   TimeFormat
		expected time.Time
		hasError bool
	}{
		{float64(baseTime.Unix()) + 0.123, FormatUnix, baseTime.Add(123 * time.Millisecond), false},
		{float64(baseTime.UnixMilli()), FormatUnixMilli, baseTime, false},
		{float64(baseTime.UnixMicro()), FormatUnixMicro, baseTime, false},
		{float64(baseTime.UnixNano()), FormatUnixNano, baseTime, false},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()
			result, err := ParseTimeWithFormat(tt.value, tt.format)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v, got: %v", tt.hasError, err)
			}
			if !tt.hasError {
				// Allow a 1 microsecond difference
				diff := result.Sub(tt.expected)
				if diff < -time.Microsecond || diff > time.Microsecond {
					t.Errorf("expected: %v, got: %v, diff: %v", tt.expected, result, diff)
				}
			}
		})
	}
}

func TestParseTimeWithFormat_AutoDetect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		value    string
		expected time.Time
		hasError bool
	}{
		{"2023-10-01 12:34:56", time.Date(2023, 10, 1, 12, 34, 56, 0, time.Local), false},
		{"2023-10-01", time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local), false},
		{"12:34:56", time.Date(0, 1, 1, 12, 34, 56, 0, time.Local), false},
		{"2023/10/01", time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local), false},
		{"2023年10月01日", time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local), false},
		{"invalid", time.Time{}, true},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.value, func(t *testing.T) {
			t.Parallel()
			result, err := ParseTimeWithFormat(tt.value, "")
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v, got: %v", tt.hasError, err)
			}
			if !tt.hasError && !result.Equal(tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}
