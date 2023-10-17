package drawables

import (
	"image/color"

	"goki.dev/mat32/v2"
)

type MandelBulb struct {
	Iterations int
	Bailout    float32
	Power      float32
	id         string
	color      color.RGBA
	pos        mat32.Vec3
	repeating  bool
}

func NewMandelB(id string, iter int, bail float32, pow float32, pos mat32.Vec3, color color.RGBA, repeating bool) MandelBulb {
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
	dr := float32(1.0)
	r := float32(0.0)

	for i := 0; i < b.Iterations; i++ {
		r = z.Length()
		if r > b.Bailout {
			break
		}

		theta := mat32.Acos(z.Z / r)
		phi := mat32.Atan2(z.Y, z.X)
		dr = mat32.Pow(r, b.Power-1)*b.Power*dr + 1

		zr := mat32.Pow(r, b.Power)
		theta = theta * b.Power
		phi = phi * b.Power

		z = mat32.Vec3{
			X: (mat32.Sin(theta) * mat32.Cos(phi)),
			Y: (mat32.Sin(phi) * mat32.Sin(theta)),
			Z: (mat32.Cos(theta)),
		}.MulScalar(zr)
		z = z.Add(pt)
	}
	return 0.5 * mat32.Log(r) * r / dr
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
