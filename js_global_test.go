package jsergo

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseInt(t *testing.T) {
	assert.Equal(t, global.ParseInt("123"), JsNumber(123))
	assert.Equal(t, global.ParseInt(123), JsNumber(123))
	assert.Equal(t, global.ParseInt(123.3), JsNumber(123))
	assert.Equal(t, global.ParseInt("123.3"), JsNumber(123))
	assert.Equal(t, global.ParseInt("123a"), JsNumber(123))

	assert.Equal(t, JsNumber(0), global.ParseInt("0"))
	assert.Equal(t, JsNumber(0), global.ParseInt(" 0"))
	assert.Equal(t, JsNumber(0), global.ParseInt(" 0 "))

	assert.Equal(t, JsNumber(77), global.ParseInt("077"))
	assert.Equal(t, JsNumber(77), global.ParseInt("  077"))
	assert.Equal(t, JsNumber(77), global.ParseInt("  077   "))
	assert.Equal(t, JsNumber(-77), global.ParseInt("  -077"))

	assert.Equal(t, JsNumber(3), global.ParseInt("11", 2))
	assert.Equal(t, JsNumber(4), global.ParseInt("11", 3))
	assert.Equal(t, JsNumber(4), global.ParseInt("11", 3.8))

	assert.Equal(t, JsNumber(0x12), global.ParseInt("0x12"))
	assert.Equal(t, JsNumber(0x12), global.ParseInt("0x12", 16))
	assert.Equal(t, JsNumber(0x12), global.ParseInt("0x12", 16.1))
	assert.Equal(t, JsNumber(0x12), global.ParseInt("0x12", NaN))
	assert.True(t, global.IsNaN(global.ParseInt("0x  ")))
	assert.True(t, global.IsNaN(global.ParseInt("0x")))
	assert.True(t, global.IsNaN(global.ParseInt("0x  ", 16)))
	assert.True(t, global.IsNaN(global.ParseInt("0x", 16)))
	assert.Equal(t, JsNumber(12), global.ParseInt("12aaa"))

	assert.Equal(t, JsNumber(1.1), global.ParseFloat("1.1"))
	assert.Equal(t, JsNumber(1.1), global.ParseFloat("1.1aaa"))
	assert.Equal(t, JsNumber(0.1), global.ParseFloat("0.1"))
	assert.Equal(t, JsNumber(0.1), global.ParseFloat("0.1aaa"))
	assert.Equal(t, JsNumber(0.1), global.ParseFloat("00000.1aaa"))
	assert.Equal(t, JsNumber(0), global.ParseFloat("0aaa"))
	assert.Equal(t, JsNumber(0), global.ParseFloat("0x12"))
	assert.Equal(t, JsNumber(77), global.ParseFloat("077"))
}
