package jsergo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_array_Split(t *testing.T) {
	var arr = make(JsArray, 5)
	arr[0] = 1
	arr[1] = 2
	arr[2] = 3
	arr[3] = 4
	arr[4] = 5
	assert.Equal(t, fmt.Sprint(arr), "[1 2 3 4 5]", "before splice arr content")
	var removed_content = arr.Splice(1, 1, 20, 21)
	assert.Equal(t, fmt.Sprint(arr), "[1 20 21 3 4 5]", "after splice arr content")
	assert.Equal(t, fmt.Sprint(*removed_content), "[2]", "removed content")
}

func Test_Push(t *testing.T) {
	var arr = make(JsArray, 0)
	var arr_len = arr.Push(1, 2, 3)
	assert.Equal(t, fmt.Sprint(arr), "[1 2 3]", "pushed arr content")
	assert.Equal(t, arr_len, 3, "pushed arr length")
}

func Test_Join(t *testing.T) {
	var arr = make(JsArray, 0)
	arr.Push(1, "QAQ", 3)
	assert.Equal(t, arr.Join(), "1,QAQ,3", "arr joined with ,")
	assert.Equal(t, arr.Join("|"), "1|QAQ|3", "arr joined with |")
}

func Test_Map(t *testing.T) {
	var arr = JsArray{1, "QAQ", 3}
	var res_arr = arr.Map(func(v interface{}, i int, _ *JsArray) interface{} {
		return fmt.Sprint(v) + " " + fmt.Sprint(i)
	})
	assert.Equal(t, res_arr.Join(), "1 0,QAQ 1,3 2", "arr Map")
}

func Test_Length(t *testing.T) {
	var arr = JsArray{nil, nil}
	assert.Equal(t, arr.Length(), 2, "arr get length")
	assert.Equal(t, arr.Length(5), 5, "arr set length")
	assert.Equal(t, arr.Length(), 5, "arr get length")
}
func Test_Every(t *testing.T) {
	var arr = JsArray{-1, -2}
	var every_handler = func(v interface{}, _ int, _ *JsArray) bool {
		return v.(int) < 0
	}
	assert.Equal(t, arr.Every(every_handler), true, "arr every item < 0")

	arr.Push(1)
	assert.Equal(t, arr.Every(every_handler), false, "arr no every item < 0")
}

func Test_Some(t *testing.T) {
	var arr = JsArray{-1, -2}
	var some_handler = func(v interface{}, _ int, _ *JsArray) bool {
		return v.(int) > 0
	}
	assert.Equal(t, arr.Some(some_handler), false, "arr no some item > 0")

	arr.Push(1)
	assert.Equal(t, arr.Some(some_handler), true, "arr has some item >= 0")
}

func Test_Fill(t *testing.T) {
	var arr = JsArray{0, 0, 0}
	arr.Fill(1)
	assert.Equal(t, arr.Join(), "1,1,1")
	arr.Fill(-1, 1, 1)
	assert.Equal(t, arr.Join(), "1,1,1")
	arr.Fill(-1, 1, 2)
	assert.Equal(t, arr.Join(), "1,-1,1")
}

func Test_Filter(t *testing.T) {
	var arr = JsArray{-1, -2, 0, 1, 2}
	var filter_handler = func(v interface{}, _ int, _ *JsArray) bool {
		return v.(int) >= 0
	}
	assert.Equal(t, arr.Filter(filter_handler).Join(), "0,1,2")
}
func Test_FindIndex(t *testing.T) {
	var arr = JsArray{-1, -2, 0, 1, 2}
	var findIndex_handler = func(v interface{}, _ int, _ *JsArray) bool {
		return v.(int) == 1
	}
	assert.Equal(t, arr.FindIndex(findIndex_handler), 3)
}

func Test_Pop(t *testing.T) {
	var arr = JsArray{1, 2, 3}
	assert.Equal(t, arr.Pop(), 3)
	assert.Equal(t, arr.Join(), "1,2")
}

func Test_Reverse(t *testing.T) {
	var arr = JsArray{1, 2, 3}
	assert.Equal(t, arr.Reverse().Join(), "3,2,1")
}

func Test_Reduce_Right(t *testing.T) {
	var arr = JsArray{1, 2, 3}
	var reduce_handler = func(pre_val interface{}, cur_val interface{}, i int, _ *JsArray) interface{} {
		return (pre_val.(int) + i) * cur_val.(int)
	}
	assert.Equal(t, arr.Reduce(reduce_handler, 0), 12, "reduce with initialValue")
	assert.Equal(t, arr.Reduce(reduce_handler), 18, "reduce without initialValue")
	assert.Equal(t, arr.ReduceRight(reduce_handler, 0), 14, "reduceRight with initialValue")
	assert.Equal(t, arr.ReduceRight(reduce_handler), 8, "reduceRight without initialValue")
}

func Test_Shift(t *testing.T) {
	var arr = JsArray{1, 2, 3}
	assert.Equal(t, arr.Shift(), 1)
	assert.Equal(t, arr.Join(), "2,3")
}

func Test_Slice(t *testing.T) {
	var arr = JsArray{1, 2, 3}
	assert.Equal(t, arr.Slice(1).Join(), "2,3")
	assert.Equal(t, arr.Slice(1, 2).Join(), "2")
	assert.Equal(t, arr.Slice(1, 1).Join(), "")
	assert.Equal(t, arr.Slice().Join(), "1,2,3")
}

func Test_Sort(t *testing.T) {
	var arr = JsArray{11, 22, 33, 14, 25, 36, 17, 28, 39, 110, 211, 312, 113, 214, 17, 3}
	arr.Sort(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})
	assert.Equal(t, arr.Join(" "), "3 11 14 17 17 22 25 28 33 36 39 110 113 211 214 312")
}

func Test_FasterSort(t *testing.T) {
	var arr = JsArray{11, 22, 33, 14, 25, 36, 17, 28, 39, 110, 211, 312, 113, 214, 17, 3}
	arr.FasterSort(func(a, b interface{}) bool {
		return a.(int) < b.(int)
	})
	assert.Equal(t, arr.Join(" "), "3 11 14 17 17 22 25 28 33 36 39 110 113 211 214 312")
}

func Test_Concat(t *testing.T) {
	var arr = JsArray{1, 2, 3}
	assert.Equal(t, arr.Concat().Join(), "1,2,3")
	assert.Equal(t, arr.Concat(JsArray{5}).Join(), "1,2,3,5")
	assert.Equal(t, arr.Concat(JsArray{5}, 6).Join(), "1,2,3,5,6")
	assert.Equal(t, arr.Concat(JsArray{5}, 6, JsArray{JsArray{7}}).Join(), "1,2,3,5,6,[7]")
}
