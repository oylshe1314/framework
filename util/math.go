package util

import "math"

// Round 四舍五入
func Round[T Float](x T, n uint) T {
	var f = math.Pow(10, float64(n))
	return T(float64(int(math.Round(float64(x)*f))) / f)
}
