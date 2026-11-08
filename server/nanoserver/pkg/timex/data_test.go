package timex

import (
	"fmt"
	"testing"
	"time"
)

func TestBeginOfWeek(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Monday",
			input:    time.Date(2021, 10, 4, 15, 20, 0, 123, time.FixedZone("CST", 8*3600)),
			expected: time.Date(2021, 10, 4, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		},
		{
			name:     "Wednesday",
			input:    time.Date(2021, 10, 6, 15, 20, 0, 123, time.FixedZone("CST", 8*3600)),
			expected: time.Date(2021, 10, 4, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		},
		{
			name:     "Sunday",
			input:    time.Date(2021, 10, 10, 15, 20, 0, 123, time.FixedZone("CST", 8*3600)),
			expected: time.Date(2021, 10, 4, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		},
		{
			name:     "Saturday",
			input:    time.Date(2021, 10, 9, 15, 20, 0, 123, time.FixedZone("CST", 8*3600)),
			expected: time.Date(2021, 10, 4, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BeginOfWeek(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("BeginOfWeek(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSameDay(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		{
			name:     "Same day",
			t1:       time.Date(2021, 10, 1, 15, 20, 0, 0, time.UTC),
			t2:       time.Date(2021, 10, 1, 10, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Different days",
			t1:       time.Date(2021, 10, 1, 23, 59, 59, 999999999, time.UTC),
			t2:       time.Date(2021, 10, 2, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Different months",
			t1:       time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Different years",
			t1:       time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Same day different time zones",
			t1:       time.Date(2021, 10, 1, 23, 59, 59, 999999999, time.FixedZone("CST", 8*3600)),
			t2:       time.Date(2021, 10, 1, 0, 0, 0, 0, time.FixedZone("PST", -8*3600)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSameDay(tt.t1, tt.t2)
			if result != tt.expected {
				t.Errorf("IsSameDay(%v, %v) = %v; want %v", tt.t1, tt.t2, result, tt.expected)
			}
		})
	}
}

func TestBeginOfYear(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Start of the year",
			input:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Middle of the year",
			input:    time.Date(2021, 6, 15, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "End of the year",
			input:    time.Date(2021, 12, 31, 23, 59, 59, 999999999, time.UTC),
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Leap year",
			input:    time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Different time zone",
			input:    time.Date(2021, 7, 1, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BeginOfYear(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("BeginOfYear(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBeginOfDay(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Start of the day",
			input:    time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Middle of the day",
			input:    time.Date(2021, 10, 1, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "End of the day",
			input:    time.Date(2021, 10, 1, 23, 59, 59, 999999999, time.UTC),
			expected: time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Different time zone",
			input:    time.Date(2021, 10, 1, 12, 0, 0, 0, time.FixedZone("CST", 8*3600)),
			expected: time.Date(2021, 10, 1, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		},
		{
			name:     "Leap year",
			input:    time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BeginOfDay(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("BeginOfDay(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEndOfDay(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Start of the day",
			input:    time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 10, 1, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:     "Middle of the day",
			input:    time.Date(2021, 10, 1, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 10, 1, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:     "End of the day",
			input:    time.Date(2021, 10, 1, 23, 59, 59, 999999999, time.UTC),
			expected: time.Date(2021, 10, 1, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:     "Different time zone",
			input:    time.Date(2021, 10, 1, 12, 0, 0, 0, time.FixedZone("CST", 8*3600)),
			expected: time.Date(2021, 10, 1, 23, 59, 59, 999999999, time.FixedZone("CST", 8*3600)),
		},
		{
			name:     "Leap year",
			input:    time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2020, 2, 29, 23, 59, 59, 999999999, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EndOfDay(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("EndOfDay(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBeginOfNextDay(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Start of the day",
			input:    time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 10, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Middle of the day",
			input:    time.Date(2021, 10, 1, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 10, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "End of the day",
			input:    time.Date(2021, 10, 1, 23, 59, 59, 999999999, time.UTC),
			expected: time.Date(2021, 10, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Different time zone",
			input:    time.Date(2021, 10, 1, 12, 0, 0, 0, time.FixedZone("CST", 8*3600)),
			expected: time.Date(2021, 10, 2, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		},
		{
			name:     "Leap year",
			input:    time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BeginOfNextDay(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("BeginOfNextDay(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBeginOfNextWeek(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Start of the week",
			input:    time.Date(2021, 10, 4, 0, 0, 0, 0, time.UTC),  // Monday
			expected: time.Date(2021, 10, 11, 0, 0, 0, 0, time.UTC), // Next Monday
		},
		{
			name:     "Middle of the week",
			input:    time.Date(2021, 10, 6, 12, 0, 0, 0, time.UTC), // Wednesday
			expected: time.Date(2021, 10, 11, 0, 0, 0, 0, time.UTC), // Next Monday
		},
		{
			name:     "End of the week",
			input:    time.Date(2021, 10, 10, 23, 59, 59, 999999999, time.UTC), // Sunday
			expected: time.Date(2021, 10, 11, 0, 0, 0, 0, time.UTC),            // Next Monday
		},
		{
			name:     "Different time zone",
			input:    time.Date(2021, 10, 6, 12, 0, 0, 0, time.FixedZone("CST", 8*3600)), // Wednesday in CST
			expected: time.Date(2021, 10, 11, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)), // Next Monday in CST
		},
		{
			name:     "Leap year",
			input:    time.Date(2020, 2, 24, 12, 0, 0, 0, time.UTC), // Monday in leap year
			expected: time.Date(2020, 3, 2, 0, 0, 0, 0, time.UTC),   // Next Monday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BeginOfNextWeek(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("BeginOfNextWeek(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAddDays(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		days     int
		expected time.Time
	}{
		{
			name:     "Start of the month",
			input:    time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
			days:     7,
			expected: time.Date(2021, 10, 7, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "Middle of the month",
			input:    time.Date(2021, 10, 15, 12, 0, 0, 0, time.UTC),
			days:     7,
			expected: time.Date(2021, 10, 21, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "End of the month",
			input:    time.Date(2021, 10, 31, 23, 59, 59, 999999999, time.UTC),
			days:     7,
			expected: time.Date(2021, 11, 6, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "Leap year",
			input:    time.Date(2020, 2, 24, 0, 0, 0, 0, time.UTC),
			days:     7,
			expected: time.Date(2020, 3, 1, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "Different time zone",
			input:    time.Date(2021, 10, 1, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
			days:     7,
			expected: time.Date(2021, 10, 7, 23, 59, 59, 0, time.FixedZone("CST", 8*3600)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddDays(tt.input, tt.days)
			if !result.Equal(tt.expected) {
				t.Errorf("AddDays(%v, %d) = %v; want %v", tt.input, tt.days, result, tt.expected)
			}
		})
	}
}

func TestDaysFromTime(t *testing.T) {
	tests := []struct {
		name    string
		input   time.Time
		want    int
		wantErr bool
	}{
		{
			name:    "today",
			input:   time.Now().UTC(),
			want:    1,
			wantErr: false,
		},
		{
			name:    "yesterday",
			input:   time.Now().UTC().AddDate(0, 0, -1),
			want:    2,
			wantErr: false,
		},
		{
			name:    "future",
			input:   time.Now().UTC().AddDate(0, 0, 1),
			want:    0,
			wantErr: true,
		},
		{
			name:    "different timezone",
			input:   time.Now().In(time.FixedZone("UTC+8", 8*60*60)),
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DaysFromTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DaysFromTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("DaysFromTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDayZeroTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		days     int
		expected time.Time
	}{
		{
			days:     0,
			expected: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		},
		{
			days:     1,
			expected: time.Date(time.Now().AddDate(0, 0, 1).Year(), time.Now().AddDate(0, 0, 1).Month(), time.Now().AddDate(0, 0, 1).Day(), 0, 0, 0, 0, time.Local),
		},
		{
			days:     -1,
			expected: time.Date(time.Now().AddDate(0, 0, -1).Year(), time.Now().AddDate(0, 0, -1).Month(), time.Now().AddDate(0, 0, -1).Day(), 0, 0, 0, 0, time.Local),
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(fmt.Sprintf("days=%d", tt.days), func(t *testing.T) {
			t.Parallel()
			result := GetDayZeroTime(tt.days)
			if !result.Equal(tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestIsToday(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected bool
	}{
		{
			name:     "today",
			input:    now,
			expected: true,
		},
		{
			name:     "today different hour",
			input:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local),
			expected: true,
		},
		{
			name:     "yesterday",
			input:    now.AddDate(0, 0, -1),
			expected: false,
		},
		{
			name:     "tomorrow",
			input:    now.AddDate(0, 0, 1),
			expected: false,
		},
		{
			name:     "last month",
			input:    now.AddDate(0, -1, 0),
			expected: false,
		},
		{
			name:     "last year same day",
			input:    now.AddDate(-1, 0, 0),
			expected: false,
		},
		{
			name:     "different timezone",
			input:    now.UTC(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsToday(tt.input)
			if got != tt.expected {
				t.Errorf("IsToday() = %v, want %v", got, tt.expected)
			}
		})
	}
}
