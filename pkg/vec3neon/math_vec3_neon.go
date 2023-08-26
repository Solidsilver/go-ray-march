package vec3neon

import (
	"math"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/alivanz/go-simd/arm"
	"github.com/alivanz/go-simd/arm/neon"
)

func (v1 Vec3Neon) Add(v2 Vec3Neon) Vec3Neon {
	neon.VaddqF32(&v1.val, &v1.val, &v2.val)
	return v1
}

func (v1 Vec3Neon) Sub(v2 Vec3Neon) Vec3Neon {
	neon.VsubqF32(&v1.val, &v1.val, &v2.val)
	return v1
}

func (v Vec3Neon) Plus(num float32) Vec3Neon {
	neon.VaddqF32(&v.val, &v.val, &arm.Float32X4{arm.Float32(num), arm.Float32(num), arm.Float32(num), arm.Float32(0)})
	return v
}

func (v Vec3Neon) Minus(num float32) Vec3Neon {
	neon.VsubqF32(&v.val, &v.val, &arm.Float32X4{arm.Float32(num), arm.Float32(num), arm.Float32(num), arm.Float32(0)})
	return v
}

func (v Vec3Neon) Mult(num float32) Vec3Neon {
	neon.VmulqF32(&v.val, &v.val, &arm.Float32X4{arm.Float32(num), arm.Float32(num), arm.Float32(num), arm.Float32(0)})
	return v
}

func (v1 Vec3Neon) MultComp(v2 Vec3Neon) Vec3Neon {
	neon.VmulqF32(&v1.val, &v1.val, &v2.val)
	return v1
}

func (v Vec3Neon) Div(num float32) Vec3Neon {
	neon.VdivqF32(&v.val, &v.val, &arm.Float32X4{arm.Float32(num), arm.Float32(num), arm.Float32(num), arm.Float32(0)})
	return v
}

func (v Vec3Neon) div(num arm.Float32) Vec3Neon {
	neon.VdivqF32(&v.val, &v.val, &arm.Float32X4{arm.Float32(num), arm.Float32(num), arm.Float32(num), arm.Float32(0)})
	return v
}

func Cross(v1, v2 Vec3Neon) Vec3Neon {
	c1 := &arm.Float32X4{v1.y(), v1.z(), v1.x()}
	c2 := &arm.Float32X4{v2.z(), v2.x(), v2.y()}
	neon.VmulqF32(c1, c1, c2)
	c3 := &arm.Float32X4{v1.z(), v1.x(), v1.y()}
	c4 := &arm.Float32X4{v2.y(), v2.z(), v2.x()}
	neon.VmulqF32(c3, c3, c4)
	neon.VsubqF32(c1, c1, c3)
	return Vec3Neon{*c1}
}

func Dot(v1, v2 Vec3Neon) float32 {
	neon.VmulqF32(&v1.val, &v1.val, &v2.val)
	return float32(v1.val[0] + v1.val[1] + v1.val[2])
}

func (v Vec3Neon) Abs() Vec3Neon {
	neon.VabsqF32(&v.val, &v.val)
	return v
}

func (v Vec3Neon) Unit() Vec3Neon {
	return v.div(v.norm())
}

func (v Vec3Neon) Sqrt() Vec3Neon {
	neon.VsqrtqF32(&v.val, &v.val)
	return v
}

func DirFromPos(v1, v2 Vec3Neon) Vec3Neon {
	return v2.Sub(v1).Unit()
}

func Angle(v1, v2 Vec3Neon) float32 {
	rad := Dot(v1, v2) / (v1.Norm() * v2.Norm())
	return utils.RadToDegF32(float32(math.Acos(float64(rad))))
}
