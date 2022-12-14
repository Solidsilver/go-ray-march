package drawables

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/Solidsilver/go-ray-marching/pkg/vec3"
)

type Sphere struct {
	Center vec3.Vec3
	Rad    float64
	color  color.RGBA
	id     string
}

func NewSphere(pos vec3.Vec3, rad float64, color color.RGBA) Sphere {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "sph", idNum)
	return Sphere{pos, rad, color, id}
}

func NewNamedSphere(id string, pos vec3.Vec3, rad float64, color color.RGBA) Sphere {
	return Sphere{pos, rad, color, id}
}

func (s Sphere) Dist(pt vec3.Vec3) float64 {
	vecToSph := pt.Sub(s.Center)
	vecLen := vecToSph.Norm()
	return vecLen - s.Rad
}

func (s Sphere) Color() color.RGBA {
	return s.color
}

func (s Sphere) Pos() vec3.Vec3 {
	return s.Center
}

func (s Sphere) Equals(d Drawable) bool {
	return s.id == d.ID()
}

func (s Sphere) ID() string {
	return s.id
}
