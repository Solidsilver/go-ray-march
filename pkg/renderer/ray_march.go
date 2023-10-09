package renderer

import (
	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type MarchResult struct {
	HitObject drawables.Drawable
	HitPos    vec3.Vec3
	Steps     int
	Distance  float64
}

func minDistSlope(rbf *utils.RingBuffer[float64]) float64 {
	sum := 0.0
	for i := 0; i > 1-rbf.Size; i-- {
		slope := rbf.Get(i) - rbf.Get(i-1)
		sum += slope
	}

	return sum / float64(rbf.Size-1)
}

func RayMarch(ray Ray, scene *Scene) MarchResult {
	totalDistTraveled := 0.0
	curPos := ray.origin
	totalMin := MAXIMUM_TRACE_DISTANCE
	var closest drawables.Drawable
	steps := 0
	rbf := utils.NewRingBufferFilled[float64](3, -1)

	for totalDistTraveled < MAXIMUM_TRACE_DISTANCE {
		minDist := MAXIMUM_TRACE_DISTANCE
		for _, obj := range scene.Drawables {
			dist := obj.Dist(curPos)
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		// if !ignoreLights {
		for _, obj := range scene.Lights {
			dist := obj.Dist(curPos)
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		// }

		rbf.Push(minDist)
		mds := minDistSlope(rbf)

		if steps == MAX_STEPS {
			return MarchResult{closest, curPos, MAX_STEPS, totalDistTraveled}
		}

		if mds < 0 && minDist < MINIMUM_HIT_DISTANCE {

			// println(minDistSlope(rbf))
			retPos := curPos
			if minDist < 0 {
				retPos = curPos.Add(ray.dir.Mult(minDist))
				retPos = retPos.Sub(ray.dir.Mult(MINIMUM_HIT_DISTANCE))
			}

			return MarchResult{closest, retPos, steps, totalDistTraveled}
		}
		distP := minDist * 0.95

		curPos = curPos.Add(ray.dir.Mult(distP))
		steps++

		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist

		}

	}
	return MarchResult{nil, curPos, steps, totalDistTraveled}

}

func SurfaceNormal(p vec3.Vec3, obj drawables.Drawable) vec3.Vec3 {
	epsilon := MINIMUM_HIT_DISTANCE // arbitrary — should be smaller than any surface detail in your distance function, but not so small as to get lost in float precision
	centerDistance := obj.Dist(p)
	grad := vec3.Vec3{
		X: obj.Dist(p.Add(vec3.Vec3{X: epsilon, Y: 0, Z: 0})),
		Y: obj.Dist(p.Add(vec3.Vec3{X: 0, Y: epsilon, Z: 0})),
		Z: obj.Dist(p.Add(vec3.Vec3{X: 0, Y: 0, Z: epsilon})),
	}
	normal := grad.Minus(centerDistance).Div(epsilon)

	return normal
}
