package renderer

import (
	"math"

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
			dist := 0.0
			if renderer.scene.options.trace.fastMath {
				dist = obj.FastDist(curPos)
			} else {
				dist = obj.Dist(curPos)
			}
			// dist := obj.Dist(curPos)
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		if renderer.scene.options.shadows && showLight {
			for _, obj := range scene.Lights {
				dist := 0.0
				if renderer.scene.options.trace.fastMath {
					dist = obj.FastDist(curPos)
				} else {
					dist = obj.Dist(curPos)
				}
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

func RayMarchP(ray Ray, renderer *Renderer, showLight bool) MarchResult {
	// scene := renderer.scene
	// inside := false
	totalDistTraveled := 0.0
	curPos := vec3.NewCp(ray.origin)
	totalMin := renderer.scene.options.trace.maxDist
	var closest drawables.Drawable
	steps := 0
	minDistAvg := 0.0
	maxTraceCubed := renderer.scene.options.trace.maxDist * renderer.scene.options.trace.maxDist //* MAXIMUM_TRACE_DISTANCE

	for totalDistTraveled < renderer.scene.options.trace.maxDist {
		minDist := renderer.scene.options.trace.maxDist
		for _, obj := range renderer.scene.Drawables {
			dist := 0.0
			if renderer.scene.options.trace.fastMath {
				dist = obj.FastDist(*curPos)
			} else {
				dist = obj.Dist(*curPos)
			}
			// dist := obj.Dist(curPos)
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		if renderer.scene.options.shadows && showLight {
			for _, obj := range renderer.scene.Lights {
				dist := 0.0
				if renderer.scene.options.trace.fastMath {
					dist = obj.FastDist(*curPos)
				} else {
					dist = obj.Dist(*curPos)
				}
				if dist < minDist {
					minDist = dist
					closest = obj
				}
			}
		}

		// if steps == 0 && minDist < 0 {
		// 	inside = true
		// }
		// if minDist < 0 &&
		//  {
		// 	minDist = -minDist
		// }

		oldAvg := minDistAvg
		minDistAvg -= minDistAvg / 3
		minDistAvg += minDist / 3
		minDistSlope := minDistAvg - oldAvg

		if steps == renderer.scene.options.trace.maxSteps {
			return MarchResult{closest, *curPos, renderer.scene.options.trace.maxSteps, totalDistTraveled, renderer.scene.options.trace.minHitDist}
		}

		minHitDist := renderer.scene.options.trace.minHitDist
		if renderer.scene.options.trace.LOD {
			distFromCamera := curPos.Sub(renderer.camera.Pos).Norm()
			minHitDist += (distFromCamera * distFromCamera /* * distFromCamera */ / maxTraceCubed * renderer.scene.options.trace.maxHitDist)
		}
		if minDistSlope < 0 && minDist < minHitDist {

			// retPos := curPos
			if minDist < 0 {
				moveBackMinDist := ray.dir.Mult(minHitDist)
				curPos.SubSet(&moveBackMinDist)
				// retPos = curPos.Add(ray.dir.Mult(minDist))
				// retPos = retPos.Sub(ray.dir.Mult(minHitDist))
			}

			return MarchResult{closest, *curPos, steps, totalDistTraveled, minHitDist}
		}
		distP := minDist * 0.95

		moveForwardMinDist := ray.dir.Mult(distP)
		curPos.AddSet(&moveForwardMinDist)
		steps++

		totalDistTraveled += distP
		// if minDist < totalMin {
		// 	totalMin = minDist

		// }
		totalMin = math.Min(totalMin, minDist)

	}
	return MarchResult{nil, *curPos, steps, totalDistTraveled, renderer.scene.options.trace.minHitDist}

}

func SurfaceNormal(hitRslt MarchResult, fast bool) vec3.Vec3 {
	obj := hitRslt.HitObject
	dx := hitRslt.HitPos.Add(vec3.NewX(hitRslt.Mhd))
	dy := hitRslt.HitPos.Add(vec3.NewY(hitRslt.Mhd))
	dz := hitRslt.HitPos.Add(vec3.NewZ(hitRslt.Mhd))
	var normal *vec3.Vec3
	switch fast {
	case true:
		normal = vec3.NewP(
			obj.FastDist(dx),
			obj.FastDist(dy),
			obj.FastDist(dz),
		)
	case false:
		normal = vec3.NewP(
			obj.Dist(dx),
			obj.Dist(dy),
			obj.Dist(dz),
		)
	}
	normal.MinusSet(hitRslt.Distance)
	normal.ToUnitSet()
	return *normal
}
