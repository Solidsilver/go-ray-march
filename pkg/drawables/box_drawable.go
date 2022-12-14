package drawables

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/Solidsilver/go-ray-marching/pkg/utils"
	"github.com/Solidsilver/go-ray-marching/pkg/vec3"
)

type Box struct {
	center vec3.Vec3
	bounds vec3.Vec3
	color  color.RGBA
	id     string
}

func NewBox(pos, bounds vec3.Vec3, color color.RGBA) Box {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "sph", idNum)
	return Box{pos, bounds, color, id}
}

func NewNamedBox(id string, pos, bounds vec3.Vec3, color color.RGBA) Box {
	return Box{pos, bounds, color, id}
}

func NewCube(pos vec3.Vec3, dim float64, color color.RGBA) Box {
	return NewBox(pos, vec3.OfSize(dim), color)
}

func NewNamedCube(id string, pos vec3.Vec3, dim float64, color color.RGBA) Box {
	return NewNamedBox(id, pos, vec3.OfSize(dim), color)
}

func (b Box) Dist(pt vec3.Vec3) float64 {
	shiftPt := pt.Sub(b.center)
	q := shiftPt.Abs().Sub(b.bounds)
	return vec3.Max(q, vec3.Zero()).Norm() + math.Min(utils.Max(q.X, q.Y, q.Z), 0.0)
	// return utils.NewNorm(vec3.Vec3Max(*q, vec3.Vec3Zero())) + math.Min(utils.Max(q.X, q.Y, q.Z), 0.0)
}

func (b Box) Color() color.RGBA {
	return b.color
}

func (b Box) Pos() vec3.Vec3 {
	return b.center
}

func (b Box) ID() string {
	return b.id
}
