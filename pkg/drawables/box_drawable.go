package drawables

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"solidsilver.dev/go-ray-marching/pkg/utils"
)

type Box struct {
	center utils.Vec3
	bounds utils.Vec3
	color  color.RGBA
	id     string
}

func NewBox(pos, bounds utils.Vec3, color color.RGBA) Box {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "sph", idNum)
	return Box{pos, bounds, color, id}
}

func NewNamedBox(id string, pos, bounds utils.Vec3, color color.RGBA) Box {
	return Box{pos, bounds, color, id}
}

func NewCube(pos utils.Vec3, dim float64, color color.RGBA) Box {
	return NewBox(pos, utils.Vec3Size(dim), color)
}

func NewNamedCube(id string, pos utils.Vec3, dim float64, color color.RGBA) Box {
	return NewNamedBox(id, pos, utils.Vec3Size(dim), color)
}

func (b Box) Dist(pt utils.Vec3) float64 {
	shiftPt := *utils.NewSub(pt, b.center)
	q := utils.NewSub(utils.Abs(shiftPt), b.bounds)
	return utils.NewNorm(utils.Vec3Max(*q, utils.Vec3Zero())) + math.Min(utils.Max(q.X, q.Y, q.Z), 0.0)
}

func (b Box) Color() color.RGBA {
	return b.color
}

func (b Box) Pos() utils.Vec3 {
	return b.center
}

func (b Box) ID() string {
	return b.id
}
