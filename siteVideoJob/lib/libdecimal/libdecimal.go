// 浮点计算库
package libdecimal

import (
	"github.com/shopspring/decimal"
)

// 两个浮点数相加
func FloatAdd(f1, f2 float64) float64 {
	f11 := decimal.NewFromFloat(f1)
	f22 := decimal.NewFromFloat(f2)

	f := f11.Add(f22)
	val, _ := f.Float64()

	return val
}

// 两个浮点数相减
func FloatSub(f1, f2 float64) float64 {
	f11 := decimal.NewFromFloat(f1)
	f22 := decimal.NewFromFloat(f2)

	f := f11.Sub(f22)
	val, _ := f.Float64()

	return val
}

// 两个浮点数相乘
func FloatMul(f1, f2 float64) float64 {
	f11 := decimal.NewFromFloat(f1)
	f22 := decimal.NewFromFloat(f2)

	f := f11.Mul(f22)
	val, _ := f.Float64()

	return val
}

// 两个浮点数相除
func FloatDiv(f1, f2 float64) float64 {
	if f2 == 0 {
		return 0
	}

	f11 := decimal.NewFromFloat(f1)
	f22 := decimal.NewFromFloat(f2)

	f := f11.Div(f22)
	val, _ := f.Float64()

	return val
}

func Cmp(f1, f2 float64) int {
	f11 := decimal.NewFromFloat(f1)
	f22 := decimal.NewFromFloat(f2)

	return f11.Cmp(f22)
}
