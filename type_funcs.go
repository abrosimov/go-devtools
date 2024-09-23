package devtools

import (
	"fmt"
	"reflect"
)

func IsInterface[T any]() bool {
	t := reflect.TypeFor[T]()
	return t.Kind() == reflect.Interface
}

func IfTypeHasPtrToV[T any, V any]() bool {
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

func GetTypeName[T any]() string {
	t := reflect.TypeFor[T]()
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}
