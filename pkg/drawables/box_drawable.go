package drawables

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec"
	"goki.dev/mat32/v2"
)

type Box struct {
	center mat32.Vec3
	bounds mat32.Vec3
	color  color.RGBA
	id     string
}

func NewBox(pos, bounds mat32.Vec3, color color.RGBA) Box {
	idNum := rand.Intn(1000)
	id := fmt.Sprintf("%s-%d", "sph", idNum)
	return Box{pos, bounds, color, id}
}

func NewNamedBox(id string, pos, bounds mat32.Vec3, color color.RGBA) Box {
	return Box{pos, bounds, color, id}
}

func NewCube(pos mat32.Vec3, dim float32, color color.RGBA) Box {
	return NewBox(pos, mat32.NewVec3Scalar(dim), color)
}

func NewNamedCube(id string, pos mat32.Vec3, dim float32, color color.RGBA) Box {
	return NewNamedBox(id, pos, mat32.NewVec3Scalar(dim), color)
}

func (b Box) Dist(pt mat32.Vec3) float32 {
	shiftPt := pt.Sub(b.center)
	q := shiftPt.Abs().Sub(b.bounds)
	return vec.MaxVec3(q, mat32.Vec3Zero).Length() + float32(math.Min(utils.Max(float64(q.X), float64(q.Y), float64(q.Z)), 0.0))
	// return utils.NewNorm(mat32.Vec3Max(*q, mat32.Vec3Zero())) + math.Min(utils.Max(q.X, q.Y, q.Z), 0.0)
}

func (b Box) Color() color.RGBA {
	return b.color
}

func (b Box) Pos() mat32.Vec3 {
	return b.center
}

func (b Box) ID() string {
	return b.id
}
