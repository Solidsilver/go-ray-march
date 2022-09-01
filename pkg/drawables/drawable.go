package drawables

import (
	"image/color"

	"solidsilver.dev/go-ray-marching/pkg/utils"
)

type Drawable interface {
	Dist(pt utils.Vec3) float64
	Color() color.Color
}
