package renderer

import "github.com/Solidsilver/go-ray-march/pkg/drawables"

type Scene struct {
	Drawables []drawables.Drawable
	Lights    []drawables.Drawable
	options   LightingOpts
}

func (s *Scene) AddDrawables(draws ...drawables.Drawable) {
	s.Drawables = append(s.Drawables, draws...)
}

func (s *Scene) AddLights(draws ...drawables.Drawable) {
	s.Lights = append(s.Lights, draws...)
}

func NewBlankScene() *Scene {
	return NewScene([]drawables.Drawable{}, []drawables.Drawable{})
}

func NewScene(draws []drawables.Drawable, lights []drawables.Drawable) *Scene {
	return NewSceneWithOpts(DefaultLightingOpts(), draws, lights)
}

func NewSceneWithOpts(opts LightingOpts, draws []drawables.Drawable, lights []drawables.Drawable) *Scene {
	scn := new(Scene)
	scn.Drawables = draws
	scn.Lights = lights
	scn.options = opts
	return scn
}
