package typesx

import "sync"

// TypedSyncMap provides a type-safe wrapper around sync.Map.
// It ensures compile-time type safety for map operations while maintaining
// the performance characteristics of sync.Map for concurrent access patterns.
//
// TypedSyncMap must not be copied after first use.
type TypedSyncMap[K comparable, V any] struct {
	m sync.Map
}

// NewTypedSyncMap creates a new TypedSyncMap.
// The returned value must not be copied after first use.
func NewTypedSyncMap[K comparable, V any]() TypedSyncMap[K, V] {
	return TypedSyncMap[K, V]{}
}

// Load returns the value stored in the map for a key, or the zero value if no
// value is present. The ok result indicates whether value was found in the map.
func (tm *TypedSyncMap[K, V]) Load(key K) (V, bool) {
	v, ok := tm.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	//nolint:errcheck // type assertion is safe by design: all stores are controlled by this wrapper
	return v.(V), true
}

// Store sets the value for a key.
func (tm *TypedSyncMap[K, V]) Store(key K, value V) {
	tm.m.Store(key, value)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (tm *TypedSyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := tm.m.LoadOrStore(key, value)
	//nolint:errcheck // type assertion is safe by design: all stores are controlled by this wrapper
	return v.(V), loaded
}

// Delete deletes the value for a key.
func (tm *TypedSyncMap[K, V]) Delete(key K) {
	tm.m.Delete(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (tm *TypedSyncMap[K, V]) Range(f func(key K, value V) bool) {
	tm.m.Range(func(k, v any) bool {
		//nolint:errcheck // type assertion is safe by design: all stores are controlled by this wrapper
		return f(k.(K), v.(V))
	})
}

// Clear deletes all entries from the map.
func (tm *TypedSyncMap[K, V]) Clear() {
	tm.m.Clear()
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (tm *TypedSyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := tm.m.LoadAndDelete(key)
	if !loaded {
		var zero V
		return zero, false
	}
	//nolint:errcheck // type assertion is safe by design: all stores are controlled by this wrapper
	return v.(V), true
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (tm *TypedSyncMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, loaded := tm.m.Swap(key, value)
	if !loaded {
		var zero V
		return zero, false
	}
	//nolint:errcheck // type assertion is safe by design: all stores are controlled by this wrapper
	return v.(V), true
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (tm *TypedSyncMap[K, V]) CompareAndSwap(key K, old, newVal V) bool {
	return tm.m.CompareAndSwap(key, old, newVal)
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete returns false
// (even if the old value is the zero value).
func (tm *TypedSyncMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return tm.m.CompareAndDelete(key, old)
}
