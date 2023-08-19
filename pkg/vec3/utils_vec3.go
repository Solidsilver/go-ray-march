package vec3

import (
	"image/color"
	"math"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
)

/*
Returns a vector with each component set to the given value.
*/
func OfSize(num float64) Vec3 {
	return Vec3{
		num,
		num,
		num,
	}
}

/*
Returns the degree of the angle between the two vectors.
*/
func Angle(v1 Vec3, v2 Vec3) float64 {
	val := utils.RadToDeg(math.Acos(Dot(v1, v2) / (v1.Norm() * v2.Norm())))
	return val
}

func Angle2(v1 Vec3, v2 Vec3) float64 {
	v1Unit := v1.Unit()
	v2Unit := v2.Unit()
	val := utils.RadToDeg(math.Acos(Dot(v1Unit, v2Unit) / (v1Unit.Norm() * v2Unit.Norm())))
	return val
}

/*
Returns the vector that points from p1 to p2 as a unit vector.
*/
func DirFromPos(p1 Vec3, p2 Vec3) Vec3 {
	dir := p2.Sub(p1)
	return dir.Unit()
}

/*
Returns the maximum value for each component of the two vectors.
*/
func Max(v1, v2 Vec3) Vec3 {
	return Vec3{
		math.Max(v1.X, v2.X),
		math.Max(v1.Y, v2.Y),
		math.Max(v1.Z, v2.Z),
	}
}

/*
Returns the minimum value for each component of the two vectors.
*/
func Min(v1, v2 Vec3) Vec3 {
	return Vec3{
		math.Min(v1.X, v2.X),
		math.Min(v1.Y, v2.Y),
		math.Min(v1.Z, v2.Z),
	}
}

func RGBAToVec3(color color.RGBA) Vec3 {
	vect := Vec3{
		float64(color.R),
		float64(color.G),
		float64(color.B),
	}
	return vect.Div(255)
}

func Vec3ToRGBA(vec Vec3, a uint8) color.RGBA {
	vec = vec.Mult(255)
	return color.RGBA{
		R: uint8(vec.X),
		G: uint8(vec.Y),
		B: uint8(vec.Z),
		A: a,
	}
}

/*
Returns a vector with each component modded by the given number.
*/
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
