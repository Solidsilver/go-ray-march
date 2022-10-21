package utils

import "math"


/*type Vec2 struct {
    X, Y float64
}
*/



type Vec2 [2]float64

func (v Vec2) X() float64 {
    return v[0]
}

func (v Vec2) Y() float64 {
    return v[1]
}

func NewVec2(x, y float64) Vec2 {
    return [2]float64{x, y}
}


func (v Vec2) Norm() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1])
}

func (v Vec2) Div2(num float64) Vec2 {
    v[0] /= num
    v[1] /= num
    return v
}

func (v Vec2) Unit2() Vec2 {
	vNorm := v.Norm()
	return v.Div2(vNorm)
}