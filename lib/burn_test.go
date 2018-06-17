package trunks

import (
	"reflect"
	"testing"
)

type TestInterface interface {
	Hello() string
}

type TestStruct struct {
	Name string
}

// TestStruct implemntes the interface
func (ts TestStruct) Hello() string {
	return "hello, " + ts.Name
}

type TestStruct2 struct {
	Name string
}

// *TestStruct2 implements the interface
func (ts *TestStruct2) Hello() string {
	return "hello2, " + ts.Name
}

func cloneNew(obj TestInterface) TestInterface {
	t := reflect.ValueOf(obj).Type()
	return reflect.New(t).Elem().Interface().(*TestStruct2)
}

func testFunc(obj TestInterface) string {
	return obj.Hello()
}

func TestCloneObj(t *testing.T) {
	obj := TestStruct{Name: "dave"}
	t.Logf(testFunc(obj))

}
