package drawables

import (
	"fmt"
	"image/color"
	"math/rand"

	"solidsilver.dev/go-ray-marching/pkg/utils"
)

type Sphere struct {
	Center utils.Vec3
	Rad    float64
	color  color.RGBA
	id     string
}

func NewSphere(pos utils.Vec3, rad float64, color color.RGBA) Sphere {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "sph", idNum)
	return Sphere{pos, rad, color, id}
}

func NewNamedSphere(id string, pos utils.Vec3, rad float64, color color.RGBA) Sphere {
	// idNum := rand.Intn(1000)
	// id := fmt.Sprintf("%s-%d", "sph", idNum)
	return Sphere{pos, rad, color, id}
}

func (s Sphere) Dist(pt utils.Vec3) float64 {
	// distToSph :=
	vecToSph := utils.NewSub(pt, s.Center)
	vecLen := vecToSph.Norm()
	return vecLen - s.Rad
}

func (s Sphere) Color() color.RGBA {
	return s.color
}

func (s Sphere) Pos() utils.Vec3 {
	return s.Center
}

// func (s Sphere) Equals(s2 Sphere) bool {
// 	return s.Center.Equals(s2.Center) && s.Rad == s2.Rad && s.color == s2.color
// }

func (s Sphere) Equals(d Drawable) bool {
	return s.id == d.ID()
}

func (s Sphere) ID() string {
	return s.id
}
