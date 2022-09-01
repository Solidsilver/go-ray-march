package drawables

import (
	"image/color"

	"solidsilver.dev/go-ray-marching/pkg/utils"
)

type Sphere struct {
	Center utils.Vec3
	Rad    float64
	color  color.Color
}

func NewSphere(pos utils.Vec3, rad float64, color color.Color) Sphere {
	return Sphere{pos, rad, color}
}

func (s Sphere) Dist(pt utils.Vec3) float64 {
	// distToSph :=
	vecToSph := utils.NewSub(pt, s.Center)
	vecLen := vecToSph.Norm()
	return vecLen - s.Rad
}

func (s Sphere) Color() color.Color {
	return s.color
}
