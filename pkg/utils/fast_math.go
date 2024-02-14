package utils

import (
	"math"
	"unsafe"
)

func FastAcos(x float64) float64 {
	return (-0.69813170079773212*x*x-0.87266462599716477)*x + 1.5707963267948966
}

// ApproximateAtan2 computes the approximate arctangent of y/x.
func FastAtan2(y, x float64) float64 {
	// Handle special cases
	if x == 0 {
		if y > 0 {
			return math.Pi / 2
		} else if y < 0 {
			return -math.Pi / 2
		}
		return 0 // undefined for (0, 0)
	}

	// Calculate the approximate arctangent
	z := y / x
	atanZ := (math.Pi / 4) - (z*z)/(1+z*z)

	// Adjust the result based on the quadrant
	if x < 0 {
		return atanZ + math.Pi
	}
	return atanZ
}

const FastCosTP = 1.0 / (2.0 * math.Pi)

// FastCos computes the cosine of x using a Taylor series approximation.
func FastCos(x float64) float64 {
	x *= FastCosTP
	x -= 0.25 + math.Floor(x+0.25)
	x *= 16.0 * (math.Abs(x) - 0.5)
	// Comment below for faster rendering
	x += 0.225 * x * (math.Abs(x) - 1.0)
	// end
	return x
}

// FastCos2 is an even faster
// but less accurate version of FastCos
func FastCos2(x float64) float64 {
	x *= FastCosTP
	x -= 0.25 + math.Floor(x+0.25)
	x *= 16.0 * (math.Abs(x) - 0.5)

	return x
}

func FastSin(x float64) float64 {
	return FastCos(x - math.Pi/2)
}

func fastlog2(x float32) float32 {
	vx := *(*uint32)(unsafe.Pointer(&x))
	mx := (vx & 0x007FFFFF) | (0x7e << 23)
	y := float32(vx) * (1.0 / (1 << 23))
	mf := *(*float32)(unsafe.Pointer(&mx))

	return y - 124.22544637 - 1.498030302*mf - 1.72587999/(0.3520887068+mf)
}

func Fastlog(x float32) float32 {
	return 0.69314718 * fastlog2(x)
}

func FastLog64(x float64) float64 {
	return float64(Fastlog(float32(x)))
}

func fastpow2(p float32) float32 {
	vp := (*struct {
		f float32
		i uint32
	})(unsafe.Pointer(&p))
	sign := int32(vp.i >> 31)
	w := int32(p)
	z := p - float32(w) + float32(sign)
	v := struct {
		i uint32
		f float32
	}{
		i: uint32((1 << 23) * (p + 121.2740838 + 27.7280233/(4.84252568-z) - 1.49012907*z)),
	}
	return v.f
}

func FastExp(p float32) float32 {
	return fastpow2(1.442695040 * p)
}

/*
fastpow (float x,
         float p)
{
  return fastpow2 (p * fastlog2 (x));
}
*/

func FastPow(x, p float32) float32 {
	return fastpow2(p * fastlog2(x))
}

// FastPow64 is a wrapper for FastPow that takes and returns float64s
func FastPow64(x, p float64) float64 {
	return float64(FastPow(float32(x), float32(p)))
}
