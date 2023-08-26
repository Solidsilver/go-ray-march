package vec3neon

import (
	"image/color"

	"github.com/Solidsilver/go-ray-march/pkg/vec3"
	"github.com/alivanz/go-simd/arm"
	"github.com/alivanz/go-simd/arm/neon"
)

func OfSize(size float32) Vec3Neon {
	return Vec3Neon{arm.Float32X4{arm.Float32(size), arm.Float32(size), arm.Float32(size), arm.Float32(0)}}
}

func Max(v1, v2 Vec3Neon) Vec3Neon {
	neon.VmaxqF32(&v1.val, &v1.val, &v2.val)
	return v1
}

func Min(v1, v2 Vec3Neon) Vec3Neon {
	neon.VminqF32(&v1.val, &v1.val, &v2.val)
	return v1
}

func (v Vec3Neon) Reflect(normal Vec3Neon) Vec3Neon {
	return v.Sub(normal.Mult(2 * Dot(v, normal)))
}

func (v Vec3Neon) Reverse() Vec3Neon {
	return Vec3Neon{arm.Float32X4{-v.x(), -v.y(), -v.z(), arm.Float32(0)}}
}

func (v Vec3Neon) ToVec3() vec3.Vec3 {
	return vec3.Vec3{X: float64(v.x()), Y: float64(v.y()), Z: float64(v.z())}
}

func FromVec3(v vec3.Vec3) Vec3Neon {
	return Vec3Neon{arm.Float32X4{arm.Float32(v.X), arm.Float32(v.Y), arm.Float32(v.Z), arm.Float32(0)}}
}

func RGBAToVecNeon(c color.RGBA) Vec3Neon {
	colorVec := Vec3Neon{arm.Float32X4{arm.Float32(c.R), arm.Float32(c.G), arm.Float32(c.B), arm.Float32(0)}}
	return colorVec.Div(255)
}

func VecNeonToRGBA(vec Vec3Neon, a uint8) color.RGBA {
	vec = vec.Mult(255)
	return color.RGBA{
		R: uint8(vec.x()),
		G: uint8(vec.y()),
		B: uint8(vec.z()),
		A: a,
	}
}
