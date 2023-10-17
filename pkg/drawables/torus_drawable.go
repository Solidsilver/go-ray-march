package drawables

import (
	"fmt"
	"image/color"
	"math/rand"

	"goki.dev/mat32/v2"
)

type Torus struct {
	Center    mat32.Vec3
	Diameters mat32.Vec2
	color     color.RGBA
	id        string
}

func NewTorus(pos mat32.Vec3, majorD, minorD float32, color color.RGBA) Torus {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "tor", idNum)
	return Torus{
		Center:    pos,
		Diameters: mat32.NewVec2(majorD, minorD),
		color:     color,
		id:        id,
	}
}

func NewNamedTorus(id string, pos mat32.Vec3, majorD, minorD float32, color color.RGBA) Torus {
	return Torus{
		Center:    pos,
		Diameters: mat32.NewVec2(majorD, minorD),
		color:     color,
		id:        id,
	}
}

func (t Torus) Dist(pt mat32.Vec3) float32 {
	vecToTorus := pt.Sub(t.Center)
	q := mat32.NewVec2(vecToTorus.Length()-t.Diameters.X, vecToTorus.Y)
	return q.Length() - t.Diameters.Y
	// return 0.0
}

func (t Torus) Color() color.RGBA {
	return t.color
}

func (t Torus) Pos() mat32.Vec3 {
	return t.Center
}

func (t Torus) Equals(d Drawable) bool {
	return t.id == d.ID()
}

func (t Torus) ID() string {
	return t.id
}
