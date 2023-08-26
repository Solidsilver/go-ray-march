package drawables

import (
	"image/color"
	"math"

	"github.com/Solidsilver/go-ray-march/pkg/vec3"
	"github.com/Solidsilver/go-ray-march/pkg/vec3neon"
)

type Drawable interface {
	Dist(pt vec3.Vec3) float64
	DistN(pt vec3neon.Vec3Neon) float32
	Color() color.RGBA
	Pos() vec3.Vec3
	Reflectivity() float64
	ID() string
	ReflectionProperties() ReflectionProperties
}

type ReflectionProperties struct {
	// Ambient reflectance is always visible, regardless of lights in a scene.
	Ambient float64
	// Lambertian reflectance is matte reflection directly related to
	// light falling onto the object from light source
	Lambertian float64
	// Specular term ms is the mirror-like reflection of light off an object to the eye
	Specular float64
	// Metalness msm controls the color of the specular highlights. msm = 0 means the highlight is the color
	// of the lightsource, msm = 1 means the highlight is the color of the object.
	Metalness float64
	// msp characterizes the smoothness (i.e., the sharpness of the
	// highlight spot) of a material, and forms an exponent in the calculation of the specular term.
	Smoothness float64

	Reflection float64
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

// func RepeatingPosN(pt vec3neon.Vec3Neon, domain float64) vec3neon.Vec3Neon {
// 	pt.X = modByDomainF32(pt.X, float32(domain))
// 	pt.Y = modByDomainF32(pt.Y, float32(domain))
// 	pt.Z = modByDomainF32(pt.Z, float32(domain))
// 	return pt
// }

func modByDomainF32(in, domain float32) float64 {
	out := math.Mod(math.Abs(float64(in+(domain/2))), float64(domain))
	out -= float64(domain / 2)
	return out
}

func modByDomain(in, domain float64) float64 {
	out := math.Mod(math.Abs(in+(domain/2)), domain)
	out -= domain / 2
	return out
}

func DefaultRefProps() ReflectionProperties {
	return ReflectionProperties{
		Ambient:    0,
		Lambertian: 1,
		Specular:   0,
		Metalness:  1,
		Smoothness: 1,
	}
}
