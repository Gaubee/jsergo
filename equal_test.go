package jsergo

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Func_isEqual(t *testing.T) {
	var source_fun = func() string {
		return "QAQ"
	}
	var target_fun = source_fun
	var fun_ref = &source_fun
	assert.Equal(t, true, IsFuncEqual(source_fun, target_fun), "别名引用的两个函数相等")
	assert.Equal(t, true, IsFuncEqual(source_fun, *fun_ref), "别名引用的两个函数相等")
}
