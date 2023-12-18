package renderer

import (
	"fmt"
	"image/color"
	"math"
)

type VignetteOpts struct {
	enabled  bool
	strength float64
}

func (vo VignetteOpts) String() string {
	return fmt.Sprintf("vignett: {enabled: %t, strength: %f}", vo.enabled, vo.strength)
}

type BGOpts struct {
	color color.RGBA
	show  bool
}

func (bg BGOpts) String() string {
	return fmt.Sprintf("bg: {color: %s, show: %t}", ColorString(bg.color), bg.show)
}

type AmbientOcclusionOpts struct {
	enabled  bool
	inverted bool
	maxSteps float64
}

func (ao AmbientOcclusionOpts) String() string {
	return fmt.Sprintf("ambient_occ: {enabled: %t, inverted: %t, maxSteps: %f}", ao.enabled, ao.inverted, ao.maxSteps)
}

type DropoffOpts struct {
	enabled  bool
	color    color.RGBA
	distance float64
}

func (do DropoffOpts) String() string {
	return fmt.Sprintf("dropoff: {enabled: %t, color: %s, distance: %f}", do.enabled, ColorString(do.color), do.distance)
}

func ColorString(c color.RGBA) string {
	return fmt.Sprintf("rgba(%d, %d, %d, %d)", c.R, c.G, c.B, c.A)
}

type TraceOpts struct {
	LOD        bool
	minHitDist float64
	maxHitDist float64
	maxSteps   int
	maxDist    float64
}

func (to TraceOpts) String() string {
	return fmt.Sprintf("trace: {LOD: %t, minHitDist: %f, maxHitDist: %f, maxSteps: %d, maxDist: %f}", to.LOD, to.minHitDist, to.maxHitDist, to.maxSteps, to.maxDist)
}

type LightingOpts struct {
	shadows  bool
	vignette VignetteOpts
	bg       BGOpts
	ao       AmbientOcclusionOpts
	dropoff  DropoffOpts
	trace    TraceOpts
}

func DefaultLightingOpts() LightingOpts {
	maxTraceDist := 5000.0
	minHitDist := 0.0001
	lopts := LightingOpts{
		shadows: true,
		vignette: VignetteOpts{
			enabled:  false,
			strength: 0.05,
		},
		bg: BGOpts{
			color: BG_COLOR,
			show:  true,
		},
		ao: AmbientOcclusionOpts{
			enabled:  true,
			inverted: false,
			maxSteps: math.Sqrt(maxTraceDist*10)/10 + 150,
		},
		dropoff: DropoffOpts{
			enabled: true,
			color: color.RGBA{
				0, 0, 0, 0,
			},
			distance: maxTraceDist / 25,
		},
		trace: TraceOpts{
			LOD:        true,
			minHitDist: minHitDist,
			maxHitDist: 10.0,
			maxSteps:   100000,
			maxDist:    maxTraceDist,
		},
	}

	return lopts
}

func (lopts LightingOpts) String() string {
	return fmt.Sprintf("LightingOpts{shadows: %t, vignette: %s, bg: %s, ao: %s, dropoff: %s, trace: %s}", lopts.shadows, lopts.vignette, lopts.bg, lopts.ao, lopts.dropoff, lopts.trace)
}

func (lopts LightingOpts) JsonString() string {
	return fmt.Sprintf("LightingOpts{shadows: %t, vignette: %s, bg: %s, ao: %s, dropoff: %s, trace: %s}", lopts.shadows, lopts.vignette, lopts.bg, lopts.ao, lopts.dropoff, lopts.trace)
}
