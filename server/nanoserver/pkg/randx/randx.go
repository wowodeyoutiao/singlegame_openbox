package randx

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
)

// randUint32 generates a random uint32 less than max
func RandUint32(max uint32) (uint32, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return 0, fmt.Errorf("randUint32 err %v ", err)
	}
	return binary.BigEndian.Uint32(b) % max, nil
}

// randInt64 generates a random int64 less than max
func RandInt64(max int64) (int64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, fmt.Errorf("RandInt64 err %v ", err)
	}
	return n.Int64(), nil
}

// Integer is a constraint that permits any integer type.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// RandIntRange generates a random integer of type T in the closed interval [min, max].
func RandIntRange[T Integer](min, max T) (T, error) {
	if min > max {
		return 0, fmt.Errorf("min cannot be greater than max")
	}

	// Convert min and max to big.Int
	bigMin := new(big.Int)
	bigMax := new(big.Int)
	// Use string conversion to handle all integer types safely
	bigMin.SetString(fmt.Sprintf("%d", min), 10)
	bigMax.SetString(fmt.Sprintf("%d", max), 10)

	// Calculate range: max - min + 1
	bigRange := new(big.Int).Sub(bigMax, bigMin)
	bigRange.Add(bigRange, big.NewInt(1))

	// Generate random number in [0, bigRange)
	n, err := rand.Int(rand.Reader, bigRange)
	if err != nil {
		return 0, fmt.Errorf("RandIntRange err %v", err)
	}

	// Shift the result to [min, max]: result = n + min
	n.Add(n, bigMin)

	// Convert the result back to type T
	var result T
	// This is a simple way to convert back, assuming the value fits.
	// For production code, more robust error handling for conversion might be needed.
	_, err = fmt.Sscan(n.String(), &result)
	if err != nil {
		return 0, fmt.Errorf("failed to convert big.Int back to type T: %v", err)
	}

	return result, nil
}
