package renderer

import (
	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type MarchResult struct {
	HitObject drawables.Drawable
	HitPos    vec3.Vec3
	Steps     int
	Distance  float64
	Mhd       float64
}

func RayMarch(ray Ray, renderer *Renderer, showLight bool) MarchResult {
	scene := renderer.scene
	totalDistTraveled := 0.0
	curPos := ray.origin
	totalMin := renderer.scene.options.trace.maxDist
	var closest drawables.Drawable
	steps := 0
	minDistAvg := 0.0
	maxTraceCubed := renderer.scene.options.trace.maxDist * renderer.scene.options.trace.maxDist //* MAXIMUM_TRACE_DISTANCE

	for totalDistTraveled < renderer.scene.options.trace.maxDist {
		minDist := renderer.scene.options.trace.maxDist
		for _, obj := range scene.Drawables {
			dist := obj.Dist(curPos)
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		if renderer.scene.options.shadows && showLight {
			for _, obj := range scene.Lights {
				dist := obj.Dist(curPos)
				if dist < minDist {
					minDist = dist
					closest = obj
				}
			}
		}

		oldAvg := minDistAvg
		minDistAvg -= minDistAvg / 3
		minDistAvg += minDist / 3
		minDistSlope := minDistAvg - oldAvg

		if steps == renderer.scene.options.trace.maxSteps {
			return MarchResult{closest, curPos, renderer.scene.options.trace.maxSteps, totalDistTraveled, renderer.scene.options.trace.minHitDist}
		}

		minHitDist := renderer.scene.options.trace.minHitDist
		if renderer.scene.options.trace.LOD {
			distFromCamera := curPos.Sub(renderer.camera.Pos).Norm()
			minHitDist += (distFromCamera * distFromCamera /* * distFromCamera */ / maxTraceCubed * renderer.scene.options.trace.maxHitDist)
		}
		if minDistSlope < 0 && minDist < minHitDist {

			retPos := curPos
			if minDist < 0 {
				// retPos = curPos.Add(ray.dir.Mult(minDist))
				retPos = retPos.Sub(ray.dir.Mult(minHitDist))
			}

			return MarchResult{closest, retPos, steps, totalDistTraveled, minHitDist}
		}
		distP := minDist * 0.95

		curPos = curPos.Add(ray.dir.Mult(distP))
		steps++

		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist

		}

	}
	return MarchResult{nil, curPos, steps, totalDistTraveled, renderer.scene.options.trace.minHitDist}

}

func SurfaceNormal(p vec3.Vec3, obj drawables.Drawable, epsilon float64) vec3.Vec3 {
	// epsilon := 0.0001 // arbitrary — should be smaller than any surface detail in your distance function, but not so small as to get lost in float precision
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
