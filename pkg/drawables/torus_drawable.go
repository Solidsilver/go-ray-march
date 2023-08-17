package drawables

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type Torus struct {
	Center    vec3.Vec3
	Diameters utils.Vec2
	color     color.RGBA
	id        string
}

func NewTorus(pos vec3.Vec3, majorD, minorD float64, color color.RGBA) Torus {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "tor", idNum)
	return Torus{
		Center:    pos,
		Diameters: utils.NewVec2(majorD, minorD),
		color:     color,
		id:        id,
	}
}

func NewNamedTorus(id string, pos vec3.Vec3, majorD, minorD float64, color color.RGBA) Torus {
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

func (t Torus) Color() color.RGBA {
	return t.color
}

func (t Torus) Pos() vec3.Vec3 {
	return t.Center
}

func (t Torus) Equals(d Drawable) bool {
	return t.id == d.ID()
}

func (t Torus) ID() string {
	return t.id
}
