package jsergo

import (
// "fmt"
// "math/rand"
// "sync"
// "time"
)

type Generator struct {
	ch               chan interface{}
	handle           func(*Generator)
	is_started       bool
	done_cache_value interface{}
	is_done          bool
}
type GeneratorData struct {
	value interface{}
	done  bool
}

func newGenerator(handle func(*Generator)) *Generator {
	var res = &Generator{
		ch:         make(chan interface{}),
		handle:     handle,
		is_started: false,
		is_done:    false,
	}
	return res
}

func (this *Generator) Next(input_args ...interface{}) *GeneratorData {
	if this.is_done {
		return &GeneratorData{
			value: this.done_cache_value,
			done:  true,
		}
	}
	var res = &GeneratorData{}
	if this.is_started {
		// 输入参数
		if len(input_args) == 0 {
			this.ch <- nil
		} else {
			this.ch <- input_args[0]
		}
		// 等待Yield/Return输出返回值
		res.value = <-this.ch
	} else {
		this.is_started = true
		go this.handle(this)
		// 等待Yield/Return输出返回值
		res.value = <-this.ch
	}
	res.done = this.is_done
	return res
}

func (this *Generator) Yield(output_args ...interface{}) interface{} {
	// 输出返回值
	if len(output_args) == 0 {
		this.ch <- nil
	} else {
		this.ch <- output_args[0]
	}
	// 等待Next输入
	input_arg := <-this.ch
	return input_arg
}

func (this *Generator) Return(output_args ...interface{}) {
	// 输出返回值
	if len(output_args) == 0 {
		this.done_cache_value = nil
		this.ch <- this.done_cache_value
	} else {
		this.done_cache_value = output_args[0]
		this.ch <- this.done_cache_value
	}

	this.is_done = true
	close(this.ch)
}
