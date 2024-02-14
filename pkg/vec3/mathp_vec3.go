package vec3

func (v1 *Vec3) AddSet(v2 *Vec3) {
	v1.X = v1.X + v2.X
	v1.Y = v1.Y + v2.Y
	v1.Z = v1.Z + v2.Z
}

func (v1 *Vec3) SubSet(v2 *Vec3) {
	v1.X = v1.X - v2.X
	v1.Y = v1.Y - v2.Y
	v1.Z = v1.Z - v2.Z
}

func (v *Vec3) PlusSet(num float64) {
	v.X += num
	v.Y += num
	v.Z += num
}

func (v *Vec3) MinusSet(num float64) {
	v.X -= num
	v.Y -= num
	v.Z -= num
}

func (v *Vec3) MultSet(num float64) {
	v.X *= num
	v.Y *= num
	v.Z *= num
}

func (v1 *Vec3) MultCompSet(v2 *Vec3) {
	v1.X *= v2.X
	v1.Y *= v2.Y
	v1.Z *= v2.Z
}

func (v *Vec3) DivSet(num float64) {
	v.X = v.X / num
	v.Y = v.Y / num
	v.Z = v.Z / num
}

func (v1 *Vec3) CrossSet(v2 *Vec3) {
	v1.X = v1.Y*v2.Z - v1.Z*v2.Y
	v1.Y = v1.Z*v2.X - v1.X*v2.Z
	v1.Z = v1.X*v2.Y - v1.Y*v2.X
}

func (v *Vec3) ToUnitSet() {
	v.DivSet(v.Mag())
}
