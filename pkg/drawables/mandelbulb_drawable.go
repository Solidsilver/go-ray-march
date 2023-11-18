package drawables

import (
	"image/color"
	"math"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type MandelBulb struct {
	Iterations int
	Bailout    float64
	Power      float64
	id         string
	color      color.RGBA
	pos        vec3.Vec3
	repeating  bool
}

func NewMandelB(id string, iter int, bail float64, pow float64, pos vec3.Vec3, color color.RGBA, repeating bool) MandelBulb {
	return MandelBulb{
		iter,
		bail,
		pow,
		id,
		color,
		pos,
		repeating,
	}
}

func (b MandelBulb) Dist(pt vec3.Vec3) float64 {
	if b.repeating {
		pt = RepeatingPos(pt, 10.0)
	}
	z := pt
	dr := 1.0
	r := 0.0

	for i := 0; i < b.Iterations; i++ {
		r = z.Norm()
		if r > b.Bailout {
			break
		}

		// theta := math.Acos(z.Z / r)
		theta := utils.FastAcos(z.Z / r)
		phi := math.Atan2(z.Y, z.X)
		// phi := utils.FastAtan2(z.Y, z.X)
		dr = math.Pow(r, b.Power-1)*float64(b.Power)*dr + 1

		zr := math.Pow(r, b.Power)
		theta = theta * b.Power
		phi = phi * b.Power

		z = vec3.Vec3{
			// X: math.Sin(theta) * math.Cos(phi),
			X: utils.FastSin(theta) * utils.FastCos(phi),
			// Y: math.Sin(phi) * math.Sin(theta),
			Y: utils.FastSin(phi) * utils.FastSin(theta),
			// Z: math.Cos(theta),
			Z: utils.FastCos(theta),
		}.Mult(zr)
		z = z.Add(pt)
	}
	return 0.5 * math.Log(r) * r / dr
}

func (b MandelBulb) Color() color.RGBA {
	return b.color
}

func (b MandelBulb) Pos() vec3.Vec3 {
	return b.pos
}
func (b MandelBulb) ID() string {
	return b.id
}
