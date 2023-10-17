package renderer

import (
	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"goki.dev/mat32/v2"
)

type MarchResult struct {
	HitObject drawables.Drawable
	HitPos    mat32.Vec3
	Steps     int
	Distance  float32
	Mhd       float32
}

func RayMarch(ray Ray, renderer *Renderer) MarchResult {
	scene := renderer.scene
	totalDistTraveled := float32(0.0)
	curPos := ray.origin
	totalMin := float32(MAXIMUM_TRACE_DISTANCE)
	var closest drawables.Drawable
	steps := 0
	minDistAvg := float32(0.0)
	maxTraceCubed := float32(MAXIMUM_TRACE_DISTANCE * MAXIMUM_TRACE_DISTANCE) //* MAXIMUM_TRACE_DISTANCE

	for totalDistTraveled < MAXIMUM_TRACE_DISTANCE {
		minDist := float32(MAXIMUM_TRACE_DISTANCE)
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

		oldAvg := minDistAvg
		minDistAvg -= minDistAvg / 3.0
		minDistAvg += minDist / 3.0
		minDistSlope := minDistAvg - oldAvg

		if steps == MAX_STEPS {
			return MarchResult{closest, curPos, MAX_STEPS, totalDistTraveled, MINIMUM_HIT_DISTANCE}
		}

		minHitDist := float32(MINIMUM_HIT_DISTANCE)
		if LOD {
			distFromCamera := curPos.DistTo(renderer.camera.Pos)
			minHitDist += (distFromCamera * distFromCamera /* * distFromCamera */ / maxTraceCubed * float32(MAX_HIT_DISTANCE))
		}
		if minDistSlope < 0 && minDist < minHitDist {

			retPos := curPos
			if minDist < 0 {
				retPos = retPos.Sub(ray.dir.MulScalar(minHitDist))
			}

			return MarchResult{closest, retPos, steps, totalDistTraveled, minHitDist}
		}
		distP := minDist * 0.95

		curPos = curPos.Add(ray.dir.MulScalar(distP))
		steps++

		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist

		}

	}
	return MarchResult{nil, curPos, steps, totalDistTraveled, MINIMUM_HIT_DISTANCE}

}

func SurfaceNormal(p mat32.Vec3, obj drawables.Drawable, epsilon float32) mat32.Vec3 {
	// epsilon := 0.0001 // arbitrary — should be smaller than any surface detail in your distance function, but not so small as to get lost in float precision
	centerDistance := obj.Dist(p)
	grad := mat32.Vec3{
		X: obj.Dist(p.Add(mat32.Vec3{X: epsilon, Y: 0, Z: 0})),
		Y: obj.Dist(p.Add(mat32.Vec3{X: 0, Y: epsilon, Z: 0})),
		Z: obj.Dist(p.Add(mat32.Vec3{X: 0, Y: 0, Z: epsilon})),
	}
	normal := grad.SubScalar(centerDistance).DivScalar(epsilon)

	return normal
}
