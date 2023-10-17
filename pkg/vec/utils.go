package vec

import (
	"image/color"

	"goki.dev/mat32/v2"
)

func MaxVec3(a, b mat32.Vec3) mat32.Vec3 {
	return mat32.Vec3{X: mat32.Max(a.X, b.X), Y: mat32.Max(a.Y, b.Y), Z: mat32.Max(a.Z, b.Z)}
}

func MinVec3(a, b mat32.Vec3) mat32.Vec3 {
	return mat32.Vec3{X: mat32.Min(a.X, b.X), Y: mat32.Min(a.Y, b.Y), Z: mat32.Min(a.Z, b.Z)}
}

func DirFromPos(pt1, p2 mat32.Vec3) mat32.Vec3 {
	return p2.Sub(pt1).Normal()
}

func RGBAToVec3(color color.RGBA) mat32.Vec3 {
	return mat32.Vec3{X: float32(color.R) / 255.0, Y: float32(color.G) / 255.0, Z: float32(color.B) / 255.0}
}

func Vec3ToRGBA(colorVec mat32.Vec3, a uint8) color.RGBA {
	colorVec = colorVec.MulScalar(255)
	return color.RGBA{R: uint8(colorVec.X * 255), G: uint8(colorVec.Y * 255), B: uint8(colorVec.Z * 255), A: a}
}
