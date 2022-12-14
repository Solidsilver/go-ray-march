package drawables

import (
	"image/color"

	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type Drawable interface {
	Dist(pt vec3.Vec3) float64
	Color() color.RGBA
	Pos() vec3.Vec3
	ID() string
}

func Equals(d1, d2 Drawable) bool {
	if d1 != nil && d2 != nil {
		return d1.ID() == d2.ID()
	}
	return false
}
