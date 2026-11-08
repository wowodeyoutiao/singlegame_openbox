package utils

import (
	"strconv"
	"sync"
	"testing"
)

func TestSafeMap_NewAndBasicOps(t *testing.T) {
	t.Parallel()

	sm := NewSafeMap[string, int]()
	if sm == nil {
		t.Fatal("expected non-nil SafeMap")
	}
	if sm.Len() != 0 {
		t.Fatalf("expected Len=0, got=%d", sm.Len())
	}

	// Get missing
	if v, ok := sm.Get("missing"); ok || v != 0 {
		t.Fatalf("expected missing key => (0,false), got (%d,%v)", v, ok)
	}
	// GetOrDefault missing
	if v := sm.GetOrDefault("missing", 42); v != 42 {
		t.Fatalf("expected default 42, got %d", v)
	}
	// Set and Get
	sm.Set("a", 1)
	sm.Set("b", 2)
	if !sm.Has("a") || !sm.Has("b") {
		t.Fatal("expected keys to exist")
	}
	if v, ok := sm.Get("a"); !ok || v != 1 {
		t.Fatalf("expected a=1, got v=%d ok=%v", v, ok)
	}
	if sm.Len() != 2 {
		t.Fatalf("expected Len=2, got=%d", sm.Len())
	}
	// Delete existing
	sm.Delete("a")
	if sm.Has("a") {
		t.Fatal("expected 'a' deleted")
	}
	// Delete non-existing should not panic
	sm.Delete("zzz")
	if sm.Len() != 1 {
		t.Fatalf("expected Len=1, got=%d", sm.Len())
	}
}

func TestSafeMap_NewWithCapacity_ZeroAndNegative(t *testing.T) {
	t.Parallel()

	// Zero capacity works
	sm0 := NewSafeMapWithCapacity[int, string](0)
	if sm0 == nil || sm0.Len() != 0 {
		t.Fatal("expected zero-capacity SafeMap to be usable")
	}

	// Negative capacity should panic (makemap with negative hint panics)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative capacity")
		}
	}()
	_ = NewSafeMapWithCapacity[int, string](-1)
}

func TestSafeMap_SetIfAbsent(t *testing.T) {
	t.Parallel()

	sm := NewSafeMap[string, string]()
	if !sm.SetIfAbsent("x", "v1") {
		t.Fatal("expected SetIfAbsent to succeed on new key")
	}
	if sm.SetIfAbsent("x", "v2") {
		t.Fatal("expected SetIfAbsent to fail on existing key")
	}
	v, ok := sm.Get("x")
	if !ok || v != "v1" {
		t.Fatalf("expected x='v1', got v=%q ok=%v", v, ok)
	}
}

func TestSafeMap_SetIfAbsent_ConcurrentSingleWinner(t *testing.T) {
	t.Parallel()

	sm := NewSafeMap[string, int]()
	const goroutines = 16

	var wg sync.WaitGroup
	var winners int64
	results := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			ok := sm.SetIfAbsent("key", val)
			results <- ok
		}(i)
	}
	wg.Wait()
	close(results)

	for ok := range results {
		if ok {
			winners++
		}
	}
	if winners != 1 {
		t.Fatalf("expected exactly 1 SetIfAbsent winner, got %d", winners)
	}
	if !sm.Has("key") {
		t.Fatal("expected key to exist after SetIfAbsent race")
	}
}

func TestSafeMap_KeysValuesContent(t *testing.T) {
	t.Parallel()

	sm := NewSafeMap[int, string]()
	const n = 10
	exp := make(map[int]string, n)
	for i := 0; i < n; i++ {
		k := i * 2
		v := "v" + strconv.Itoa(k)
		exp[k] = v
		sm.Set(k, v)
	}

	keys := sm.Keys()
	if len(keys) != n {
		t.Fatalf("expected %d keys, got %d", n, len(keys))
	}
	for _, k := range keys {
		if _, ok := exp[k]; !ok {
			t.Fatalf("unexpected key %d in Keys()", k)
		}
	}

	values := sm.Values()
	if len(values) != n {
		t.Fatalf("expected %d values, got %d", n, len(values))
	}
	// Quick membership check by building frequency map
	freq := map[string]int{}
	for _, v := range values {
		freq[v]++
	}
	for _, v := range exp {
		if freq[v] != 1 {
			t.Fatalf("expected value %q exactly once, got %d", v, freq[v])
		}
	}
}

