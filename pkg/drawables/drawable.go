package drawables

import (
	"image/color"

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

func RepeatingPos(pt mat32.Vec3, domain float32) mat32.Vec3 {
	pt.X = modByDomain(pt.X, domain)
	pt.Y = modByDomain(pt.Y, domain)
	pt.Z = modByDomain(pt.Z, domain)
	return pt
}

func modByDomain(in, domain float32) float32 {
	out := mat32.Mod(mat32.Abs(in+(domain/2)), domain)
	out -= domain / 2
	return float32(out)
}
