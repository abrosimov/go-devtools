package typesx_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/abrosimov/go-devtools/typesx"
)

func TestTypedSyncMap_Store_Load(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     int
		wantFound bool
	}{
		{
			name:      "store and load existing key",
			key:       "test-key",
			value:     42,
			wantFound: true,
		},
		{
			name:      "load non-existent key",
			key:       "missing-key",
			value:     0,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm typesx.TypedSyncMap[string, int]

			if tt.wantFound {
				tm.Store(tt.key, tt.value)
			}

			got, ok := tm.Load(tt.key)
			if ok != tt.wantFound {
				t.Errorf("Load() ok = %v, want %v", ok, tt.wantFound)
			}
			if tt.wantFound && got != tt.value {
				t.Errorf("Load() got = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestTypedSyncMap_LoadOrStore(t *testing.T) {
	//nolint:govet // test struct alignment is acceptable
	tests := []struct {
		name       string
		setupKey   string
		setupValue int
		testKey    string
		testValue  int
		wantLoaded bool
		wantActual int
	}{
		{
			name:       "load existing value",
			setupKey:   "existing",
			setupValue: 100,
			testKey:    "existing",
			testValue:  200,
			wantLoaded: true,
			wantActual: 100,
		},
		{
			name:       "store new value",
			setupKey:   "other",
			setupValue: 100,
			testKey:    "new-key",
			testValue:  300,
			wantLoaded: false,
			wantActual: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm typesx.TypedSyncMap[string, int]
			tm.Store(tt.setupKey, tt.setupValue)

			actual, loaded := tm.LoadOrStore(tt.testKey, tt.testValue)
			if loaded != tt.wantLoaded {
				t.Errorf("LoadOrStore() loaded = %v, want %v", loaded, tt.wantLoaded)
			}
			if actual != tt.wantActual {
				t.Errorf("LoadOrStore() actual = %v, want %v", actual, tt.wantActual)
			}
		})
	}
}

func TestTypedSyncMap_Delete(t *testing.T) {
	var tm typesx.TypedSyncMap[string, int]

	tm.Store("key1", 100)
	tm.Store("key2", 200)

	if _, ok := tm.Load("key1"); !ok {
		t.Fatal("key1 should exist before delete")
	}

	tm.Delete("key1")

	if _, ok := tm.Load("key1"); ok {
		t.Error("key1 should not exist after delete")
	}

	if val, ok := tm.Load("key2"); !ok || val != 200 {
		t.Error("key2 should still exist after deleting key1")
	}
}

func TestTypedSyncMap_Range(t *testing.T) {
	var tm typesx.TypedSyncMap[string, int]

	expected := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range expected {
		tm.Store(k, v)
	}

	got := make(map[string]int)
	tm.Range(func(key string, value int) bool {
		got[key] = value
		return true
	})

	if len(got) != len(expected) {
		t.Errorf("Range() collected %d items, want %d", len(got), len(expected))
	}

	for k, wantV := range expected {
		if gotV, ok := got[k]; !ok {
			t.Errorf("Range() missing key %q", k)
		} else if gotV != wantV {
			t.Errorf("Range() key %q = %v, want %v", k, gotV, wantV)
		}
	}
}

func TestTypedSyncMap_Range_EarlyStop(t *testing.T) {
	var tm typesx.TypedSyncMap[string, int]

	for i := 0; i < 10; i++ {
		tm.Store(string(rune('a'+i)), i)
	}

	count := 0
	tm.Range(func(key string, value int) bool {
		count++
		return count < 5
	})

	if count != 5 {
		t.Errorf("Range() stopped after %d iterations, want 5", count)
	}
}

func TestTypedSyncMap_Clear(t *testing.T) {
	var tm typesx.TypedSyncMap[string, int]

	tm.Store("key1", 1)
	tm.Store("key2", 2)
	tm.Store("key3", 3)

	tm.Clear()

	count := 0
	tm.Range(func(key string, value int) bool {
		count++
		return true
	})

	if count != 0 {
		t.Errorf("Clear() left %d items in map, want 0", count)
	}

	if _, ok := tm.Load("key1"); ok {
		t.Error("key1 should not exist after Clear()")
	}
}

func TestTypedSyncMap_ZeroValues(t *testing.T) {
	var tm typesx.TypedSyncMap[string, int]

	tm.Store("zero", 0)

	val, ok := tm.Load("zero")
	if !ok {
		t.Error("Load() should return true for stored zero value")
	}
	if val != 0 {
		t.Errorf("Load() got = %v, want 0", val)
	}
}

func TestTypedSyncMap_ConcurrentAccess(t *testing.T) {
	var tm typesx.TypedSyncMap[int, int]
	var wg sync.WaitGroup

	numGoroutines := 100
	numOperations := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				tm.Store(key, key*2)
				tm.Load(key)
				if j%2 == 0 {
					tm.Delete(key)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestTypedSyncMap_DifferentTypes(t *testing.T) {
	t.Run("string to struct", func(t *testing.T) {
		type person struct {
			name string
			age  int
		}

		var tm typesx.TypedSyncMap[string, person]

		p := person{name: "Alice", age: 30}
		tm.Store("alice", p)

		got, ok := tm.Load("alice")
		if !ok {
			t.Fatal("Load() should find stored value")
		}
		if got != p {
			t.Errorf("Load() got = %+v, want %+v", got, p)
		}
	})

	t.Run("int to pointer", func(t *testing.T) {
		var tm typesx.TypedSyncMap[int, *string]

		str := "test"
		tm.Store(1, &str)

		got, ok := tm.Load(1)
		if !ok {
			t.Fatal("Load() should find stored value")
		}
		if got == nil {
			t.Fatal("Load() returned nil pointer")
		}
		if *got != str {
			t.Errorf("Load() got = %q, want %q", *got, str)
		}
	})
}

// ExampleTypedSyncMap_Store demonstrates storing values in a TypedSyncMap.
func ExampleTypedSyncMap_Store() {
	var m typesx.TypedSyncMap[string, int]

	m.Store("age", 30)
	m.Store("count", 42)

	val, ok := m.Load("age")
	fmt.Println(val, ok)
	// Output: 30 true
}

// ExampleTypedSyncMap_Load demonstrates loading values from a TypedSyncMap.
func ExampleTypedSyncMap_Load() {
	var m typesx.TypedSyncMap[string, string]

	m.Store("name", "Alice")

	// Load existing key
	val, ok := m.Load("name")
	fmt.Println(val, ok)

	// Load non-existent key
	val, ok = m.Load("missing")
	fmt.Println(val, ok)

	// Output:
	// Alice true
	//  false
}

// ExampleTypedSyncMap_LoadOrStore demonstrates LoadOrStore functionality.
func ExampleTypedSyncMap_LoadOrStore() {
	var m typesx.TypedSyncMap[string, int]

	// Store new value
	actual, loaded := m.LoadOrStore("count", 100)
	fmt.Println(actual, loaded)

	// Try to store again - returns existing value
	actual, loaded = m.LoadOrStore("count", 200)
	fmt.Println(actual, loaded)

	// Output:
	// 100 false
	// 100 true
}

// ExampleTypedSyncMap_Delete demonstrates deleting entries from a TypedSyncMap.
func ExampleTypedSyncMap_Delete() {
	var m typesx.TypedSyncMap[string, int]

	m.Store("key1", 1)
	m.Store("key2", 2)

	m.Delete("key1")

	_, ok := m.Load("key1")
	fmt.Println("key1 exists:", ok)

	val, ok := m.Load("key2")
	fmt.Println("key2 value:", val, "exists:", ok)

	// Output:
	// key1 exists: false
	// key2 value: 2 exists: true
}

// ExampleTypedSyncMap_Range demonstrates iterating over a TypedSyncMap.
func ExampleTypedSyncMap_Range() {
	var m typesx.TypedSyncMap[string, int]

	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	sum := 0
	m.Range(func(key string, value int) bool {
		sum += value
		return true
	})

	fmt.Println(sum)
	// Output: 6
}

// ExampleTypedSyncMap_Clear demonstrates clearing all entries from a TypedSyncMap.
func ExampleTypedSyncMap_Clear() {
	var m typesx.TypedSyncMap[string, int]

	m.Store("a", 1)
	m.Store("b", 2)

	m.Clear()

	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})

	fmt.Println("entries after clear:", count)
	// Output: entries after clear: 0
}

func TestTypedSyncMap_LoadAndDelete(t *testing.T) {
	//nolint:govet // test struct alignment is acceptable
	tests := []struct {
		name       string
		setupKey   string
		setupValue int
		deleteKey  string
		wantValue  int
		wantLoaded bool
	}{
		{
			name:       "delete existing key",
			setupKey:   "key1",
			setupValue: 100,
			deleteKey:  "key1",
			wantValue:  100,
			wantLoaded: true,
		},
		{
			name:       "delete non-existent key",
			setupKey:   "key1",
			setupValue: 100,
			deleteKey:  "key2",
			wantValue:  0,
			wantLoaded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm typesx.TypedSyncMap[string, int]
			tm.Store(tt.setupKey, tt.setupValue)

			gotValue, gotLoaded := tm.LoadAndDelete(tt.deleteKey)
			if gotLoaded != tt.wantLoaded {
				t.Errorf("LoadAndDelete() loaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
			if gotValue != tt.wantValue {
				t.Errorf("LoadAndDelete() value = %v, want %v", gotValue, tt.wantValue)
			}

			// Verify key is actually deleted
			if tt.wantLoaded {
				if _, ok := tm.Load(tt.deleteKey); ok {
					t.Error("key should be deleted after LoadAndDelete")
				}
			}
		})
	}
}

func TestTypedSyncMap_Swap(t *testing.T) {
	//nolint:govet // test struct alignment is acceptable
	tests := []struct {
		name         string
		setupKey     string
		setupValue   int
		swapKey      string
		swapValue    int
		wantPrevious int
		wantLoaded   bool
	}{
		{
			name:         "swap existing key",
			setupKey:     "key1",
			setupValue:   100,
			swapKey:      "key1",
			swapValue:    200,
			wantPrevious: 100,
			wantLoaded:   true,
		},
		{
			name:         "swap non-existent key",
			setupKey:     "key1",
			setupValue:   100,
			swapKey:      "key2",
			swapValue:    200,
			wantPrevious: 0,
			wantLoaded:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm typesx.TypedSyncMap[string, int]
			tm.Store(tt.setupKey, tt.setupValue)

			gotPrevious, gotLoaded := tm.Swap(tt.swapKey, tt.swapValue)
			if gotLoaded != tt.wantLoaded {
				t.Errorf("Swap() loaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
			if gotPrevious != tt.wantPrevious {
				t.Errorf("Swap() previous = %v, want %v", gotPrevious, tt.wantPrevious)
			}

			// Verify new value is stored
			if val, ok := tm.Load(tt.swapKey); !ok || val != tt.swapValue {
				t.Errorf("Swap() should store new value, got %v", val)
			}
		})
	}
}

//nolint:govet // test struct alignment is acceptable
type compareAndSwapTestCase struct {
	name       string
	setupKey   string
	setupValue int
	casKey     string
	oldValue   int
	newValue   int
	wantSwap   bool
}

func getCompareAndSwapTestCases() []compareAndSwapTestCase {
	return []compareAndSwapTestCase{
		{
			name:       "successful compare and swap",
			setupKey:   "key1",
			setupValue: 100,
			casKey:     "key1",
			oldValue:   100,
			newValue:   200,
			wantSwap:   true,
		},
		{
			name:       "failed compare and swap - wrong old value",
			setupKey:   "key1",
			setupValue: 100,
			casKey:     "key1",
			oldValue:   999,
			newValue:   200,
			wantSwap:   false,
		},
		{
			name:       "failed compare and swap - key doesn't exist",
			setupKey:   "key1",
			setupValue: 100,
			casKey:     "key2",
			oldValue:   0,
			newValue:   200,
			wantSwap:   false,
		},
	}
}

func TestTypedSyncMap_CompareAndSwap(t *testing.T) {
	tests := getCompareAndSwapTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm typesx.TypedSyncMap[string, int]
			tm.Store(tt.setupKey, tt.setupValue)

			gotSwap := tm.CompareAndSwap(tt.casKey, tt.oldValue, tt.newValue)
			if gotSwap != tt.wantSwap {
				t.Errorf("CompareAndSwap() = %v, want %v", gotSwap, tt.wantSwap)
			}

			// Verify value is correct
			if val, ok := tm.Load(tt.casKey); ok {
				expectedValue := tt.setupValue
				if tt.wantSwap && tt.casKey == tt.setupKey {
					expectedValue = tt.newValue
				}
				if val != expectedValue {
					t.Errorf("After CompareAndSwap(), value = %v, want %v", val, expectedValue)
				}
			}
		})
	}
}

//nolint:govet // test struct alignment is acceptable
type compareAndDeleteTestCase struct {
	name        string
	setupKey    string
	setupValue  int
	deleteKey   string
	oldValue    int
	wantDeleted bool
}

func getCompareAndDeleteTestCases() []compareAndDeleteTestCase {
	return []compareAndDeleteTestCase{
		{
			name:        "successful compare and delete",
			setupKey:    "key1",
			setupValue:  100,
			deleteKey:   "key1",
			oldValue:    100,
			wantDeleted: true,
		},
		{
			name:        "failed compare and delete - wrong old value",
			setupKey:    "key1",
			setupValue:  100,
			deleteKey:   "key1",
			oldValue:    999,
			wantDeleted: false,
		},
		{
			name:        "failed compare and delete - key doesn't exist",
			setupKey:    "key1",
			setupValue:  100,
			deleteKey:   "key2",
			oldValue:    0,
			wantDeleted: false,
		},
	}
}

func TestTypedSyncMap_CompareAndDelete(t *testing.T) {
	tests := getCompareAndDeleteTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm typesx.TypedSyncMap[string, int]
			tm.Store(tt.setupKey, tt.setupValue)

			gotDeleted := tm.CompareAndDelete(tt.deleteKey, tt.oldValue)
			if gotDeleted != tt.wantDeleted {
				t.Errorf("CompareAndDelete() = %v, want %v", gotDeleted, tt.wantDeleted)
			}

			// Verify key state
			_, exists := tm.Load(tt.deleteKey)
			if tt.wantDeleted && exists {
				t.Error("key should be deleted after successful CompareAndDelete")
			}
			if !tt.wantDeleted && tt.deleteKey == tt.setupKey && !exists {
				t.Error("key should still exist after failed CompareAndDelete")
			}
		})
	}
}