func TestSafeMap_RangeAndEarlyStop(t *testing.T) {
	t.Parallel()

	sm := NewSafeMap[int, int]()
	for i := 0; i < 50; i++ {
		sm.Set(i, i*i)
	}

	// Full range count
	count := 0
	sm.Range(func(k, v int) bool {
		count++
		return true
	})
	if count != 50 {
		t.Fatalf("expected to range 50 items, got %d", count)
	}

	// Early stop after 5 items
	count = 0
	sm.Range(func(k, v int) bool {
		count++
		return count < 5
	})
	if count != 5 {
		t.Fatalf("expected early stop at 5, got %d", count)
	}
}

func TestSafeMap_GetAndDelete(t *testing.T) {
	t.Parallel()

	sm := NewSafeMap[string, string]()
	sm.Set("k", "v")
	v, ok := sm.GetAndDelete("k")
	if !ok || v != "v" {
		t.Fatalf("expected GetAndDelete to return (v,true), got (%q,%v)", v, ok)
	}
	if sm.Has("k") {
		t.Fatal("expected key removed after GetAndDelete")
	}
	// Missing key returns zero, false
	v, ok = sm.GetAndDelete("missing")
	if ok || v != "" {
		t.Fatalf("expected missing => (\"\",false), got (%q,%v)", v, ok)
	}
}

func TestSafeMap_CopyIsSnapshot(t *testing.T) {
	t.Parallel()

	sm := NewSafeMap[int, string]()
	for i := 0; i < 5; i++ {
		sm.Set(i, "v"+strconv.Itoa(i))
	}
	cp := sm.Copy()
	if len(cp) != 5 {
		t.Fatalf("expected copy len=5, got %d", len(cp))
	}
	// Mutate original, copy should not change
	sm.Set(100, "x")
	sm.Delete(0)
	if len(sm.Copy()) != 5 { // current map has keys: 1..4 and 100
		t.Fatalf("sanity: expected live map len=5")
	}
	if len(cp) != 5 {
		t.Fatalf("expected snapshot copy to remain 5, got %d", len(cp))
	}
	if _, ok := cp[100]; ok {
		t.Fatal("unexpected key 100 in snapshot copy")
	}
	if _, ok := cp[0]; !ok {
		t.Fatal("expected key 0 to remain in snapshot copy")
	}
}

func TestSafeMap_ClearAndReuse(t *testing.T) {
	t.Parallel()

	sm := NewSafeMapWithCapacity[string, int](5)
	sm.Set("a", 1)
	sm.Set("b", 2)
	if sm.Len() != 2 {
		t.Fatalf("expected Len=2, got=%d", sm.Len())
	}

	sm.Clear()
	if sm.Len() != 0 {
		t.Fatalf("expected Len=0 after Clear, got=%d", sm.Len())
	}
	if _, ok := sm.Get("a"); ok {
		t.Fatal("expected key 'a' to be cleared")
	}

	// Reuse after clear
	sm.Set("c", 3)
	v, ok := sm.Get("c")
	if !ok || v != 3 {
		t.Fatalf("expected key 'c' => 3, got v=%d ok=%v", v, ok)
	}
}

func TestSafeMap_GrowthBeyondHint(t *testing.T) {
	t.Parallel()

	const hint = 4
	const total = 3 * hint
	sm := NewSafeMapWithCapacity[int, int](hint)

	for i := 0; i < total; i++ {
		sm.Set(i, i*i)
	}

	if got := sm.Len(); got != total {
		t.Fatalf("expected Len=%d, got=%d", total, got)
	}

	// Validate few entries
	for _, i := range []int{0, hint - 1, hint, total - 1} {
		v, ok := sm.Get(i)
		if !ok || v != i*i {
			t.Fatalf("expected key %d => %d, got v=%d ok=%v", i, i*i, v, ok)
		}
	}
}

