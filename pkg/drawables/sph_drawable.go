package drawables

import (
	"fmt"
	"image/color"
	"math/rand"

	"goki.dev/mat32/v2"
)

type Sphere struct {
	Center    mat32.Vec3
	Rad       float32
	color     color.RGBA
	id        string
	repeating bool
}

func NewSphere(pos mat32.Vec3, rad float32, color color.RGBA, repeating bool) Sphere {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "sph", idNum)
	return Sphere{pos, rad, color, id, repeating}
}

func NewNamedSphere(id string, pos mat32.Vec3, rad float32, color color.RGBA, repeating bool) Sphere {
	return Sphere{pos, rad, color, id, repeating}
}

func (s Sphere) Dist(pt mat32.Vec3) float32 {
	if s.repeating {
		pt = RepeatingPos(pt, 20.0)
	}
	vecToSph := pt.Sub(s.Center)
	vecLen := vecToSph.Length()
	return vecLen - float32(s.Rad)
}

func (s Sphere) Color() color.RGBA {
	return s.color
}

func (s Sphere) Pos() mat32.Vec3 {
	return s.Center
}

func (s Sphere) Equals(d Drawable) bool {
	return s.id == d.ID()
}

func (s Sphere) ID() string {
	return s.id
}
