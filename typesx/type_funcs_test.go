package typesx_test

import (
	"fmt"
	"testing"

	"github.com/abrosimov/go-devtools/typesx"
)

type ExportedStruct struct{}
type ExportedInterface interface{}
type ExportedGeneric[T any] struct{}

type unexportedStruct struct{}
type unexportedInterface interface{}
type unexportedGeneric[T any] struct{}

type typeNameTestCase struct {
	name     string
	getType  func() string
	wantType string
}

func getTypeNameTestCases() []typeNameTestCase {
	return []typeNameTestCase{
		{
			name:     "unexported struct",
			getType:  func() string { return typesx.GetTypeFQN[unexportedStruct]() },
			wantType: "github.com/abrosimov/go-devtools/typesx_test.unexportedStruct",
		},
		{
			name:     "unexported interface",
			getType:  func() string { return typesx.GetTypeFQN[unexportedInterface]() },
			wantType: "github.com/abrosimov/go-devtools/typesx_test.unexportedInterface",
		},
		{
			name:     "unexported generic with int",
			getType:  func() string { return typesx.GetTypeFQN[unexportedGeneric[int]]() },
			wantType: "github.com/abrosimov/go-devtools/typesx_test.unexportedGeneric[int]",
		},
		{
			name:     "unexported generic with float32",
			getType:  func() string { return typesx.GetTypeFQN[unexportedGeneric[float32]]() },
			wantType: "github.com/abrosimov/go-devtools/typesx_test.unexportedGeneric[float32]",
		},
		{
			name:     "exported struct",
			getType:  func() string { return typesx.GetTypeFQN[ExportedStruct]() },
			wantType: "github.com/abrosimov/go-devtools/typesx_test.ExportedStruct",
		},
		{
			name:     "exported interface",
			getType:  func() string { return typesx.GetTypeFQN[ExportedInterface]() },
			wantType: "github.com/abrosimov/go-devtools/typesx_test.ExportedInterface",
		},
		{
			name:     "exported generic with int",
			getType:  func() string { return typesx.GetTypeFQN[ExportedGeneric[int]]() },
			wantType: "github.com/abrosimov/go-devtools/typesx_test.ExportedGeneric[int]",
		},
		{
			name:     "exported generic with float32",
			getType:  func() string { return typesx.GetTypeFQN[ExportedGeneric[float32]]() },
			wantType: "github.com/abrosimov/go-devtools/typesx_test.ExportedGeneric[float32]",
		},
	}
}

func TestGetTypeName(t *testing.T) {
	tests := getTypeNameTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.getType()
			if got != tt.wantType {
				t.Errorf("GetTypeFQN() = %q, want %q", got, tt.wantType)
			}
		})
	}
}

func TestIsInterface(t *testing.T) {
	//nolint:govet // test struct alignment is acceptable
	tests := []struct {
		name        string
		checkType   func() bool
		wantIsIface bool
	}{
		{
			name:        "struct is not interface",
			checkType:   func() bool { return typesx.IsInterface[ExportedStruct]() },
			wantIsIface: false,
		},
		{
			name:        "int is not interface",
			checkType:   func() bool { return typesx.IsInterface[int]() },
			wantIsIface: false,
		},
		{
			name:        "generic struct is not interface",
			checkType:   func() bool { return typesx.IsInterface[ExportedGeneric[int]]() },
			wantIsIface: false,
		},
		{
			name:        "exported interface is interface",
			checkType:   func() bool { return typesx.IsInterface[ExportedInterface]() },
			wantIsIface: true,
		},
		{
			name:        "fmt.Stringer is interface",
			checkType:   func() bool { return typesx.IsInterface[fmt.Stringer]() },
			wantIsIface: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.checkType()
			if got != tt.wantIsIface {
				t.Errorf("IsInterface() = %v, want %v", got, tt.wantIsIface)
			}
		})
	}
}

func TestHasPtrFieldOfType(t *testing.T) {
	type expectedStruct struct{}
	type structWithPtr struct {
		*expectedStruct //nolint:unused // don't need linter here.
	}
	type structWithVal struct {
		expectedStruct //nolint:unused // don't need linter here.
	}

	//nolint:govet // test struct alignment is acceptable
	tests := []struct {
		name    string
		check   func() bool
		wantHas bool
	}{
		{
			name:    "struct with pointer field",
			check:   func() bool { return typesx.HasPtrFieldOfType[structWithPtr, expectedStruct]() },
			wantHas: true,
		},
		{
			name:    "struct with value field",
			check:   func() bool { return typesx.HasPtrFieldOfType[structWithVal, expectedStruct]() },
			wantHas: false,
		},
		{
			name:    "generic type without pointer",
			check:   func() bool { return typesx.HasPtrFieldOfType[ExportedGeneric[int], expectedStruct]() },
			wantHas: false,
		},
		{
			name:    "exported struct without pointer",
			check:   func() bool { return typesx.HasPtrFieldOfType[ExportedStruct, expectedStruct]() },
			wantHas: false,
		},
		{
			name:    "interface type",
			check:   func() bool { return typesx.HasPtrFieldOfType[fmt.Stringer, expectedStruct]() },
			wantHas: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.check()
			if got != tt.wantHas {
				t.Errorf("HasPtrFieldOfType() = %v, want %v", got, tt.wantHas)
			}
		})
	}
}

// ExampleIsInterface demonstrates checking if a type is an interface.
func ExampleIsInterface() {
	// Check if fmt.Stringer is an interface
	fmt.Println(typesx.IsInterface[fmt.Stringer]())

	// Check if a struct is an interface
	fmt.Println(typesx.IsInterface[ExportedStruct]())

	// Check if int is an interface
	fmt.Println(typesx.IsInterface[int]())

	// Output:
	// true
	// false
	// false
}

// ExampleGetTypeFQN demonstrates getting fully qualified type names.
func ExampleGetTypeFQN() {
	// Get FQN for a local struct
	fmt.Println(typesx.GetTypeFQN[ExportedStruct]())

	// Get FQN for a generic type
	fmt.Println(typesx.GetTypeFQN[ExportedGeneric[int]]())

	// Output:
	// github.com/abrosimov/go-devtools/typesx_test.ExportedStruct
	// github.com/abrosimov/go-devtools/typesx_test.ExportedGeneric[int]
}

// ExampleHasPtrFieldOfType demonstrates checking if a struct has a pointer field of specific type.
func ExampleHasPtrFieldOfType() {
	type Inner struct {
		value int //nolint:unused // used for type demonstration
	}

	type WithPtr struct {
		data *Inner //nolint:unused // used for type demonstration
	}

	type WithValue struct {
		data Inner //nolint:unused // used for type demonstration
	}

	// Check if WithPtr has a pointer to Inner
	fmt.Println(typesx.HasPtrFieldOfType[WithPtr, Inner]())

	// Check if WithValue has a pointer to Inner
	fmt.Println(typesx.HasPtrFieldOfType[WithValue, Inner]())

	// Output:
	// true
	// false
}
