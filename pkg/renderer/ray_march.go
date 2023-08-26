package renderer

import (
	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
	"github.com/Solidsilver/go-ray-march/pkg/vec3neon"
)

type MarchResult struct {
	HitObject          drawables.Drawable
	HitPos             vec3.Vec3
	Steps              int
	Distance           float64
	DidHit             bool
	Dir                vec3.Vec3
	ReachedMaxSteps    bool
	ReachedMaxDistance bool
}

type MarchResultN struct {
	HitObject          drawables.Drawable
	HitPos             vec3neon.Vec3Neon
	Steps              int
	Distance           float64
	DidHit             bool
	Dir                vec3neon.Vec3Neon
	ReachedMaxSteps    bool
	ReachedMaxDistance bool
}

func minDistSlope(rbf *utils.RingBuffer[float64]) float64 {
	sum := 0.0
	for i := 0; i > -(rbf.Size() - 1); i-- {
		slope := rbf.Get(i) - rbf.Get(i-1)
		sum += slope
	}

	return sum / float64(rbf.Size()-1)
}

func minDistSlopeF32(rbf *utils.RingBuffer[float32]) float32 {
	sum := float32(0.0)
	for i := 0; i > -(rbf.Size() - 1); i-- {
		slope := rbf.Get(i) - rbf.Get(i-1)
		sum += slope
	}

	return sum / float32(rbf.Size()-1)
}

func RayMarch(ray Ray, scene *Scene) MarchResult {
	totalDistTraveled := 0.0
	curPos := ray.origin
	totalMin := MAXIMUM_TRACE_DISTANCE
	var closest drawables.Drawable
	steps := 0
	rbf := utils.NewRingBuffer[float64](3)
	for i := 0; i < rbf.Size(); i++ {
		rbf.Push(-1)
	}

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
			return MarchResult{closest, curPos, MAX_STEPS, totalDistTraveled, false, ray.dir, true, false}
		}

		if mds < 0 && minDist < MINIMUM_HIT_DISTANCE {

			// println(minDistSlope(rbf))
			retPos := curPos
			if minDist < 0 {
				retPos = curPos.Add(ray.dir.Mult(minDist))
				retPos = retPos.Sub(ray.dir.Mult(MINIMUM_HIT_DISTANCE))
			}

			return MarchResult{closest, retPos, steps, totalDistTraveled, true, ray.dir, false, false}
		}
		distP := minDist * 0.95

		curPos = curPos.Add(ray.dir.Mult(distP))
		steps++

		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist

		}

	}
	return MarchResult{nil, curPos, steps, totalDistTraveled, false, ray.dir, false, true}

}

