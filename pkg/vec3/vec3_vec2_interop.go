package vec3

import "github.com/Solidsilver/go-ray-march/pkg/utils"

func (v *Vec3) XY() utils.Vec2 {
	return utils.NewVec2(v.X, v.Y)
}

func (v *Vec3) XZ() utils.Vec2 {
	return utils.NewVec2(v.X, v.Z)
}

func (v *Vec3) YZ() utils.Vec2 {
	return utils.NewVec2(v.Y, v.Z)
}
