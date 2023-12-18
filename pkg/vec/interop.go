package vec

import (
	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

func (v *Vec[T]) ToVec2() utils.Vec2 {
	return utils.Vec2{float64(v.GetOrZero(0)), float64(v.GetOrZero(1))}
}

func (v *Vec[T]) ToVec3() vec3.Vec3 {
	return vec3.Vec3{X: float64(v.GetOrZero(0)), Y: float64(v.GetOrZero(1)), Z: float64(v.GetOrZero(2))}
}