func TestSafeMap_ConcurrentReadWriteSmoke(t *testing.T) {
	t.Parallel()

	const (
		hint        = 64
		writers     = 8
		readsPerKey = 10
		keysPerG    = 100
	)

	sm := NewSafeMapWithCapacity[int, string](hint)

	var wg sync.WaitGroup

	// Writers
	for g := 0; g < writers; g++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := 0; i < keysPerG; i++ {
				k := base*100000 + i
				sm.Set(k, "v"+strconv.Itoa(k))
			}
		}(g)
	}

	// Readers
	for g := 0; g < writers; g++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			start := base * 100000
			for i := 0; i < keysPerG; i++ {
				k := start + i
				for r := 0; r < readsPerKey; r++ {
					_, _ = sm.Get(k)
					_ = sm.Has(k)
				}
			}
		}(g)
	}

	wg.Wait()

	expected := writers * keysPerG
	if got := sm.Len(); got != expected {
		t.Fatalf("expected Len=%d after concurrent ops, got=%d", expected, got)
	}

	// Range and Copy consistency
	seen := 0
	sm.Range(func(k int, v string) bool {
		seen++
		return true
	})
	if seen != expected {
		t.Fatalf("expected Range to see %d items, saw %d", expected, seen)
	}
	copied := sm.Copy()
	if len(copied) != expected {
		t.Fatalf("expected Copy to have %d items, got %d", expected, len(copied))
	}
}

func TestSafeMap_RangeSafe_SnapshotAndMutation(t *testing.T) {
	t.Parallel()
	sm := NewSafeMap[string, int]()
	sm.Set("a", 1)
	sm.Set("b", 2)
	sm.Set("c", 3)

	seen := make(map[string]int, 4)

	sm.RangeSafe(func(k string, v int) bool {
		// Mutations during iteration should not deadlock and should not affect the current snapshot
		sm.Set("d", 4)
		sm.Delete("c")

		// Reentrant reads should work fine
		if !sm.Has("a") {
			t.Errorf("expected 'a' to exist during RangeSafe callback")
		}
		if got := sm.GetOrDefault("b", -1); got != 2 {
			t.Errorf("expected 'b' to be 2 via GetOrDefault, got %d", got)
		}

		seen[k] = v
		return true
	})

	// Post-iteration mutations must be visible on the live map
	if !sm.Has("d") {
		t.Fatal("expected key 'd' to be added during RangeSafe")
	}
	if sm.Has("c") {
		t.Fatal("expected key 'c' to be deleted during RangeSafe")
	}

	// Snapshot semantics: the iteration should have seen only the original keys
	if len(seen) != 3 {
		t.Fatalf("expected to see 3 items from snapshot, saw %d", len(seen))
	}
	for _, k := range []string{"a", "b", "c"} {
		if _, ok := seen[k]; !ok {
			t.Fatalf("expected to see key %q from snapshot", k)
		}
	}
	if _, ok := seen["d"]; ok {
		t.Fatal("did not expect to see key 'd' added during iteration")
	}
}

func TestSafeMap_RangeSafe_EarlyStop(t *testing.T) {
	t.Parallel()
	sm := NewSafeMap[int, int]()
	for i := 0; i < 100; i++ {
		sm.Set(i, i*i)
	}

	count := 0
	sm.RangeSafe(func(k, v int) bool {
		count++
		if count == 10 {
			// Mutate map on early stop to ensure no deadlock and independence from current snapshot
			sm.Set(1000, 1)
			sm.Delete(0)
			return false
		}
		return true
	})

	if count != 10 {
		t.Fatalf("expected early stop at 10, got %d", count)
	}
	if !sm.Has(1000) {
		t.Fatal("expected key 1000 to be added during RangeSafe early stop")
	}
	if sm.Has(0) {
		t.Fatal("expected key 0 to be deleted during RangeSafe early stop")
	}
}
