package drawables

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type Sphere struct {
	Center    vec3.Vec3
	Rad       float64
	color     color.RGBA
	id        string
	repeating bool
	refProps  ReflectionProperties
}

func NewSphere(pos vec3.Vec3, rad float64, color color.RGBA, repeating bool) Sphere {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "sph", idNum)
	return Sphere{pos, rad, color, id, repeating, DefaultRefProps()}
}

func NewNamedSphere(id string, pos vec3.Vec3, rad float64, color color.RGBA, repeating bool, props ...ReflectionProperties) Sphere {
	if len(props) > 0 {
		return Sphere{pos, rad, color, id, repeating, props[0]}
	}
	return Sphere{pos, rad, color, id, repeating, DefaultRefProps()}
}

func (s Sphere) Dist(pt vec3.Vec3) float64 {
	if s.repeating {
		pt = RepeatingPos(pt, 20.0)
	}
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

func (s Sphere) Reflectivity() float64 {
	return 0
}

func (s Sphere) ReflectionProperties() ReflectionProperties {
	return s.refProps
}

func DefaultRefProps() ReflectionProperties {
	return ReflectionProperties{
		Ambient:    0,
		Lambertian: 1,
		Specular:   0,
		Metalness:  1,
		Smoothness: 1,
	}
}