// ExampleNewTypedSyncMap demonstrates creating a new TypedSyncMap.
func ExampleNewTypedSyncMap() {
	m := typesx.NewTypedSyncMap[string, int]()

	m.Store("count", 42)
	val, ok := m.Load("count")
	fmt.Println(val, ok)
	// Output: 42 true
}

// ExampleTypedSyncMap_LoadAndDelete demonstrates atomic load and delete.
func ExampleTypedSyncMap_LoadAndDelete() {
	m := typesx.NewTypedSyncMap[string, int]()

	m.Store("temp", 100)

	// Load and delete in one atomic operation
	val, loaded := m.LoadAndDelete("temp")
	fmt.Println("value:", val, "loaded:", loaded)

	// Verify it's deleted
	_, exists := m.Load("temp")
	fmt.Println("exists:", exists)

	// Output:
	// value: 100 loaded: true
	// exists: false
}

// ExampleTypedSyncMap_Swap demonstrates atomic value swap.
func ExampleTypedSyncMap_Swap() {
	m := typesx.NewTypedSyncMap[string, int]()

	m.Store("counter", 10)

	// Swap with new value, get old value
	old, loaded := m.Swap("counter", 20)
	fmt.Println("old:", old, "loaded:", loaded)

	// Verify new value
	val, _ := m.Load("counter")
	fmt.Println("new:", val)

	// Output:
	// old: 10 loaded: true
	// new: 20
}

