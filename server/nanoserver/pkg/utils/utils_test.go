package utils

import (
	"testing"
)

func TestConvertJSONToMap(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "valid JSON object",
			jsonStr:  `{"name":"John","age":30,"active":true}`,
			expected: map[string]interface{}{"name": "John", "age": float64(30), "active": true},
			wantErr:  false,
		},
		{
			name:     "empty JSON object",
			jsonStr:  `{}`,
			expected: map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "nested JSON object",
			jsonStr:  `{"user":{"name":"Jane","details":{"age":25}}}`,
			expected: map[string]interface{}{"user": map[string]interface{}{"name": "Jane", "details": map[string]interface{}{"age": float64(25)}}},
			wantErr:  false,
		},
		{
			name:     "empty string",
			jsonStr:  "",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid JSON",
			jsonStr:  `{"name":"John","age":}`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "JSON array",
			jsonStr:  `[1,2,3]`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "JSON string",
			jsonStr:  `"hello"`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "malformed JSON",
			jsonStr:  `{name:"John"}`,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertJSONToMap(tt.jsonStr)

			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertJSONToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(result) != len(tt.expected) {
					t.Errorf("ConvertJSONToMap() = %v, want %v", result, tt.expected)
					return
				}

				for key, expectedValue := range tt.expected {
					if actualValue, ok := result[key]; !ok {
						t.Errorf("ConvertJSONToMap() missing key %s", key)
					} else if !deepEqual(actualValue, expectedValue) {
						t.Errorf("ConvertJSONToMap() key %s = %v, want %v", key, actualValue, expectedValue)
					}
				}
			}
		})
	}
}

func TestRemoveElementInPlace(t *testing.T) {
	tests := []struct {
		name        string
		slice       []int
		index       int
		expected    []int
		wantErr     bool
		errContains string
	}{
		{
			name:     "remove first element",
			slice:    []int{1, 2, 3, 4, 5},
			index:    0,
			expected: []int{2, 3, 4, 5},
			wantErr:  false,
		},
		{
			name:     "remove middle element",
			slice:    []int{1, 2, 3, 4, 5},
			index:    2,
			expected: []int{1, 2, 4, 5},
			wantErr:  false,
		},
		{
			name:     "remove last element",
			slice:    []int{1, 2, 3, 4, 5},
			index:    4,
			expected: []int{1, 2, 3, 4},
			wantErr:  false,
		},
		{
			name:     "remove only element",
			slice:    []int{42},
			index:    0,
			expected: []int{},
			wantErr:  false,
		},
		{
			name:        "empty slice",
			slice:       []int{},
			index:       0,
			expected:    []int{},
			wantErr:     true,
			errContains: "index 0 out of range",
		},
		{
			name:        "negative index",
			slice:       []int{1, 2, 3},
			index:       -1,
			expected:    []int{1, 2, 3},
			wantErr:     true,
			errContains: "index -1 out of range",
		},
		{
			name:        "index too large",
			slice:       []int{1, 2, 3},
			index:       3,
			expected:    []int{1, 2, 3},
			wantErr:     true,
			errContains: "index 3 out of range",
		},
		{
			name:        "index way out of bounds",
			slice:       []int{1, 2},
			index:       10,
			expected:    []int{1, 2},
			wantErr:     true,
			errContains: "index 10 out of range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the original test data
			input := make([]int, len(tt.slice))
			copy(input, tt.slice)

			result, err := RemoveElementInPlace(input, tt.index)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveElementInPlace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err != nil && tt.errContains != "" {
					if err.Error() == "" || len(err.Error()) == 0 {
						t.Errorf("RemoveElementInPlace() error message is empty")
					}
				}
			} else {
				if len(result) != len(tt.expected) {
					t.Errorf("RemoveElementInPlace() result length = %d, want %d", len(result), len(tt.expected))
					return
				}

				for i, v := range tt.expected {
					if result[i] != v {
						t.Errorf("RemoveElementInPlace() result[%d] = %v, want %v", i, result[i], v)
					}
				}
			}
		})
	}
}

func TestRemoveElementInPlaceString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		index    int
		expected []string
		wantErr  bool
	}{
		{
			name:     "remove string from middle",
			slice:    []string{"apple", "banana", "cherry", "date"},
			index:    1,
			expected: []string{"apple", "cherry", "date"},
			wantErr:  false,
		},
		{
			name:     "remove last string",
			slice:    []string{"hello", "world"},
			index:    1,
			expected: []string{"hello"},
			wantErr:  false,
		},
		{
			name:    "empty string slice",
			slice:   []string{},
			index:   0,
			wantErr: true,
		},

		{
			name:     "remove only string",
			slice:    []string{"only"},
			index:    0,
			expected: []string{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]string, len(tt.slice))
			copy(input, tt.slice)

			result, err := RemoveElementInPlace(input, tt.index)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveElementInPlace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(result) != len(tt.expected) {
					t.Errorf("RemoveElementInPlace() result length = %d, want %d", len(result), len(tt.expected))
					return
				}

				for i, v := range tt.expected {
					if result[i] != v {
						t.Errorf("RemoveElementInPlace() result[%d] = %v, want %v", i, result[i], v)
					}
				}
			}
		})
	}
}

type testStruct struct {
	ID   int
	Name string
}

func TestRemoveElementInPlaceStruct(t *testing.T) {
	tests := []struct {
		name     string
		slice    []testStruct
		index    int
		expected []testStruct
		wantErr  bool
	}{
		{
			name: "remove struct from slice",
			slice: []testStruct{
				{ID: 1, Name: "first"},
				{ID: 2, Name: "second"},
				{ID: 3, Name: "third"},
			},
			index: 1,
			expected: []testStruct{
				{ID: 1, Name: "first"},
				{ID: 3, Name: "third"},
			},
			wantErr: false,
		},
		{
			name: "remove only struct",
			slice: []testStruct{
				{ID: 1, Name: "only"},
			},
			index:    0,
			expected: []testStruct{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]testStruct, len(tt.slice))
			copy(input, tt.slice)

			result, err := RemoveElementInPlace(input, tt.index)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveElementInPlace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(result) != len(tt.expected) {
					t.Errorf("RemoveElementInPlace() result length = %d, want %d", len(result), len(tt.expected))
					return
				}

				for i, v := range tt.expected {
					if result[i] != v {
						t.Errorf("RemoveElementInPlace() result[%d] = %v, want %v", i, result[i], v)
					}
				}
			}
		})
	}
}

func deepEqual(a, b interface{}) bool {
	switch va := a.(type) {
	case map[string]interface{}:
		if vb, ok := b.(map[string]interface{}); ok {
			if len(va) != len(vb) {
				return false
			}
			for k, v := range va {
				if !deepEqual(v, vb[k]) {
					return false
				}
			}
			return true
		}
		return false
	default:
		return a == b
	}
}
