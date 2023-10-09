package mcopy

import (
	"testing"
)

func TestFieldCopy(t *testing.T) {
	type A struct {
		Name string
		Age  int
	}
	type B struct {
		Name string
		Num  int
	}
	a := &A{}
	err := CopySameNamedField(a, B{Name: "mack", Num: 10})
	if err != nil {
		panic(err)
	}
	t.Log(a)
}
