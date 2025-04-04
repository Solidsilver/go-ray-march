package utils

import (
	"math"
)

/*type Vec2 struct {
    X, Y float64
}
*/

type Vec2[T Number] [2]T

func (v *Vec2[T]) X() T {
	return v[0]
}

func (v *Vec2[T]) Y() T {
	return v[1]
}

func NewVec2[T Number](x, y T) Vec2[T] {
	return [2]T{x, y}
}

func (v Vec2[T]) Norm() float64 {
	return math.Sqrt(float64(v[0]*v[0] + v[1]*v[1]))
}

func (v Vec2[T]) Add(Vec2[T]) Vec2[T] {
	return Vec2[T]{v[0] + v[0], v[1] + v[1]}
}

func (v Vec2[T]) Sub(Vec2[T]) Vec2[T] {
	return Vec2[T]{v[0] - v[0], v[1] - v[1]}
}

func (v Vec2[T]) Div2(num T) Vec2[T] {
	v[0] /= num
	v[1] /= num
	return v
}

func Unit2(v Vec2[float64]) Vec2[float64] {
	vNorm := v.Norm()
	return v.Div2(vNorm)
}
