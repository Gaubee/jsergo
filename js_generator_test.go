package jsergo

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Generator(t *testing.T) {
	assert.Equal(t, 1, 1)
	var generator = newGenerator(func(this *Generator) {
		var index = 0
		for index < 3 {
			index += 1
			this.Yield(index)
		}
		this.Return(666)
	})

	assert.Equal(t, generator.Next().value, 1)
	assert.Equal(t, generator.Next().value, 2)
	assert.Equal(t, generator.Next().value, 3)
	assert.Equal(t, generator.Next().done, true)
	assert.Equal(t, generator.Next().value, 666)
}
