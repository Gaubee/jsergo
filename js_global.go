package jsergo

import (
	// "errors"
	"fmt"
	"math"
	"reflect"
)

var num_table = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'z', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
var lower_table = [...]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'z', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
var upper_table = [...]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'Z', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

var upper_len = len(upper_table)
var upper_map = make(map[byte]int)

func init() {
	for i := 0; i < upper_len; i += 1 {
		upper_map[upper_table[i]] = i
	}
}

func get_char_lower(_char byte) byte {
	var _char_num = upper_map[_char]
	if 0 <= _char_num && _char_num < upper_len && upper_table[_char_num] == _char {
		return lower_table[_char_num]
	}
	return _char
}

func upper_to_lower(str string) string {
	var res = ""
	for i, len := 0, len(str); i < len; i += 1 {
		res += string(get_char_lower(str[i]))
	}
	return res
}

/*
 * trim-left
 */
var whitespace_table = " \n\r\t\f\x0b\xa0\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u200b\u2028\u2029\u3000"

//[...]byte{' ', '\n', '\r', '\t', '\f', '\x0b', '\xa0', '\u2000', '\u2001', '\u2002', '\u2003', '\u2004', '\u2005', '\u2006', '\u2007', '\u2008', '\u2009', '\u200a', '\u200b', '\u2028', '\u2029', '\u3000'}
var whitespace_len = len(whitespace_table)
var whitespace_map = make(map[byte]int)

func init() {
	for i := 0; i < whitespace_len; i += 1 {
		whitespace_map[whitespace_table[i]] = i
	}
}
func whitespace_has(_char byte) bool {
	var _char_num = whitespace_map[_char]
	if 0 <= _char_num && _char_num < whitespace_len && whitespace_table[_char_num] == _char {
		return true
	}
	return false
}

func left_trim_index(str string) int {
	i, len := 0, len(str)
	for ; i < len; i += 1 {
		if !whitespace_has(str[i]) {
			return i
		}
	}
	return -1
}

func left_trim_with_zero_index(str string) int {
	i, len := 0, len(str)
	for ; i < len; i += 1 {
		c := str[i]
		if !whitespace_has(c) && c != '0' {
			return i
		}
	}
	return -1
}

func (this *Global) IsNaN(v JsNumber) bool {
	return math.IsNaN(float64(v))
}

var num_map_cache = make(map[int]map[byte]int)

func getNumMap(radix int) map[byte]int {
	if num_map, ok := num_map_cache[radix]; ok {
		return num_map
	}
	var num_map = make(map[byte]int)
	for i := 0; i < radix; i += 1 {
		num_map[num_table[i]] = i
	}
	num_map_cache[radix] = num_map
	return num_map
}

type Global JsObject

var global = &Global{}

var floatType = reflect.TypeOf(float64(0))

func (this *Global) ParseInt(to_num interface{}, args ...interface{}) JsNumber {
	v := reflect.ValueOf(to_num)
	v = reflect.Indirect(v)
	if v.Type().ConvertibleTo(floatType) {
		fv := v.Convert(floatType)
		num := fv.Float()
		return JsNumber(int64(num))
	}

	var numString = fmt.Sprint(to_num)
	// fmt.Println("numString", numString)
	var radix = 0
	switch len(args) {
	case 0:
	default:
		if r, ok := args[0].(JsNumber); ok {
			if this.IsNaN(r) {
				radix = 0
			} else {
				radix = int(r)
			}
		} else {
			var parsed_radix = this.ParseInt(args[0])
			if this.IsNaN(parsed_radix) {
				radix = 0
			} else {
				radix = int(parsed_radix)
			}
		}
	}
	// fmt.Println("numString, radix", numString, radix)
	return toInt(numString, radix)
}

func toInt(value string, radix int) JsNumber {

	var value_str = upper_to_lower(value)

	/*
	 * 取得正确的进制数
	 * NaN: 0x*(16) , *(10)
	 * 16:0x*
	 */
	var start_index = left_trim_index(value_str)
	if start_index < 0 {
		return NaN
	}
	var base = 1
	if value_str[start_index] == '-' {
		base = -1
		start_index += 1
	}
	var value_str_len = len(value_str)
	if (radix == 0 || radix == 16) && start_index <= value_str_len-2 && value_str[start_index] == '0' && value_str[start_index+1] == 'x' {
		radix = 16
		start_index += 2 //忽略0x
	} else if radix == 0 {
		radix = 10
	}
	if radix < 2 || radix > 36 { // NaN, 0~36
		return NaN
	}

	/*
	 * 每一个字符对应的数字，临时生成，使用缓存的话可能会被通过Object.prototype来注入缓存，导致安全问题
	 */
	var res float64 = 0

	var num_map = getNumMap(radix)
	i := start_index
	for ; i < value_str_len; i += 1 {
		var c = value_str[i]
		var c_num = num_map[c]
		if 0 <= c_num && c_num < radix && num_table[c_num] == c {
			res = res*float64(radix) + float64(c_num)
		} else {
			break
		}
	}
	//如果第一个字符就是非法字符（包括空字符）的话，等于res是空的，返回NaN , 注意：'0x'返回NaN，因为被当成16进制处理
	if i == start_index {
		return NaN
	} else {
		return JsNumber(res * float64(base))
	}
}
func (this *Global) ParseFloat(to_num interface{}) JsNumber {
	v := reflect.ValueOf(to_num)
	v = reflect.Indirect(v)
	if v.Type().ConvertibleTo(floatType) {
		fv := v.Convert(floatType)
		num := fv.Float()
		return JsNumber(num)
	}

	var numString = fmt.Sprint(to_num)

	return toFloat(numString)
}
func toFloat(value string) JsNumber {
	var value_str = upper_to_lower(value)

	var start_index = left_trim_index(value_str)
	if start_index < 0 {
		return NaN
	}
	var base = 1
	if value_str[start_index] == '-' {
		base = -1
		start_index += 1
	}
	var value_str_len = len(value_str)
	var radix = 10

	var res float64 = 0
	var num_map = getNumMap(radix)
	i := start_index

	var float_dot_index = -1
	for ; i < value_str_len; i += 1 {
		var c = value_str[i]
		if c == '.' {
			if float_dot_index == -1 {
				float_dot_index = i
				continue
			} else {
				break
			}
		}
		var c_num = num_map[c]

		if 0 <= c_num && c_num < radix && num_table[c_num] == c {
			res = res*float64(radix) + float64(c_num)
		} else {
			break
		}
	}
	//如果第一个字符就是非法字符（包括空字符）的话，等于res是空的，返回NaN , 注意：'0x'返回NaN，因为被当成16进制处理
	if i == start_index {
		return NaN
	} else {
		res *= float64(base)
		if float_dot_index != -1 {
			return JsNumber(res / math.Pow10(i-float_dot_index-1))
		} else {
			return JsNumber(res)
		}
	}
}
