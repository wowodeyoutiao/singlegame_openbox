package utils

import (
	"sync"
)

// SafeMap 支持泛型的并发安全Map
type SafeMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewSafeMap 创建一个新的SafeMap实例
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		m: make(map[K]V),
	}
}

// NewSafeMapWithCapacity 创建具有指定初始容量的SafeMap
func NewSafeMapWithCapacity[K comparable, V any](capacity int) *SafeMap[K, V] {
	if capacity < 0 {
		panic("capacity must be non-negative")
	}
	return &SafeMap[K, V]{
		m: make(map[K]V, capacity),
	}
}

// Set 设置键值对
func (sm *SafeMap[K, V]) Set(key K, value V) {
	sm.mu.Lock()
	sm.m[key] = value
	sm.mu.Unlock()
}

// Get 获取指定键的值
func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	val, ok := sm.m[key]
	sm.mu.RUnlock()
	return val, ok
}

// GetOrDefault 获取值，不存在时返回默认值
func (sm *SafeMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	sm.mu.RLock()
	val, ok := sm.m[key]
	sm.mu.RUnlock()
	if !ok {
		return defaultValue
	}
	return val
}

// Delete 删除指定键
func (sm *SafeMap[K, V]) Delete(key K) {
	sm.mu.Lock()
	delete(sm.m, key)
	sm.mu.Unlock()
}

// Has 检查键是否存在
func (sm *SafeMap[K, V]) Has(key K) bool {
	sm.mu.RLock()
	_, ok := sm.m[key]
	sm.mu.RUnlock()
	return ok
}

// Len 返回元素数量
func (sm *SafeMap[K, V]) Len() int {
	sm.mu.RLock()
	length := len(sm.m)
	sm.mu.RUnlock()
	return length
}

// Clear 清空所有元素
func (sm *SafeMap[K, V]) Clear() {
	sm.mu.Lock()
	sm.m = make(map[K]V)
	sm.mu.Unlock()
}

// Keys 返回所有键的切片
func (sm *SafeMap[K, V]) Keys() []K {
	sm.mu.RLock()
	keys := make([]K, 0, len(sm.m))
	for k := range sm.m {
		keys = append(keys, k)
	}
	sm.mu.RUnlock()
	return keys
}

// Values 返回所有值的切片
func (sm *SafeMap[K, V]) Values() []V {
	sm.mu.RLock()
	values := make([]V, 0, len(sm.m))
	for _, v := range sm.m {
		values = append(values, v)
	}
	sm.mu.RUnlock()
	return values
}

// Range 遍历所有键值对，返回false可提前终止
// 注意：回调在读锁内执行，若在回调中调用需要写锁的方法（如 Set/Delete）可能造成死锁。
// 若需要在回调中对 SafeMap 进行修改，请使用 RangeSafe。
func (sm *SafeMap[K, V]) Range(f func(key K, value V) bool) {
	sm.mu.RLock()
	for k, v := range sm.m {
		if !f(k, v) {
			sm.mu.RUnlock()
			return
		}
	}
	sm.mu.RUnlock()
}

// RangeSafe 使用快照进行遍历，不在锁内执行回调，
// 因此在回调中调用 SafeMap 的其它方法（如 Set/Delete/Get 等）不会死锁。
// 遍历期间对底层 Map 的修改不会反映到本次回调序列（弱一致性）。
func (sm *SafeMap[K, V]) RangeSafe(f func(key K, value V) bool) {
	// 读锁下做一次快照
	sm.mu.RLock()
	type pair struct {
		key K
		val V
	}
	items := make([]pair, 0, len(sm.m))
	for k, v := range sm.m {
		items = append(items, pair{key: k, val: v})
	}
	sm.mu.RUnlock()

	// 在无锁状态下执行回调，避免死锁
	for _, it := range items {
		if !f(it.key, it.val) {
			return
		}
	}
}

// SetIfAbsent 仅当键不存在时设置值，返回是否设置成功
func (sm *SafeMap[K, V]) SetIfAbsent(key K, value V) bool {
	sm.mu.Lock()
	if _, exists := sm.m[key]; exists {
		sm.mu.Unlock()
		return false
	}
	sm.m[key] = value
	sm.mu.Unlock()
	return true
}

// GetAndDelete 获取并删除键值对
func (sm *SafeMap[K, V]) GetAndDelete(key K) (V, bool) {
	sm.mu.Lock()
	val, ok := sm.m[key]
	if ok {
		delete(sm.m, key)
	}
	sm.mu.Unlock()
	return val, ok
}

// Copy 返回map的浅拷贝
func (sm *SafeMap[K, V]) Copy() map[K]V {
	sm.mu.RLock()
	copied := make(map[K]V, len(sm.m))
	for k, v := range sm.m {
		copied[k] = v
	}
	sm.mu.RUnlock()
	return copied
}
