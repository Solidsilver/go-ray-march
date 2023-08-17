package renderer

import (
	"image/color"
	"math"

	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

func CalculateReflectionsForHit(marchRslt MarchResult, renderer *Renderer) MarchResult {
	if marchRslt.HitObject != nil && renderer.scene.options.shadows {
		hitPoint := marchRslt.HitPos
		surfaceNormal := SurfaceNormal(hitPoint, marchRslt.HitObject)
		surfaceNormal = surfaceNormal.Unit()

		reflectVec := marchRslt.Dir.Reflect(surfaceNormal)
		bounceDeg := vec3.Angle(reflectVec, surfaceNormal)
		if bounceDeg < 90 {
			ray := Ray{hitPoint, reflectVec}
			refMarchRslt := RayMarch(ray, renderer.scene)
			refMarchRslt.Steps += marchRslt.Steps
			refMarchRslt.Distance += marchRslt.Distance
			if refMarchRslt.Steps > MAX_STEPS {
				refMarchRslt.ReachedMaxSteps = true
			}
			if refMarchRslt.Distance > MAXIMUM_TRACE_DISTANCE {
				refMarchRslt.ReachedMaxDistance = true
			}
			return refMarchRslt

		}
		marchRslt.DidHit = false
	}
	return marchRslt
}

func CalculateReflectionColor(hitColor color.RGBA, objColor color.RGBA, refl float64, renderer *Renderer) color.RGBA {

	if !renderer.scene.options.reflections {
		return objColor
	}

	objColorVec := vec3.RGBAToVec3(objColor)
	reflectivity := math.Min(1, refl)
	refObjColVec := vec3.RGBAToVec3(hitColor)

	objColorVec = objColorVec.Mult(1 - reflectivity).Add(refObjColVec.Mult(reflectivity))
	return vec3.Vec3ToRGBA(objColorVec, 255)

}
