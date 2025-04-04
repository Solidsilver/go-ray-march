package vec3

import (
	"image/color"
	"math"
)

// Run function f on each component of the fector in place
func (v *Vec3) Op(f func(v float64) float64) {
	v.X = f(v.X)
	v.Y = f(v.Y)
	v.Z = f(v.Z)
}

func NewP(x, y, z float64) *Vec3 {
	return &Vec3{
		X: x,
		Y: y,
		Z: z,
	}
}

func NewOfSizeP(size float64) *Vec3 {
	return &Vec3{
		X: size,
		Y: size,
		Z: size,
	}
}

func NewCp(v Vec3) *Vec3 {
	return &Vec3{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}

func (v1 *Vec3) ClampSet(min, max float64) {
	v1.X = Clamp(v1.X, min, max)
	v1.Y = Clamp(v1.Y, min, max)
	v1.Z = Clamp(v1.Z, min, max)
}

// func Clamp(num, min, max float64) float64 {
// 	if num > max {
// 		return max
// 	}
// 	if num < min {
// 		return min
// 	}
// 	return num
// }

func Clamp(num, min, max float64) float64 {
	return math.Max(min, math.Min(max, num))
}

func (v *Vec3) MinSet(v2 *Vec3) {
	v.X = math.Min(v.X, v2.X)
	v.Y = math.Min(v.Y, v2.Y)
	v.Z = math.Min(v.Z, v2.Z)
}

func (v *Vec3) MaxSet(v2 *Vec3) {
	v.X = math.Max(v.X, v2.X)
	v.Y = math.Max(v.Y, v2.Y)
	v.Z = math.Max(v.Z, v2.Z)
}

// RGBAToVec3 converts a color.RGBA to a Vec3
// on the range [0, 1] for each component
// by dividing each component by 255.
// It will ignore the alpha value.
func RGBAToVec3P(color color.RGBA) *Vec3 {
	vect := &Vec3{
		float64(color.R),
		float64(color.G),
		float64(color.B),
	}
	vect.DivSet(255)
	return vect
}

func (v1 *Vec3) Set(v2 *Vec3) {
	v1.X = v2.X
	v1.Y = v2.Y
	v1.Z = v2.Z
}
