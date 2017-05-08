package jsergo

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type JsArray []interface{}

func (this *JsArray) Splice(start_index int, remove_number int, insert_args ...interface{}) *JsArray {
	var removed_slice = JsArray((*this)[start_index : start_index+remove_number])
	this.SpliceWithoutReturn(start_index, remove_number, insert_args...)
	return &removed_slice
}
func (this *JsArray) SpliceWithoutReturn(start_index int, remove_number int, insert_args ...interface{}) {
	var start_slice = JsArray((*this)[:start_index])
	var end_slice = JsArray((*this)[start_index+remove_number:])
	var res = make(JsArray, 0)
	res = append(res, start_slice...)
	res = append(res, insert_args...)
	res = append(res, end_slice...)
	*this = res
}

func (this *JsArray) Join(join_args ...interface{}) string {
	var sep = ","
	if len(join_args) > 0 {
		sep = fmt.Sprint(join_args[0])
	}
	var string_arr = make([]string, len(*this))
	for i, val := range *this {
		string_arr[i] = fmt.Sprint(val)
	}
	return strings.Join(string_arr, sep)
}

func (this *JsArray) Map(map_handle func(interface{}, int, *JsArray) interface{}) *JsArray {
	var res = make(JsArray, len(*this))
	for i, val := range *this {
		var return_item = map_handle(val, i, this)
		res[i] = return_item
	}
	return &res
}
func (this *JsArray) ForEach(map_handle func(interface{}, int, *JsArray)) {
	for i, val := range *this {
		map_handle(val, i, this)
	}
}
func (this *JsArray) Filter(map_handle func(interface{}, int, *JsArray) bool) *JsArray {
	var res = make(JsArray, len(*this))
	var res_len = 0
	for i, val := range *this {
		var return_item = map_handle(val, i, this)
		if return_item {
			res[res_len] = val
			res_len += 1
		}
	}
	res.Length(res_len)
	return &res
}

func (this *JsArray) Push(push_args ...interface{}) int {
	this.PushWithoutReturn(push_args...)
	return len(*this)
}
func (this *JsArray) PushWithoutReturn(push_args ...interface{}) {
	var res = append(*this, push_args...)
	*this = res
}
func (this *JsArray) Length(set_len ...int) int {
	if len(set_len) > 0 {
		var target_len = set_len[0]
		var current_len = len(*this)
		if target_len > current_len {
			var footer_arr = make(JsArray, target_len-current_len)
			this.PushWithoutReturn(footer_arr...)
		} else if target_len < current_len {
			this.SpliceWithoutReturn(target_len, current_len-target_len)
		}
		return target_len
	} else {
		return len(*this)
	}
}
func (this *JsArray) Every(map_handle func(interface{}, int, *JsArray) bool) bool {
	for i, val := range *this {
		var check_every = map_handle(val, i, this)
		if check_every == false {
			return false
		}
	}
	return true
}
func (this *JsArray) Some(map_handle func(interface{}, int, *JsArray) bool) bool {
	for i, val := range *this {
		var check_every = map_handle(val, i, this)
		if check_every == true {
			return true
		}
	}
	return false
}
func (this *JsArray) Fill(fill_val interface{}, args ...int) {
	var start_index, end_index int
	switch len(args) {
	case 0:
		start_index = 0
		end_index = len(*this)
	case 1:
		start_index = args[0]
	default:
		start_index = args[0]
		end_index = args[1]
	}
	if end_index <= start_index {
		return
	}
	for ; start_index < end_index; start_index += 1 {
		(*this)[start_index] = fill_val
	}
}
func (this *JsArray) FindIndex(map_handle func(interface{}, int, *JsArray) bool) int {
	for i, val := range *this {
		var is_finded = map_handle(val, i, this)
		if is_finded {
			return i
		}
	}
	return -1
}
func (this *JsArray) IndexOf(finder interface{}) int {
	finder_type := reflect.TypeOf(finder)
	if finder_type.Comparable() {
		for i, val := range *this {
			val_type := reflect.TypeOf(val)
			// val_value := reflect.ValueOf(val)
			if val_type.Comparable() && val == finder {
				return i
			}
		}
	} else {
		// 直接对比指针
		finder_value := reflect.ValueOf(finder)
		finder_pointer := finder_value.Pointer()
		for i, val := range *this {
			if reflect.TypeOf(val) == finder_type {
				val_pointer := reflect.ValueOf(val).Pointer()
				if finder_pointer == val_pointer {
					return i
				}
			}
		}
	}
	return -1
}

