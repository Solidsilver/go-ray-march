package utils

import (
	"math"

	"golang.org/x/exp/constraints"
)

func RadToDeg[F constraints.Float](rad F) F {
	return rad * 180 / math.Pi
}

func DegToRad[F constraints.Float](deg F) F {
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

/*
Finds the maximum of a list of floats
and returns it.
*/
func MaxN[F constraints.Ordered](nums ...F) F {
	max := nums[0]
	for i, num := range nums {
		if i != 0 && num > max {
			max = num
		}
	}
	return max
}

func Abs(val int) int {
	if val < 0 {
		return val * -1
	}
	return val
}

func Pow2[F constraints.Float](val F) F {
	return val * val
}
