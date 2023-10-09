package utils

import (
	"math"

	"golang.org/x/exp/constraints"
)

/*type Vec2 struct {
    X, Y float64
}
*/

type Vec2[F constraints.Float] [2]F

func (v Vec2[F]) X() F {
	return v[0]
}

func (v Vec2[F]) Y() F {
	return v[1]
}

func NewVec2[F constraints.Float](x, y F) Vec2[F] {
	return [2]F{x, y}
}

func (v Vec2[F]) Norm() F {
	return F(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1])))
}

func (v Vec2[F]) Div2(num F) Vec2[F] {
	v[0] /= num
	v[1] /= num
	return v
}

func (v Vec2[F]) Unit2() Vec2[F] {
	vNorm := v.Norm()
	return v.Div2(vNorm)
}
