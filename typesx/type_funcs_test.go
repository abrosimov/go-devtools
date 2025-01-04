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

func TestGetTypeName(t *testing.T) {
	unexportedStructTypeName := typesx.GetTypeFQN[unexportedStruct]()
	if unexportedStructTypeName != "github.com/abrosimov/go-devtools/typesx_test.unexportedStruct" {
		t.Errorf("unexportedStruct (func scope) has invalid type name")
	}

	myLocalInterfaceType := typesx.GetTypeFQN[unexportedInterface]()
	if myLocalInterfaceType != "github.com/abrosimov/go-devtools/typesx_test.unexportedInterface" {
		t.Errorf("unexportedInterface has invalid type name")
	}

	myLocalGenericTypeForInt := typesx.GetTypeFQN[unexportedGeneric[int]]()
	if myLocalGenericTypeForInt != "github.com/abrosimov/go-devtools/typesx_test.unexportedGeneric[int]" {
		t.Errorf("unexportedGeneric[int] has invalid type name")
	}

	myLocalGenericTypeForFloat := typesx.GetTypeFQN[unexportedGeneric[float32]]()
	if myLocalGenericTypeForFloat != "github.com/abrosimov/go-devtools/typesx_test.unexportedGeneric[float32]" {
		t.Errorf("unexportedGeneric[float32] has invalid type name")
	}

	ExportedStructTypeName := typesx.GetTypeFQN[ExportedStruct]()
	if ExportedStructTypeName != "github.com/abrosimov/go-devtools/typesx_test.ExportedStruct" {
		t.Errorf("ExportedStruct has invalid type name")
	}

	myGlobalInterfaceType := typesx.GetTypeFQN[ExportedInterface]()
	if myGlobalInterfaceType != "github.com/abrosimov/go-devtools/typesx_test.ExportedInterface" {
		t.Errorf("myGlobalInterfaceType has invalid type name")
	}

	myGlobalGenericTypeForInt := typesx.GetTypeFQN[ExportedGeneric[int]]()
	if myGlobalGenericTypeForInt != "github.com/abrosimov/go-devtools/typesx_test.ExportedGeneric[int]" {
		t.Errorf("myGlobalGenericTypeForInt has invalid type name")
	}

	myGlobalGenericTypeForFloat := typesx.GetTypeFQN[ExportedGeneric[float32]]()
	if myGlobalGenericTypeForFloat != "github.com/abrosimov/go-devtools/typesx_test.ExportedGeneric[float32]" {
		t.Errorf("myGlobalGenericTypeForFloat has invalid type name")
	}
}

func TestIsInterface(t *testing.T) {
	if typesx.IsInterface[ExportedStruct]() {
		t.Errorf("ExportedStruct is not an interface")
	}
	if typesx.IsInterface[int]() {
		t.Errorf("int is not an interface")
	}

	if typesx.IsInterface[ExportedGeneric[int]]() {
		t.Errorf("ExportedGeneric[int] is not an interface")
	}

	if !typesx.IsInterface[ExportedInterface]() {
		t.Errorf("ExportedInterface is an interface")
	}

	if !typesx.IsInterface[fmt.Stringer]() {
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

	if !typesx.IfTypeHasPtrToV[structWithPtr, expectedStruct]() {
		t.Errorf("structWithPtr should have a pointer to expectedStruct")
	}

	if typesx.IfTypeHasPtrToV[structWithVal, expectedStruct]() {
		t.Errorf("structWithVal should not have a pointer to expectedStruct")
	}

	if typesx.IfTypeHasPtrToV[ExportedGeneric[int], expectedStruct]() {
		t.Errorf("ExportedGeneric[int] should not have a pointer to expectedStruct")
	}

	if typesx.IfTypeHasPtrToV[ExportedStruct, expectedStruct]() {
		t.Errorf("ExportedStruct should not have a pointer to expectedStruct")
	}
	if typesx.IfTypeHasPtrToV[fmt.Stringer, expectedStruct]() {
		t.Errorf("fmt.Stringer should not have a pointer to expectedStruct")
	}
}
