package randx

import (
	"crypto/rand"
	"math/big"
	"sort"
)

// Picker represents a weighted random picker.
type Picker[T any] struct {
	items       []T
	prefixSums  []int
	totalWeight int
}

// New creates a new Picker with the given items and weight function.
func New[T any](items []T, weightFunc func(T) int) *Picker[T] {
	p := &Picker[T]{
		items:      items,
		prefixSums: make([]int, len(items)),
	}

	sum := 0
	for i, item := range items {
		weight := weightFunc(item)
		sum += weight
		p.prefixSums[i] = sum
	}
	p.totalWeight = sum

	return p
}

// Pick returns a randomly selected item based on weights.
func (p *Picker[T]) Pick() (T, error) {
	if len(p.items) == 0 {
		var zero T
		return zero, &ErrEmptyPicker{}
	}

	// Generate a random number between 1 and totalWeight
	n, err := rand.Int(rand.Reader, big.NewInt(int64(p.totalWeight)))
	if err != nil {
		var zero T
		return zero, err
	}
	x := int(n.Int64()) + 1

	index := sort.SearchInts(p.prefixSums, x)
	return p.items[index], nil
}

// ErrEmptyPicker is returned when trying to pick from an empty picker.
type ErrEmptyPicker struct{}

func (e *ErrEmptyPicker) Error() string {
	return "picker is empty"
}
