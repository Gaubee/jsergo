package jsergo

import (
	"reflect"
)

type anyFun func(interface{}) interface{}
type any interface{}

func IsFuncEqual(this any, fun any) bool {
	var this_pointer = reflect.ValueOf(this).Pointer()
	var fun_pointer = reflect.ValueOf(fun).Pointer()
	return this_pointer == fun_pointer
}
