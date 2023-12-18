package vec3

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func New(x, y, z float64) *Vec3 {
	vec := new(Vec3)
	vec.X = x
	vec.Y = y
	vec.Z = z
	return vec
}

func NewCopy(v1 Vec3) *Vec3 {
	v := new(Vec3)
	v.X = v1.X
	v.Y = v1.Y
	v.Z = v1.Z
	return v
}

func Zero() Vec3 {
	return Vec3{0, 0, 0}
}

func UnitX() Vec3 {
	return Vec3{1, 0, 0}
}

func UnitY() Vec3 {
	return Vec3{0, 1, 0}
}

func UnitZ() Vec3 {
	return Vec3{0, 0, 1}
}

// Return the euclidean length of the vector
func (v Vec3) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v1 Vec3) Eq(v2 Vec3) bool {
	return v1.X == v2.X && v1.Y == v2.Y && v1.Z == v2.Z
}
