package drawables

import (
	"image/color"
	"math"

	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type Drawable interface {
	Dist(pt vec3.Vec3) float64
	Color() color.RGBA
	Pos() vec3.Vec3
	Reflectivity() float64
	ID() string
}

func Equals(d1, d2 Drawable) bool {
	if d1 != nil && d2 != nil {
		return d1.ID() == d2.ID()
	}
	return false
}

func RepeatingPos(pt vec3.Vec3, domain float64) vec3.Vec3 {
	pt.X = modByDomain(pt.X, domain)
	pt.Y = modByDomain(pt.Y, domain)
	pt.Z = modByDomain(pt.Z, domain)
	return pt
}

func modByDomain(in, domain float64) float64 {
	out := math.Mod(math.Abs(in+(domain/2)), domain)
	out -= domain / 2
	return out
}
