package drawables

import (
	"image/color"
	"math/rand"

	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type Sphere struct {
	Center    vec3.Vec3
	Rad       float64
	color     color.RGBA
	id        int64
	repeating bool
	invisible bool
	colorVec  vec3.Vec3
}

func NewSphere(pos vec3.Vec3, rad float64, color color.RGBA, repeating bool) Sphere {
	id := rand.Int63()
	return Sphere{pos, rad, color, id, repeating, false, vec3.RGBAToVec3(color)}
}

func NewLight(pos vec3.Vec3, rad float64, color color.RGBA, repeating bool) Sphere {
	id := rand.Int63()
	return Sphere{pos, rad, color, id, repeating, true, vec3.RGBAToVec3(color)}
}

func NewNamedSphere(id int64, pos vec3.Vec3, rad float64, color color.RGBA, invisible bool, repeating bool) Sphere {
	return Sphere{pos, rad, color, id, repeating, invisible, vec3.RGBAToVec3(color)}
}

func (s Sphere) Dist(pt vec3.Vec3) float64 {
	if s.repeating {
		pt = RepeatingPos(pt, 30.0)
	}
	vecToSph := pt.Sub(s.Center)
	vecLen := vecToSph.Norm()
	return vecLen - s.Rad
}

func (s Sphere) FastDist(pt vec3.Vec3) float64 {
	return s.Dist(pt)
}

func (s Sphere) Color() color.RGBA {
	return s.color
}

func (s Sphere) ColorVec() vec3.Vec3 {
	return s.colorVec
}

func (s Sphere) Pos() vec3.Vec3 {
	return s.Center
}

func (s Sphere) Equals(d Drawable) bool {
	return s.id == d.ID()
}

func (s Sphere) ID() int64 {
	return s.id
}

func (s Sphere) IsLight() bool {
	return s.invisible
}
