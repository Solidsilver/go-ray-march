package utils

import "math"

type Vec3 struct {
	X, Y, Z float64
}

func NewVec3(x, y, z float64) *Vec3 {
	vec := new(Vec3)
	vec.X = x
	vec.Y = y
	vec.Z = z
	return vec
}

func Vec3Size(n float64) Vec3 {
	return Vec3{n, n, n}
}

func Vec3Zero() Vec3 {
	return Vec3{0, 0, 0}
}

func Vec3UnitX() Vec3 {
	return Vec3{1, 0, 0}
}

func Vec3UnitY() Vec3 {
	return Vec3{0, 1, 0}
}

func Vec3UnitZ() Vec3 {
	return Vec3{0, 0, 1}
}

func (v *Vec3) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func NewNorm(v Vec3) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func Vec3Max(v1, v2 Vec3) Vec3 {
	return Vec3{
		math.Max(v1.X, v2.X),
		math.Max(v1.Y, v2.Y),
		math.Max(v1.Z, v2.Z),
	}
}

// func Vec3Len(v Vec3) float64 {
// 	return v.Mag()
// }

/*
Adds v1 to v2 and saves the result to the current vector.
The result is also returned
*/
func (v *Vec3) Add(v1, v2 Vec3) *Vec3 {
	v.X = v1.X + v2.X
	v.Y = v1.Y + v2.Y
	v.Z = v1.Z + v2.Z
	return v
}

func (v Vec3) Equals(v1 Vec3) bool {
	return v.X == v1.X && v.Y == v1.Y && v.Z == v1.Z
}

/*
Adds v1 to v2 and returns the result as a newly allocated vector.
*/
func NewAdd(v1, v2 Vec3) *Vec3 {
	v := new(Vec3)
	return v.Add(v1, v2)
}

func NewAdd2(v1, v2 Vec3) Vec3 {
	return Vec3{
		v1.X + v2.X,
		v1.Y + v2.Y,
		v1.Z + v2.Z,
	}
}

func (v *Vec3) Plus(num float64) *Vec3 {
	v.X += num
	v.Y += num
	v.Z += num
	return v
}

func NewPlus(v1 Vec3, num float64) Vec3 {
	return Vec3{
		v1.X + num,
		v1.Y + num,
		v1.Z + num,
	}
}

func (v *Vec3) Minus(num float64) *Vec3 {
	v.X -= num
	v.Y -= num
	v.Z -= num
	return v
}

func NewMinus(v1 Vec3, num float64) Vec3 {
	return NewPlus(v1, num*-1)
}

/*
Subtracts v1 from v2 and saves the result to the current vector.
The result is also returned
*/
func (v *Vec3) Sub(v1, v2 Vec3) *Vec3 {
	v.X = v1.X - v2.X
	v.Y = v1.Y - v2.Y
	v.Z = v1.Z - v2.Z
	return v
}

/*
Adds v1 to v2 and returns the result as a newly allocated vector.
*/
func NewSub(v1, v2 Vec3) *Vec3 {
	v := new(Vec3)
	return v.Sub(v1, v2)
}

func (v *Vec3) Mult(num float64) *Vec3 {
	v.X = v.X * num
	v.Y = v.Y * num
	v.Z = v.Z * num
	return v
}

func (v *Vec3) Div(num float64) *Vec3 {
	v.X = v.X / num
	v.Y = v.Y / num
	v.Z = v.Z / num
	return v
}

func Dot(v1, v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v *Vec3) Cross(v1, v2 Vec3) *Vec3 {
	v.X = v1.Y*v2.Z - v1.Z*v2.Y
	v.Y = v1.Z*v2.X - v1.X*v2.Z
	v.Z = v1.X*v2.Y - v1.Y*v2.X

	return v
}

func (v *Vec3) Copy(v1 Vec3) *Vec3 {
	v.X = v1.X
	v.Y = v1.Y
	v.Z = v1.Z
	return v
}

func NewCopy(v1 Vec3) *Vec3 {
	v := new(Vec3)
	v.X = v1.X
	v.Y = v1.Y
	v.Z = v1.Z
	return v
}

func (v *Vec3) Unit() *Vec3 {
	vNorm := v.Norm()
	return v.Div(vNorm)
}

func (v Vec3) NewUnit() Vec3 {
	return *v.Unit()
}

func Angle(v1 Vec3, v2 Vec3) float64 {
	return RadToDeg(math.Acos(Dot(v1, v2) / (v1.Norm() * v2.Norm())))
}

func DirFromPos(p1 Vec3, p2 Vec3) Vec3 {
	dir := new(Vec3)
	dir.Sub(p1, p2)
	dir.Unit()
	return *dir
}

func Abs(p Vec3) Vec3 {
	return Vec3{
		math.Abs(p.X),
		math.Abs(p.Y),
		math.Abs(p.Z),
	}
}
