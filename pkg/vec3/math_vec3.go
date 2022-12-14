package vec3

import (
	"math"
)

func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{
		v1.X + v2.X,
		v1.Y + v2.Y,
		v1.Z + v2.Z,
	}
}

func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{
		v1.X - v2.X,
		v1.Y - v2.Y,
		v1.Z - v2.Z,
	}
}

func (v Vec3) Plus(num float64) Vec3 {
	v.X += num
	v.Y += num
	v.Z += num
	return v
}

func (v Vec3) Minus(num float64) Vec3 {
	v.X -= num
	v.Y -= num
	v.Z -= num
	return v
}

func (v Vec3) Mult(num float64) Vec3 {
	v.X = v.X * num
	v.Y = v.Y * num
	v.Z = v.Z * num
	return v
}

func (v1 Vec3) MultComp(v2 Vec3) Vec3 {
	v1.X *= v2.X
	v1.Y *= v2.Y
	v1.Z *= v2.Z
	return v1
}

func (v Vec3) Div(num float64) Vec3 {
	v.X = v.X / num
	v.Y = v.Y / num
	v.Z = v.Z / num
	return v
}

func Dot(v1, v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	v := Vec3{
		v1.Y*v2.Z - v1.Z*v2.Y,
		v1.Z*v2.X - v1.X*v2.Z,
		v1.X*v2.Y - v1.Y*v2.X,
	}
	return v
}

func (p Vec3) Abs() Vec3 {
	return Vec3{
		math.Abs(p.X),
		math.Abs(p.Y),
		math.Abs(p.Z),
	}
}

func (v Vec3) Unit() Vec3 {
	vNorm := v.Norm()
	return v.Div(vNorm)
}
