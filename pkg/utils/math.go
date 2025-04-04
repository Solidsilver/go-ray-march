package utils

import (
	"math"

	"golang.org/x/exp/constraints"
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

// Max returns the maximum number
func Max[T Number](n1, n2 T) T {
	if n1 > n2 {
		return n1
	}
	return n2
}

// MaxN returns the maximum number
// in a list of numbers
func MaxN[T Number](nums ...T) T {
	max := nums[0]
	for i, num := range nums {
		if i != 0 {
			max = Max(max, num)
		}
	}
	return max
}

func Abs[T SignedNumber](val T) T {
	if val < 0 {
		return val * -1.0
	}
	return val
}

type SignedNumber interface {
	constraints.Signed | constraints.Float
}

type Number interface {
	constraints.Integer | constraints.Float
}
