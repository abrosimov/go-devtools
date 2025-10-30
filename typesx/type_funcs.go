// Package typesx provides some type-related functions.
package typesx

import (
	"fmt"
	"reflect"
)

// IsInterface return true if T is an interface, otherwise false.
func IsInterface[T any]() bool {
	t := reflect.TypeFor[T]()
	return t.Kind() == reflect.Interface
}

// HasPtrFieldOfType returns true if T has a field with type *V, otherwise false.
// Only checks direct fields, not embedded fields.
func HasPtrFieldOfType[T any, V any]() bool {
	t := reflect.TypeFor[T]()
	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type == reflect.TypeOf((*V)(nil)) {
			return true
		}
	}
	return false
}

// GetTypeFQN returns a fully qualified name of the type including package name.
func GetTypeFQN[T any]() string {
	t := reflect.TypeFor[T]()
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}
