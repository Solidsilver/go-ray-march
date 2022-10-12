package utils

import (
	"math"
)

func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func Sig(x, a, b, c, d float64) float64 {
	return 1/(1-a*math.Exp(b*(x-c))) + d - 0.5
}

func SigLocal(dist float64, step int, smooth float64) float64 {
	b := 1 / smooth
	stepF := float64(step)
	return Sig(dist, 1, b, stepF, stepF)
}

// func Min[T constraints.Ordered](a, b T) T {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }
