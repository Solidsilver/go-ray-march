package drawables

import (
	"image/color"
	"math"

	"goki.dev/mat32/v2"
)

type Drawable interface {
	Dist(pt mat32.Vec3) float32
	Color() color.RGBA
	Pos() mat32.Vec3
	ID() string
}

func Equals(d1, d2 Drawable) bool {
	if d1 != nil && d2 != nil {
		return d1.ID() == d2.ID()
	}
	return false
}

func RepeatingPos(pt mat32.Vec3, domain float64) mat32.Vec3 {
	pt.X = modByDomain(float64(pt.X), domain)
	pt.Y = modByDomain(float64(pt.Y), domain)
	pt.Z = modByDomain(float64(pt.Z), domain)
	return pt
}

func modByDomain(in, domain float64) float32 {
	out := math.Mod(math.Abs(in+(domain/2)), domain)
	out -= domain / 2
	return float32(out)
}
