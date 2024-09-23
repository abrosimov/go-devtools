package devtools_test

import (
	"fmt"
	"testing"

	"github.com/abrosimov/go-devtools"
)

type ExportedStruct struct{}
type ExportedInterface interface{}
type ExportedGeneric[T any] struct{}

type unexportedStruct struct{}
type unexportedInterface interface{}
type unexportedGeneric[T any] struct{}

func TestGetTypeName(t *testing.T) {
	unexportedStructTypeName := devtools.GetTypeName[unexportedStruct]()
	if unexportedStructTypeName != "github.com/abrosimov/go-devtools_test.unexportedStruct" {
		t.Errorf("unexportedStruct (func scope) has invalid type name")
	}

	myLocalInterfaceType := devtools.GetTypeName[unexportedInterface]()
	if myLocalInterfaceType != "github.com/abrosimov/go-devtools_test.unexportedInterface" {
		t.Errorf("unexportedInterface has invalid type name")
	}

	myLocalGenericTypeForInt := devtools.GetTypeName[unexportedGeneric[int]]()
	if myLocalGenericTypeForInt != "github.com/abrosimov/go-devtools_test.unexportedGeneric[int]" {
		t.Errorf("unexportedGeneric[int] has invalid type name")
	}

	myLocalGenericTypeForFloat := devtools.GetTypeName[unexportedGeneric[float32]]()
	if myLocalGenericTypeForFloat != "github.com/abrosimov/go-devtools_test.unexportedGeneric[float32]" {
		t.Errorf("unexportedGeneric[float32] has invalid type name")
	}

	ExportedStructTypeName := devtools.GetTypeName[ExportedStruct]()
	if ExportedStructTypeName != "github.com/abrosimov/go-devtools_test.ExportedStruct" {
		t.Errorf("ExportedStruct has invalid type name")
	}

	myGlobalInterfaceType := devtools.GetTypeName[ExportedInterface]()
	if myGlobalInterfaceType != "github.com/abrosimov/go-devtools_test.ExportedInterface" {
		t.Errorf("")
	}

	myGlobalGenericTypeForInt := devtools.GetTypeName[ExportedGeneric[int]]()
	if myGlobalGenericTypeForInt != "github.com/abrosimov/go-devtools_test.ExportedGeneric[int]" {
		t.Errorf("")
	}

	myGlobalGenericTypeForFloat := devtools.GetTypeName[ExportedGeneric[float32]]()
	if myGlobalGenericTypeForFloat != "github.com/abrosimov/go-devtools_test.ExportedGeneric[float32]" {
		t.Errorf("")
	}
}

func TestIsInterface(t *testing.T) {
	if devtools.IsInterface[ExportedStruct]() {
		t.Errorf("ExportedStruct is not an interface")
	}
	if devtools.IsInterface[int]() {
		t.Errorf("int is not an interface")
	}

	if devtools.IsInterface[ExportedGeneric[int]]() {
		t.Errorf("ExportedGeneric[int] is not an interface")
	}

	if !devtools.IsInterface[ExportedInterface]() {
		t.Errorf("ExportedInterface is an interface")
	}

	if !devtools.IsInterface[fmt.Stringer]() {
		t.Errorf("fmt.Stringer is an interface")
	}
}

func TestIfTypeHasPtrToV(t *testing.T) {
	type expectedStruct struct{}
	type structWithPtr struct {
		*expectedStruct //nolint:unused // don't need linter here.
	}
	type structWithVal struct {
		expectedStruct //nolint:unused // don't need linter here.
	}

	if !devtools.IfTypeHasPtrToV[structWithPtr, expectedStruct]() {
		t.Errorf("structWithPtr should have a pointer to expectedStruct")
	}

	if devtools.IfTypeHasPtrToV[structWithVal, expectedStruct]() {
		t.Errorf("structWithVal should not have a pointer to expectedStruct")
	}

	if devtools.IfTypeHasPtrToV[ExportedGeneric[int], expectedStruct]() {
		t.Errorf("ExportedGeneric[int] should not have a pointer to expectedStruct")
	}

	if devtools.IfTypeHasPtrToV[ExportedStruct, expectedStruct]() {
		t.Errorf("ExportedStruct should not have a pointer to expectedStruct")
	}
	if devtools.IfTypeHasPtrToV[fmt.Stringer, expectedStruct]() {
		t.Errorf("fmt.Stringer should not have a pointer to expectedStruct")
	}
}
