package renderer

import (
	"image/color"
	"math"

	"solidsilver.dev/go-ray-marching/pkg/drawables"
	"solidsilver.dev/go-ray-marching/pkg/utils"
)

type MarchResult struct {
	HitObject drawables.Drawable
	HitPos    utils.Vec3
	Steps     int
	Distance  float64
}

func RayMarchColor(ray *Ray, scene *Scene, rnd float64) color.RGBA {
	totalDistTraveled := 0.0
	// curPos := utils.NewCopy(ray.origin)
	curPos := ray.origin
	totalMin := MAXIMUM_TRACE_DISTANCE
	var closest drawables.Drawable
	steps := 0

	// if curPos.Equals(utils.Vec3{X: 0, Y: -5, Z: -29}) {
	// 	print("hello")
	// }

	for totalDistTraveled < MAXIMUM_TRACE_DISTANCE {
		minDist := MAXIMUM_TRACE_DISTANCE
		for _, obj := range scene.Drawables {
			dist := obj.Dist(curPos)
			// if dist < MINIMUM_HIT_DISTANCE*2 {
			// 	print("hello")
			// }
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		distP := minDist * (1 - MINIMUM_HIT_DISTANCE)
		curPos.Add(curPos, *utils.NewCopy(ray.dir).Mult(distP))
		steps++

		// if minDist < 0 {
		// 	return color.RGBA{0, 0, 0, 255}
		// }

		if minDist < 0 || minDist < MINIMUM_HIT_DISTANCE {
			// _, distF := math.Modf(totalDistTraveled / 10)
			col := closest.Color()
			// noise := (rand.Float64() - 0.5) * 3
			// rndRslt := utils.LGCRandDec(uint32(totalDistTraveled), 400)
			// stepsFl := float64(steps) + ((rndRslt.Rnd - 0.5) * 0.5) //((rnd - 0.5) * distF * 2)
			distScale := (totalDistTraveled + 20)
			distFrac := totalDistTraveled / distScale
			stepsFl := utils.SigLocal(float64(steps)+(distFrac), steps, 0.5)
			// stepsFl := float64(steps) / totalDistTraveled * 10
			// fmt.Printf("Steps: %d, dist: %f, StepsFL: %f\n", steps, totalDistTraveled, stepsFl)
			// println(stepsFl)
			harsh := 3.0
			darkP := (1 / math.Sqrt(stepsFl+math.Pow(2, harsh))) * harsh
			r := uint8(darkP * float64(col.R))
			g := uint8(darkP * float64(col.G))
			b := uint8(darkP * float64(col.B))
			// newA := darkP * float64(a)

			return color.RGBA{r, g, b, col.A}
			// return color.RGBA{100, 200, 200, 255}
		}

		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist

		}

	}
	// distPercet := MINIMUM_HIT_DISTANCE / totalMin
	// clr := 255 * distPercet
	// fmt.Println(clr)
	// return color.RGBA{100, 200, 200, uint8(clr)}
	return color.RGBA{0, 0, 0, 0}

}

func RayMarch(ray *Ray, scene *Scene) MarchResult {
	totalDistTraveled := 0.0
	curPos := ray.origin
	totalMin := MAXIMUM_TRACE_DISTANCE
	var closest drawables.Drawable
	steps := 0

	for totalDistTraveled < MAXIMUM_TRACE_DISTANCE {
		minDist := MAXIMUM_TRACE_DISTANCE
		for _, obj := range scene.Drawables {
			dist := obj.Dist(curPos)
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		for _, obj := range scene.Lights {
			dist := obj.Dist(curPos)
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}

		if minDist < 0 || minDist < MINIMUM_HIT_DISTANCE {
			retPos := utils.NewCopy(curPos)
			retPos.Sub(*retPos, *utils.NewCopy(ray.dir).Mult(MINIMUM_HIT_DISTANCE))
			return MarchResult{closest, *retPos, steps, totalDistTraveled}
		}
		distP := minDist * (1 - MINIMUM_HIT_DISTANCE)

		curPos.Add(curPos, *utils.NewCopy(ray.dir).Mult(distP))
		steps++

		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist

		}

	}
	return MarchResult{nil, curPos, steps, totalDistTraveled}

}

func SurfaceNormal(p utils.Vec3, obj drawables.Drawable) utils.Vec3 {
	epsilon := 0.0001 // arbitrary — should be smaller than any surface detail in your distance function, but not so small as to get lost in float precision
	centerDistance := obj.Dist(p)
	xDistance := obj.Dist(*utils.NewAdd(p, utils.Vec3{X: epsilon, Y: 0, Z: 0}))
	yDistance := obj.Dist(*utils.NewAdd(p, utils.Vec3{X: 0, Y: epsilon, Z: 0}))
	zDistance := obj.Dist(*utils.NewAdd(p, utils.Vec3{X: 0, Y: 0, Z: epsilon}))
	normal := utils.NewVec(xDistance, yDistance, zDistance)
	normal.Minus(centerDistance)
	normal.Div(epsilon)
	// return normal.Div(epsilon)
	return *normal
}