func (this *JsArray) Pop() interface{} {
	var total_len = len(*this)
	defer func() {
		var res = (*this)[:total_len-1]
		*this = res
	}()
	return (*this)[total_len-1]
}
func (this *JsArray) Reduce(reduce_handle func(interface{}, interface{}, int, *JsArray) interface{}, args ...interface{}) interface{} {
	var pre_val interface{}
	if len(args) > 0 {
		pre_val = args[0]
		for i, val := range *this {
			pre_val = reduce_handle(pre_val, val, i, this)
		}
	} else {
		var arr = *this
		pre_val = arr[0]
		for i, total_len := 1, len(arr); i < total_len; i += 1 {
			pre_val = reduce_handle(pre_val, arr[i], i, this)
		}
	}
	return pre_val
}
func (this *JsArray) ReduceRight(reduce_handle func(interface{}, interface{}, int, *JsArray) interface{}, args ...interface{}) interface{} {
	var pre_val interface{}
	var arr = *this
	var start_index = len(arr) - 1
	if len(args) > 0 {
		pre_val = args[0]
	} else {
		pre_val = arr[start_index]
		start_index -= 1
	}
	for i := start_index; i >= 0; i -= 1 {
		pre_val = reduce_handle(pre_val, arr[i], i, this)
	}
	return pre_val
}

func (this *JsArray) ReverseWithoutSet() *JsArray {
	var res = *this
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return &res
}

func (this *JsArray) Reverse() *JsArray {
	var res = this.ReverseWithoutSet()
	*this = *res
	return this
}

func (this *JsArray) Shift() interface{} {
	defer func() {
		var res = (*this)[1:]
		*this = res
	}()
	return (*this)[0]
}
func (this *JsArray) Unshift(push_args ...interface{}) int {
	this.UnshiftWithoutReturn(push_args...)
	return len(*this)
}
func (this *JsArray) UnshiftWithoutReturn(push_args ...interface{}) {
	var res = append(push_args, *this...)
	*this = res
}

func (this *JsArray) Slice(slice_args ...int) *JsArray {
	var start_index, end_index int
	switch len(slice_args) {
	case 0:
		start_index = 0
		end_index = len(*this)
	case 1:
		start_index = slice_args[0]
		end_index = len(*this)
	default:
		start_index = slice_args[0]
		end_index = slice_args[1]
	}
	if end_index <= start_index {
		return &JsArray{}
	}
	var res = (*this)[start_index:end_index]
	return &res
}

type goLangSort struct {
	JsArray     JsArray
	sort_handle func(a, b interface{}) bool
}

func (this *goLangSort) Len() int { return len(this.JsArray) }
func (this *goLangSort) Swap(i, j int) {
	var jsArray = this.JsArray
	jsArray[i], jsArray[j] = jsArray[j], jsArray[i]
}
func (this *goLangSort) Less(i, j int) bool {
	var jsArray = this.JsArray
	return this.sort_handle(jsArray[i], jsArray[j])
}
func (this *JsArray) Sort(sort_handle func(a, b interface{}) int) *JsArray {
	var jsarr_sorter = &goLangSort{
		JsArray: *this,
		sort_handle: func(a, b interface{}) bool {
			return sort_handle(a, b) < 0
		},
	}
	sort.Sort(jsarr_sorter)
	return this
}

func (this *JsArray) FasterSort(sort_handle func(a, b interface{}) bool) *JsArray {
	var jsarr_sorter = &goLangSort{
		JsArray:     *this,
		sort_handle: sort_handle,
	}
	sort.Sort(jsarr_sorter)
	return this
}

func (this *JsArray) Concat(concat_list ...interface{}) *JsArray {
	var res = this.Slice()
	for _, val := range concat_list {
		if jsarray, ok := val.(JsArray); ok {
			res.Push(jsarray...)
		} else {
			res.Push(val)
		}
	}
	return res
}

/*实现其它接口的方法，包括JSON接口等等*/

func (this *JsArray) GoString() string {
	return this.Join()
}

func init() {
	// fmt.Println("")
}
