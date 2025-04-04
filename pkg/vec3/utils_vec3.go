package vec3

import (
	"image/color"
	"math"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
)

// Creates a new Vec3 with each component set to ${num}
func OfSize(num float64) Vec3 {
	return Vec3{
		num,
		num,
		num,
	}
}

// Angle calcualtes the angle between two vectors
// in degrees
func Angle(v1 Vec3, v2 Vec3) float64 {
	val := utils.RadToDeg(math.Acos(Dot(v1, v2) / (v1.Norm() * v2.Norm())))
	return val
}

// Angle_fast is an alternate implementation of Angle
// using a faster but less accurate approximation of acos. See [utils.FastAcos	]
func Angle_fast(v1 Vec3, v2 Vec3) float64 {
	val := utils.RadToDeg(utils.FastAcos(Dot(v1, v2) / (v1.Norm() * v2.Norm())))
	return val
}

// AngleFromUnit calcualtes the angle between two vectors
// in degrees by first converting them to unit vectors
func AngleFromUnit(v1 Vec3, v2 Vec3) float64 {
	v1Unit := v1.ToUnit()
	v2Unit := v2.ToUnit()
	val := utils.RadToDeg(math.Acos(Dot(v1Unit, v2Unit) / (v1Unit.Norm() * v2Unit.Norm())))
	return val
}

// AngleFromUnit_fast is an alternate implementation of AngleFromUnit
// using a faster but less accurate approximation of acos
func AngleFromUnit_fast(v1 Vec3, v2 Vec3) float64 {
	v1Unit := v1.ToUnit()
	v2Unit := v2.ToUnit()
	val := utils.RadToDeg(utils.FastAcos(Dot(v1Unit, v2Unit) / (v1Unit.Norm() * v2Unit.Norm())))
	return val
}

// DirFromPos returns the direction vector pointing from
// one point to another
func DirFromPos(to Vec3, from Vec3) Vec3 {
	dir := to.Sub(from)
	return dir.ToUnit()
}

// Max returns a vector with each component
// set to the maximum of the two vectors' components
func Max(v1, v2 Vec3) Vec3 {
	return Vec3{
		math.Max(v1.X, v2.X),
		math.Max(v1.Y, v2.Y),
		math.Max(v1.Z, v2.Z),
	}
}

// Min returns a vector with each component
// set to the minimum of the two vectors' components
func Min(v1, v2 Vec3) Vec3 {
	return Vec3{
		math.Min(v1.X, v2.X),
		math.Min(v1.Y, v2.Y),
		math.Min(v1.Z, v2.Z),
	}
}

// RGBAToVec3 converts a color.RGBA to a Vec3
// on the range [0, 1] for each component
// by dividing each component by 255.
// It will ignore the alpha value.
func RGBAToVec3(color color.RGBA) Vec3 {
	vect := Vec3{
		float64(color.R),
		float64(color.G),
		float64(color.B),
	}
	return vect.Div(255)
}

// Vec3ToRGBA converts a Vec3 on the range [0, 255] to a color.RGBA
// with the given alpha value by multiplying each component by 255. It will cast
// the components of the vector to uint8.
func Vec3ToRGBA(vec Vec3, a uint8) color.RGBA {
	vec = vec.Mult(255)
	return color.RGBA{
		R: uint8(vec.X),
		G: uint8(vec.Y),
		B: uint8(vec.Z),
		A: a,
	}
}

// Mod returns the vector with each component
// modded by the given number
func (v Vec3) Mod(num float64) Vec3 {
	v.X = math.Mod(v.X, num)
	v.Y = math.Mod(v.Y, num)
	v.Z = math.Mod(v.Z, num)
	return v
}

/*
Returns a vector reflecting the given vector off a surface with the given normal.
*/
func (v Vec3) Reflect(surfaceNormal Vec3) Vec3 {
	return v.Sub(surfaceNormal.Mult(2 * Dot(v, surfaceNormal)))
}

func (v Vec3) Reverse() Vec3 {
	return v.Mult(-1)
}
