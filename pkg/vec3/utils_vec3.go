package vec3

import (
	"math"

	"solidsilver.dev/go-ray-marching/pkg/utils"
)

func OfSize(num float64) Vec3 {
	return Vec3{
		num,
		num,
		num,
	}
}

func Angle(v1 Vec3, v2 Vec3) float64 {
	val := utils.RadToDeg(math.Acos(Dot(v1, v2) / (v1.Norm() * v2.Norm())))
	return val
}

func Angle2(v1 Vec3, v2 Vec3) float64 {
	v1Unit := v1.Unit();
	v2Unit := v2.Unit();
	val := utils.RadToDeg(math.Acos(Dot(v1Unit, v2Unit) / (v1Unit.Norm() * v2Unit.Norm())))
	return val
}
func DirFromPos(p1 Vec3, p2 Vec3) Vec3 {
	dir := p1.Sub(p2)
	return dir.Unit()
}

func Max(v1, v2 Vec3) Vec3 {
	return Vec3{
		math.Max(v1.X, v2.X),
		math.Max(v1.Y, v2.Y),
		math.Max(v1.Z, v2.Z),
	}
}