// ExampleTypedSyncMap_CompareAndSwap demonstrates conditional atomic swap.
func ExampleTypedSyncMap_CompareAndSwap() {
	m := typesx.NewTypedSyncMap[string, int]()

	m.Store("version", 1)

	// Only swap if current value is 1
	swapped := m.CompareAndSwap("version", 1, 2)
	fmt.Println("swapped:", swapped)

	// Try to swap with wrong old value
	swapped = m.CompareAndSwap("version", 1, 3)
	fmt.Println("swapped:", swapped)

	val, _ := m.Load("version")
	fmt.Println("final:", val)

	// Output:
	// swapped: true
	// swapped: false
	// final: 2
}

// ExampleTypedSyncMap_CompareAndDelete demonstrates conditional atomic delete.
func ExampleTypedSyncMap_CompareAndDelete() {
	m := typesx.NewTypedSyncMap[string, int]()

	m.Store("item", 100)

	// Try to delete with wrong value
	deleted := m.CompareAndDelete("item", 999)
	fmt.Println("deleted:", deleted)

	// Delete with correct value
	deleted = m.CompareAndDelete("item", 100)
	fmt.Println("deleted:", deleted)

	_, exists := m.Load("item")
	fmt.Println("exists:", exists)

	// Output:
	// deleted: false
	// deleted: true
	// exists: false
}
