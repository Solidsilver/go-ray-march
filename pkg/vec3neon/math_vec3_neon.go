package vec3neon

import (
	"github.com/alivanz/go-simd/arm"
	"github.com/alivanz/go-simd/arm/neon"
)

func (v1 Vec3Neon) Add(v2 Vec3Neon) Vec3Neon {
	newVec := Vec3Neon{}
	neon.VaddqF32(&newVec.val, &v1.val, &v2.val)
	return newVec
}

func (v1 Vec3Neon) Sub(v2 Vec3Neon) Vec3Neon {
	newVec := Vec3Neon{}
	neon.VsubqF32(&newVec.val, &v1.val, &v2.val)
	return newVec
}

func (v Vec3Neon) Plus(num float32) Vec3Neon {
	newVec := Vec3Neon{}
	neon.VaddqF32(&newVec.val, &v.val, &arm.Float32X4{arm.Float32(num), arm.Float32(num), arm.Float32(num), arm.Float32(0)})
	return newVec
}

func (v Vec3Neon) Minus(num float32) Vec3Neon {
	newVec := Vec3Neon{}
	neon.VsubqF32(&newVec.val, &v.val, &arm.Float32X4{arm.Float32(num), arm.Float32(num), arm.Float32(num), arm.Float32(0)})
	return newVec
}

func (v Vec3Neon) Mult(num float32) Vec3Neon {
	newVec := Vec3Neon{}
	neon.VmulqF32(&newVec.val, &v.val, &arm.Float32X4{arm.Float32(num), arm.Float32(num), arm.Float32(num), arm.Float32(0)})
	return newVec
}

func (v1 Vec3Neon) MultComp(v2 Vec3Neon) Vec3Neon {
	newVec := Vec3Neon{}
	neon.VmulqF32(&newVec.val, &v1.val, &v2.val)
	return newVec
}

func (v Vec3Neon) Div(num float32) Vec3Neon {
	newVec := Vec3Neon{}
	neon.VdivqF32(&newVec.val, &v.val, &arm.Float32X4{arm.Float32(num), arm.Float32(num), arm.Float32(num), arm.Float32(0)})
	return newVec
}

func Cross(v1, v2 Vec3Neon) Vec3Neon {
	var result1 arm.Float32X4
	c1 := &arm.Float32X4{v1.y(), v1.z(), v1.x()}
	c2 := &arm.Float32X4{v2.z(), v2.x(), v2.y()}
	neon.VmulqF32(&result1, c1, c2)
	var result2 arm.Float32X4
	c3 := &arm.Float32X4{v1.z(), v1.x(), v1.y()}
	c4 := &arm.Float32X4{v2.y(), v2.z(), v2.x()}
	neon.VmulqF32(&result2, c3, c4)
	neon.VsubqF32(&result1, &result1, &result2)
	return Vec3Neon{result1}
}

func Dot(v1, v2 Vec3Neon) float32 {
	var result arm.Float32X4
	neon.VmulqF32(&result, &v1.val, &v2.val)
	return float32(result[0] + result[1] + result[2])
}


