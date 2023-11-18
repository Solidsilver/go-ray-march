package utils

import "math"

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

func FastCos(x float64) float64 {
	x *= FastCosTP
	x -= 0.25 + math.Floor(x+0.25)
	x *= 16.0 * (math.Abs(x) - 0.5)
	// Comment below for faster rendering
	x += 0.225 * x * (math.Abs(x) - 1.0)
	// end
	return x
}

func FastSin(x float64) float64 {
	return FastCos(x - math.Pi/2)
}
