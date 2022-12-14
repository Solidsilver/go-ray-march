package renderer

import "github.com/Solidsilver/go-ray-marching/pkg/drawables"

type Scene struct {
	Drawables []drawables.Drawable
	Lights    []drawables.Drawable
}

func NewScene(draws []drawables.Drawable, lights []drawables.Drawable) *Scene {
	scn := new(Scene)
	scn.Drawables = draws
	scn.Lights = lights
	return scn
}
