package drawables

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

const E_DIV_2 = math.E / 2

type MandelBulb struct {
	Iterations int
	Bailout    float64
	Power      float64
	id         int64
	color      color.RGBA
	pos        vec3.Vec3
	repeating  bool
	colorVec   vec3.Vec3
	refl       float64
}

func NewMandelB(iter int, bail float64, pow float64, pos vec3.Vec3, color color.RGBA, repeating bool, refl float64) MandelBulb {
	id := rand.Int63()
	return MandelBulb{
		iter,
		bail,
		pow,
		id,
		color,
		pos,
		repeating,
		vec3.RGBAToVec3(color),
		refl,
	}
}

func (b MandelBulb) Dist(pt vec3.Vec3) float64 {
	if b.repeating {
		pt = RepeatingPos(pt, 20.0)
	}
	z := pt
	dr := 1.0
	r := 0.0

	for i := 0; i < b.Iterations; i++ {
		r = z.Norm()
		if r > b.Bailout {
			break
		}

		theta := math.Acos(z.Z / r)
		phi := math.Atan2(z.Y, z.X)
		dr = math.Pow(r, b.Power-1)*float64(b.Power)*dr + 1

		zr := math.Pow(r, b.Power)
		theta = theta * b.Power
		phi = phi * b.Power

		z = vec3.Vec3{
			X: math.Sin(theta) * math.Cos(phi),
			Y: math.Sin(phi) * math.Sin(theta),
			Z: math.Cos(theta),
		}.Mult(zr)
		z = z.Add(pt)
	}
	if r >= math.E && dr == 1.0 {
		return r - E_DIV_2
	} else {
		return 0.5 * math.Log(r) * r / dr
	}
}

func (b MandelBulb) FastDist(pt vec3.Vec3) float64 {
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

		theta := utils.FastAcos(z.Z / r)
		phi := math.Atan2(z.Y, z.X)
		// phi := utils.FastAtan2(z.Y, z.X)
		dr = math.Pow(r, b.Power-1)*b.Power*dr + 1

		zr := math.Pow(r, b.Power)
		theta = theta * b.Power
		phi = phi * b.Power

		z = vec3.Vec3{
			X: utils.FastSin(theta) * utils.FastCos(phi),
			Y: utils.FastSin(phi) * utils.FastSin(theta),
			Z: utils.FastCos(theta),
		}.Mult(zr)
		z = z.Add(pt)
	}
	return 0.5 * utils.FastLog64(r) * r / dr
}

func (b MandelBulb) Color() color.RGBA {
	return b.color
}

func (b MandelBulb) ColorVec() vec3.Vec3 {
	return b.colorVec
}

func (b MandelBulb) Pos() vec3.Vec3 {
	return b.pos
}
func (b MandelBulb) ID() int64 {
	return b.id
}

func (b MandelBulb) IsLight() bool {
	return false
}

func (b MandelBulb) Reflectivity() float64 {
	return b.refl
}
