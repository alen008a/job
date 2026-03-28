package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}
func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// 字符串是否包含在字符数组里
func IsStringInArray(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

// 字符串MD5加密
func Md5EncodeToString(s string) string {
	hexCode := md5.Sum([]byte(s))
	return hex.EncodeToString(hexCode[:])
}
func InterfaceToInt(inter interface{}) int {
	var target int
	if interInt, ok := inter.(int); ok {
		target = interInt
	} else if targetString, ok := inter.(string); ok {
		temp, err := strconv.Atoi(targetString)
		if err == nil {
			target = temp
		}
	}
	return target
}

/**
* interface类型转成字符串
* @param interface{} 转换的值
* @return string 返回的字符串
 */
func InterfaceToString(inter interface{}) string {

	res := ""
	switch inter := inter.(type) {
	case bool:
		res = fmt.Sprintf("%t", inter)
	case int:
		res = fmt.Sprintf("%d", inter)
	case int64:
		res = fmt.Sprintf("%d", inter)
	case float64:
		res = strconv.FormatFloat(inter, 'f', -1, 64)
	case byte:
		res = fmt.Sprintf("%b", inter)
	case string:
		res = fmt.Sprintf("%s", inter)
	case *bool:
		res = fmt.Sprintf("%p", inter)
	case *int:
		res = fmt.Sprintf("%p", inter)
	case *int64:
		res = fmt.Sprintf("%p", inter)
	case *float64:
		res = fmt.Sprintf("%p", inter)
	case *string:
		res = fmt.Sprintf("%p", inter)
	}
	return res
}

func Overlay(str string, overlay string, start int, end int) string {
	if str == "" {
		return ""
	} else {
		len := len(str)
		if start < 0 {
			start = 0
		}
		if start > len {
			start = len
		}
		if end < 0 {
			end = 0
		}
		if end > len {
			end = len
		}
		if start > end {
			temp := start
			start = end
			end = temp
		}
		return Substring(str, 0, start) + overlay + Substring(str, end, len)
	}
}

func Substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)
	if start < 0 || end > length || start > end {
		return ""
	}
	if start == 0 && end == length {
		return source
	}
	return string(r[start:end])
}

// IsBlank checks if a cs is empty, or nil point or whitespace only.
func IsBlank(cs string) bool {
	if strLen := len(cs); strLen == 0 {
		return true
	} else {
		return len(strings.TrimSpace(cs)) != strLen
	}

}

func ToInt(str string, defalutValue int) (int, error) {
	if str == "" {
		return defalutValue, nil
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// 字符在另外一个字符串里，出现第num次的位置
func OrdinalIndexOf(source, str string, num int) int {
	var r = []rune(source)
	lenStr := len(r)
	variable := -1
	if num <= 0 {
		return variable
	}
	for i := 0; i < lenStr; i++ {
		if string(r[i]) == str {
			variable++
		}
		if variable == num-1 {
			return i
		}
	}
	return variable
}

// 整型是否包含在数组里
func IsIntInArray(target int, strArray []int) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

// FloatPrecision float 精度转换
func FloatPrecision(fStr string, prec int, round bool) (string, error) {
	f, err := strconv.ParseFloat(fStr, 64)
	if err != nil {
		return "", err
	}

	f = Precision(f, prec, round)
	str := strconv.FormatFloat(f, 'f', prec, 64)

	return str, nil
}

// FloatPrecisionStr float 转换为 string 精度转换
func FloatPrecisionStr(f float64, prec int, round bool) string {
	ff := Precision(f, prec, round)
	str := strconv.FormatFloat(ff, 'f', prec, 64)

	return str
}

// FloatPrecision float 精度转换
func FloatFPrecision(fStr string, prec int, round bool) (float64, error) {
	f, err := strconv.ParseFloat(fStr, 64)
	if err != nil {
		return 0, err
	}

	return Precision(f, prec, round), nil
}

// Precision 支持精度以及是否四舍五入, round: true 为四舍五入, false 不是四舍五入
func Precision(f float64, prec int, round bool) float64 {
	// 需要加上对长度的校验, 否则直接用 math.Trunc 会有bug(1.14会变成1.13)
	arr := strings.Split(strconv.FormatFloat(f, 'f', -1, 64), ".")
	if len(arr) < 2 {
		return f
	}
	if len(arr[1]) <= prec {
		return f
	}
	pow10N := math.Pow10(prec)

	if round {
		return math.Trunc((f+0.5/pow10N)*pow10N) / pow10N
	}

	return math.Trunc((f)*pow10N) / pow10N
}

// 字符串转小写去除空格
func Str2Lower(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "", -1)
}