func RayMarchNeon(rayIn Ray, scene *Scene) MarchResult {
	ray := rayIn.ToRayN()
	totalDistTraveled := float32(0.0)
	curPos := ray.origin
	totalMin := float32(MAXIMUM_TRACE_DISTANCE)
	var closest drawables.Drawable
	steps := 0
	rbf := utils.NewRingBuffer[float32](3)
	for i := 0; i < rbf.Size(); i++ {
		rbf.Push(-1)
	}

	for totalDistTraveled < MAXIMUM_TRACE_DISTANCE {
		minDist := float32(MAXIMUM_TRACE_DISTANCE)
		for _, obj := range scene.Drawables {
			dist := (obj.DistN(curPos))
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		// if !ignoreLights {
		for _, obj := range scene.Lights {
			dist := obj.DistN(curPos)
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		// }

		rbf.Push(minDist)
		mds := minDistSlopeF32(rbf)

		if steps == MAX_STEPS {
			return MarchResult{closest, curPos.ToVec3(), MAX_STEPS, float64(totalDistTraveled), false, ray.dir.ToVec3(), true, false}
		}

		if mds < 0 && minDist < MINIMUM_HIT_DISTANCE {

			// println(minDistSlope(rbf))
			retPos := curPos
			if minDist < 0 {
				retPos = curPos.Add(ray.dir.Mult(float32(minDist)))
				retPos = retPos.Sub(ray.dir.Mult(MINIMUM_HIT_DISTANCE))
			}

			return MarchResult{closest, retPos.ToVec3(), steps, float64(totalDistTraveled), true, ray.dir.ToVec3(), false, false}
		}
		distP := minDist * 0.95

		curPos = curPos.Add(ray.dir.Mult(float32(distP)))
		steps++

		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist

		}

	}
	return MarchResult{nil, curPos.ToVec3(), steps, float64(totalDistTraveled), false, ray.dir.ToVec3(), false, true}

}

func RayMarchNeon2(ray RayN, scene *Scene) MarchResultN {
	totalDistTraveled := float32(0.0)
	curPos := ray.origin
	totalMin := float32(MAXIMUM_TRACE_DISTANCE)
	var closest drawables.Drawable
	steps := 0
	rbf := utils.NewRingBuffer[float32](3)
	for i := 0; i < rbf.Size(); i++ {
		rbf.Push(-1)
	}

	for totalDistTraveled < MAXIMUM_TRACE_DISTANCE {
		minDist := float32(MAXIMUM_TRACE_DISTANCE)
		for _, obj := range scene.Drawables {
			dist := (obj.DistN(curPos))
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		// if !ignoreLights {
		for _, obj := range scene.Lights {
			dist := obj.DistN(curPos)
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		// }

		rbf.Push(minDist)
		mds := minDistSlopeF32(rbf)

		if steps == MAX_STEPS {
			return MarchResultN{closest, curPos, MAX_STEPS, float64(totalDistTraveled), false, ray.dir, true, false}
		}

		if mds < 0 && minDist < MINIMUM_HIT_DISTANCE {

			// println(minDistSlope(rbf))
			retPos := curPos
			if minDist < 0 {
				retPos = curPos.Add(ray.dir.Mult(float32(minDist)))
				retPos = retPos.Sub(ray.dir.Mult(MINIMUM_HIT_DISTANCE))
			}

			return MarchResultN{closest, retPos, steps, float64(totalDistTraveled), true, ray.dir, false, false}
		}
		distP := minDist * 0.95

		curPos = curPos.Add(ray.dir.Mult(float32(distP)))
		steps++

		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist

		}

	}
	return MarchResultN{nil, curPos, steps, float64(totalDistTraveled), false, ray.dir, false, true}

}

func SurfaceNormal(p vec3.Vec3, obj drawables.Drawable) vec3.Vec3 {
	epsilon := 0.0001 // arbitrary — should be smaller than any surface detail in your distance function, but not so small as to get lost in float precision
	centerDistance := obj.Dist(p)
	grad := vec3.Vec3{
		X: obj.Dist(p.Add(vec3.Vec3{X: epsilon, Y: 0, Z: 0})),
		Y: obj.Dist(p.Add(vec3.Vec3{X: 0, Y: epsilon, Z: 0})),
		Z: obj.Dist(p.Add(vec3.Vec3{X: 0, Y: 0, Z: epsilon})),
	}
	normal := grad.Minus(centerDistance)
	normal = normal.Div(epsilon)

	return normal
}

func SurfaceNormalNeon(p vec3neon.Vec3Neon, obj drawables.Drawable) vec3neon.Vec3Neon {
	epsilon := float32(0.0001) // arbitrary — should be smaller than any surface detail in your distance function, but not so small as to get lost in float precision
	centerDistance := obj.DistN(p)
	grad := vec3neon.New(
		obj.DistN(p.Add(*vec3neon.New(epsilon, 0, 0))),
		obj.DistN(p.Add(*vec3neon.New(0, epsilon, 0))),
		obj.DistN(p.Add(*vec3neon.New(0, 0, epsilon))),
	)
	normal := grad.Minus(centerDistance)
	normal = normal.Div(epsilon)

	return normal
}
