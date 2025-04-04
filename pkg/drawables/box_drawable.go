package drawables

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type Box struct {
	center vec3.Vec3
	bounds vec3.Vec3
	color  color.RGBA
	id     int64
}

func NewBox(pos, bounds vec3.Vec3, color color.RGBA) Box {
	id := rand.Int63()
	return Box{pos, bounds, color, id}
}

func NewNamedBox(id int64, pos, bounds vec3.Vec3, color color.RGBA) Box {
	return Box{pos, bounds, color, id}
}

func NewCube(pos vec3.Vec3, dim float64, color color.RGBA) Box {
	return NewBox(pos, vec3.OfSize(dim), color)
}

func NewNamedCube(id int64, pos vec3.Vec3, dim float64, color color.RGBA) Box {
	return NewNamedBox(id, pos, vec3.OfSize(dim), color)
}

func (b Box) Dist(pt vec3.Vec3) float64 {
	shiftPt := pt.Sub(b.center)
	q := shiftPt.Abs().Sub(b.bounds)
	return vec3.Max(q, vec3.Zero).Norm() + math.Min(utils.MaxN(q.X, q.Y, q.Z), 0.0)
}

func (b Box) FastDist(pt vec3.Vec3) float64 {
	return b.Dist(pt)
}

func (b Box) Color() color.RGBA {
	return b.color
}

func (b Box) ColorVec() vec3.Vec3 {
	return vec3.RGBAToVec3(b.color)
}

func (b Box) Pos() vec3.Vec3 {
	return b.center
}

func (b Box) ID() int64 {
	return b.id
}

func (b Box) IsLight() bool {
	return false
}
