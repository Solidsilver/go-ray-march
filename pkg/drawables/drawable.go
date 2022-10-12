package drawables

import (
	"image/color"

	"solidsilver.dev/go-ray-marching/pkg/utils"
)

type Drawable interface {
	Dist(pt utils.Vec3) float64
	Color() color.RGBA
	Pos() utils.Vec3
	ID() string
}

func Equals(d1, d2 Drawable) bool {
	if d1 != nil && d2 != nil {
		return d1.ID() == d2.ID()
	}
	return false
}
