package vec3

import (
	"math"
)

// Add adds the two vectors component-wise
// and returns the resultant vector
func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{
		v1.X + v2.X,
		v1.Y + v2.Y,
		v1.Z + v2.Z,
	}
}

// Sub subtracts the two vectors component-wise
// and returns the resultant vector
func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{
		v1.X - v2.X,
		v1.Y - v2.Y,
		v1.Z - v2.Z,
	}
}

// Mult multiplies the vector by the given scalar
// and returns the result
func (v Vec3) Plus(num float64) Vec3 {
	v.X += num
	v.Y += num
	v.Z += num
	return v
}

// Sub subtracts the vector by the given scalar
// and returns the resultant vector
func (v Vec3) Minus(num float64) Vec3 {
	v.X -= num
	v.Y -= num
	v.Z -= num
	return v
}

// Multiplies the vector by the given scalar
// and returns the resultant vector
func (v Vec3) Mult(num float64) Vec3 {
	v.X = v.X * num
	v.Y = v.Y * num
	v.Z = v.Z * num
	return v
}

// MultComp multiplies each component
// of the vectors together and returns the result
func (v1 Vec3) MultComp(v2 Vec3) Vec3 {
	v1.X *= v2.X
	v1.Y *= v2.Y
	v1.Z *= v2.Z
	return v1
}

// Div divides each component of the vector by the given scalar
// and returns the resultant vector
func (v Vec3) Div(num float64) Vec3 {
	v.X = v.X / num
	v.Y = v.Y / num
	v.Z = v.Z / num
	return v
}

// Dot returns the dot product of the two vectors
func Dot(v1, v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// Cross returns the cross product of the two vectors
func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	v := Vec3{
		v1.Y*v2.Z - v1.Z*v2.Y,
		v1.Z*v2.X - v1.X*v2.Z,
		v1.X*v2.Y - v1.Y*v2.X,
	}
	return v
}

// Abs returns the absolute value of the vector
// component-wise
func (p Vec3) Abs() Vec3 {
	return Vec3{
		math.Abs(p.X),
		math.Abs(p.Y),
		math.Abs(p.Z),
	}
}

// Returns the vector divided by its norm
func (v Vec3) ToUnit() Vec3 {
	vNorm := v.Norm()
	return v.Div(vNorm)
}

// Calculates the square root of the vector component-wise
func (v Vec3) Sqrt() Vec3 {
	v.X = math.Sqrt(v.X)
	v.Y = math.Sqrt(v.Y)
	v.Z = math.Sqrt(v.Z)
	return v
}
