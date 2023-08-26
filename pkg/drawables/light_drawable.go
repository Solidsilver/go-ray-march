package drawables

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/Solidsilver/go-ray-march/pkg/vec3"
	"github.com/Solidsilver/go-ray-march/pkg/vec3neon"
)

type Light struct {
	Center    vec3.Vec3
	Rad       float64
	color     color.RGBA
	id        string
	repeating bool
}

func NewLight(pos vec3.Vec3, rad float64, color color.RGBA, repeating bool) Light {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "sph", idNum)
	return Light{pos, rad, color, id, repeating}
}

func NewNamedLight(id string, pos vec3.Vec3, rad float64, color color.RGBA, repeating bool) Light {
	return Light{pos, rad, color, id, repeating}
}

func (s Light) Dist(pt vec3.Vec3) float64 {
	if s.repeating {
		pt = RepeatingPos(pt, 20.0)
	}
	vecToSph := pt.Sub(s.Center)
	vecLen := vecToSph.Norm()
	return vecLen - s.Rad
}

func (s Light) DistN(pt vec3neon.Vec3Neon) float32 {
	// if s.repeating {
	// 	pt = RepeatingPos(pt, 20.0)
	// }
	vecToSph := pt.Sub(vec3neon.FromVec3(s.Center))
	vecLen := vecToSph.Norm()
	return vecLen - float32(s.Rad)
}

func (s Light) Color() color.RGBA {
	return s.color
}

func (s Light) Pos() vec3.Vec3 {
	return s.Center
}

func (s Light) Equals(d Drawable) bool {
	return s.id == d.ID()
}

func (s Light) ID() string {
	return s.id
}

func (s Light) Brightness() uint8 {
	return s.color.A
}

func (s Light) Reflectivity() float64 {
	return 0
}

func (l Light) ReflectionProperties() ReflectionProperties {
	return ReflectionProperties{
		Ambient:    0,
		Lambertian: 0,
		Specular:   0,
		Metalness:  0,
	}
}
