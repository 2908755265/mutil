package mcopy

import (
	"testing"
)

func TestFieldCopy(t *testing.T) {
	type A struct {
		Name string
		Age  int
		Num  int
	}
	type B struct {
		Name string
		Num  int
	}
	a := &A{Name: "alice", Num: 10}
	err := CopySameNamedField(a, B{Name: "", Num: -1}, true)
	if err != nil {
		panic(err)
	}
	t.Log(a)
}
