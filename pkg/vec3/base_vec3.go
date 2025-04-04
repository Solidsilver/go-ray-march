// Package vec3 provides a 3D vector type and operations on it.
package vec3

import (
	"fmt"
	"math"
)

// A Vec3 is a 3D vector
// with components X, Y, and Z
type Vec3 struct {
	X, Y, Z float64
}

// Creates a new Vec3 with the given components
func New(x, y, z float64) Vec3 {
	return Vec3{
		X: x,
		Y: y,
		Z: z,
	}
}

// (0, 0, 0)
var Zero = Vec3{0, 0, 0}

// (1, 1, 1)
var One = Vec3{1, 1, 1}

// (1, 0, 0)
var UnitX = Vec3{1, 0, 0}

// (0, 1, 0)
var UnitY = Vec3{0, 1, 0}

// (0, 0, 1)
var UnitZ = Vec3{0, 0, 1}

func NewX(x float64) Vec3 {
	return Vec3{
		X: x,
		Y: 0,
		Z: 0,
	}
}

func NewY(y float64) Vec3 {
	return Vec3{
		X: 0,
		Y: y,
		Z: 0,
	}
}

func NewZ(z float64) Vec3 {
	return Vec3{
		X: 0,
		Y: 0,
		Z: z,
	}
}

func Unit() Vec3 {
	return Vec3{1, 1, 1}
}

func Unit() Vec3 {
	return Vec3{1, 1, 1}
}

// Return the euclidean length of the vector
func (v Vec3) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) Mag() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Checks if v1 equals v2 component-wise
func (v1 Vec3) Eq(v2 Vec3) bool {
	return v1.X == v2.X && v1.Y == v2.Y && v1.Z == v2.Z
}

func (v1 Vec3) String() string {
	return fmt.Sprintf("(%f, %f, %f)", v1.X, v1.Y, v1.Z)
}
