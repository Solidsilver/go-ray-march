package vec3neon

import (
	"log"

	"github.com/alivanz/go-simd/arm"
	"github.com/alivanz/go-simd/arm/neon"
)

type Vec3Neon struct {
	val arm.Float32X4
}

func New(x, y, z float32) *Vec3Neon {
	vec := new(Vec3Neon)
	vec.val = arm.Float32X4{arm.Float32(x), arm.Float32(y), arm.Float32(z), arm.Float32(0)}
	return vec
}

func NewF64(x, y, z float64) *Vec3Neon {
	vec := new(Vec3Neon)
	vec.val = arm.Float32X4{arm.Float32(x), arm.Float32(y), arm.Float32(z), arm.Float32(0)}
	return vec
}

func (v Vec3Neon) x() arm.Float32 {
	return arm.Float32(v.val[0])
}

func (v Vec3Neon) X() float32 {
	return float32(v.val[0])
}

func (v Vec3Neon) y() arm.Float32 {
	return v.val[1]
}

func (v Vec3Neon) Y() float32 {
	return float32(v.val[1])
}

func (v Vec3Neon) z() arm.Float32 {
	return v.val[2]
}

func (v Vec3Neon) Z() float32 {
	return float32(v.val[2])
}

func main() {
	var a, b arm.Int8X8
	var add, mul arm.Int16X8
	for i := 0; i < 8; i++ {
		a[i] = arm.Int8(i)
		b[i] = arm.Int8(i * i)
	}
	log.Printf("a = %+v", b)
	log.Printf("b = %+v", a)
	neon.VaddlS8(&add, &a, &b)
	neon.VmullS8(&mul, &a, &b)
	log.Printf("add = %+v", add)
	log.Printf("mul = %+v", mul)
}
