package drawables

import (
	"image/color"
	"math/rand"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type Torus struct {
	Center    vec3.Vec3
	Diameters utils.Vec2[float64]
	color     color.RGBA
	id        int64
}

func NewTorus(pos vec3.Vec3, majorD, minorD float64, color color.RGBA) Torus {

	return Torus{
		Center:    pos,
		Diameters: utils.NewVec2(majorD, minorD),
		color:     color,
		id:        rand.Int63(),
	}
}

func NewNamedTorus(id int64, pos vec3.Vec3, majorD, minorD float64, color color.RGBA) Torus {
	return Torus{
		Center:    pos,
		Diameters: utils.NewVec2(majorD, minorD),
		color:     color,
		id:        id,
	}
}

func (t Torus) Dist(pt vec3.Vec3) float64 {
	vecToTorus := pt.Sub(t.Center)
	q := utils.NewVec2(vecToTorus.Norm()-t.Diameters.X(), vecToTorus.Y)
	return q.Norm() - t.Diameters.Y()
	// return 0.0
}

func (t Torus) FastDist(pt vec3.Vec3) float64 {
	return t.Dist(pt)
}

func (t Torus) Color() color.RGBA {
	return t.color
}

func (t Torus) ColorVec() vec3.Vec3 {
	return vec3.RGBAToVec3(t.color)
}

func (t Torus) Pos() vec3.Vec3 {
	return t.Center
}

func (t Torus) ID() int64 {
	return t.id
}

func (t Torus) IsLight() bool {
	return false
}
