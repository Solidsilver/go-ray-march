package drawables

import (
	"image/color"
	"math"

	"goki.dev/mat32/v2"
)

type MandelBulb struct {
	Iterations int
	Bailout    float64
	Power      float64
	id         string
	color      color.RGBA
	pos        mat32.Vec3
	repeating  bool
}

func NewMandelB(id string, iter int, bail float64, pow float64, pos mat32.Vec3, color color.RGBA, repeating bool) MandelBulb {
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

func (b MandelBulb) Dist(pt mat32.Vec3) float32 {
	if b.repeating {
		pt = RepeatingPos(pt, 10.0)
	}
	z := pt
	dr := 1.0
	r := 0.0

	for i := 0; i < b.Iterations; i++ {
		r = float64(z.Length())
		if r > b.Bailout {
			break
		}

		theta := math.Acos(float64(z.Z) / r)
		phi := math.Atan2(float64(z.Z), float64(z.X))
		dr = math.Pow(r, b.Power-1)*float64(b.Power)*dr + 1

		zr := math.Pow(r, b.Power)
		theta = theta * b.Power
		phi = phi * b.Power

		z = mat32.Vec3{
			X: float32(math.Sin(theta) * math.Cos(phi)),
			Y: float32(math.Sin(phi) * math.Sin(theta)),
			Z: float32(math.Cos(theta)),
		}.MulScalar(float32(zr))
		z = z.Add(pt)
	}
	return float32(0.5 * math.Log(r) * r / dr)
}

func (b MandelBulb) Color() color.RGBA {
	return b.color
}

func (b MandelBulb) Pos() mat32.Vec3 {
	return b.pos
}
func (b MandelBulb) ID() string {
	return b.id
}
