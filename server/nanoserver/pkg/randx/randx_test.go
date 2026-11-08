package randx

import (

	//"math/rand/v2"
	"sync"
	"testing"
)

func TestRandUint32(t *testing.T) {
	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				max := uint32(100)
				num, err := RandUint32(max)
				if err != nil {
					t.Errorf("RandUint32 error: %v", err)
				}
				if num >= max {
					t.Errorf("Expected number less than %d, but got %d", max, num)
				}

				//t.Logf("sdfsdfsdfdsfd    %d", num)
			}
		}()
	}

	wg.Wait()
}

func TestRandInt64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		max := int64(100)
		num, err := RandInt64(max)
		if err != nil {
			t.Errorf("RandInt64 error: %v", err)
		}
		if num >= max {
			t.Errorf("Expected number less than %d, but got %d", max, num)
		}
	}
}

// testRange is a helper function to test a specific range for a given integer type.
func testRange[T Integer](t *testing.T, min, max T) {
	t.Helper()
	for i := 0; i < 100; i++ {
		n, err := RandIntRange(min, max)
		if err != nil {
			t.Errorf("RandIntRange(%v, %v) returned an error: %v", min, max, err)
			return
		}
		if n < min || n > max {
			t.Errorf("RandIntRange(%v, %v) returned %v, which is outside the valid range", min, max, n)
		}
	}
}

func TestRandIntRange(t *testing.T) {
	t.Run("int range", func(t *testing.T) {
		testRange(t, -10, 10)
	})

	t.Run("int8 range", func(t *testing.T) {
		testRange[int8](t, -50, 50)
	})

	t.Run("int16 range", func(t *testing.T) {
		testRange[int16](t, -1000, 1000)
	})

	t.Run("int32 range", func(t *testing.T) {
		testRange[int32](t, -100000, 100000)
	})

	t.Run("int64 range", func(t *testing.T) {
		testRange[int64](t, -1000000, 1000000)
	})

	t.Run("uint range", func(t *testing.T) {
		testRange[uint](t, 100, 200)
	})

	t.Run("uint8 range", func(t *testing.T) {
		testRange[uint8](t, 10, 20)
	})

	t.Run("uint16 range", func(t *testing.T) {
		testRange[uint16](t, 1000, 2000)
	})

	t.Run("uint32 range", func(t *testing.T) {
		testRange[uint32](t, 100000, 200000)
	})

	t.Run("uint64 range", func(t *testing.T) {
		testRange[uint64](t, 1000000, 2000000)
	})

	t.Run("single value range", func(t *testing.T) {
		min, max := 42, 42
		n, err := RandIntRange(min, max)
		if err != nil {
			t.Errorf("RandIntRange(%d, %d) returned an error: %v", min, max, err)
		}
		if n != min {
			t.Errorf("Expected %d, but got %d", min, n)
		}
	})

	t.Run("invalid range", func(t *testing.T) {
		min, max := 100, 10
		_, err := RandIntRange(min, max)
		if err == nil {
			t.Errorf("RandIntRange(%d, %d) should have returned an error for min > max, but it did not", min, max)
		}
		expectedErr := "min cannot be greater than max"
		if err.Error() != expectedErr {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("zero range", func(t *testing.T) {
		testRange(t, 0, 0)
	})

	t.Run("negative range", func(t *testing.T) {
		testRange(t, -100, -10)
	})
}
