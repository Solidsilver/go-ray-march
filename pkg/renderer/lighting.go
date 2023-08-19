package renderer

import (
	"image/color"
	"math"

	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

func CalculateLighting2(marchRslt MarchResult, incomingColor color.RGBA, renderer *Renderer) color.RGBA {
	pxColorVec := vec3.Zero()
	if marchRslt.HitObject != nil {
		pxColorVec = vec3.RGBAToVec3(incomingColor)
		if renderer.scene.options.shadows {
			hitPoint := marchRslt.HitPos
			colorVec := vec3.Zero()
			for _, lSource := range renderer.scene.Lights {
				lightDir := vec3.DirFromPos(lSource.Pos(), hitPoint).Unit()
				surfaceNormal := SurfaceNormal(hitPoint, marchRslt.HitObject)
				bounceDeg := vec3.Angle(lightDir, surfaceNormal)
				if bounceDeg < 90 {
					ray := Ray{hitPoint, lightDir}
					rslt := RayMarch(ray, renderer.scene)
					if drawables.Equals(rslt.HitObject, lSource) {
						brightness := float64(rslt.HitObject.Color().A) / 255
						brightness = brightness * (90 - bounceDeg) / 90
						lightColorVec := vec3.RGBAToVec3(lSource.Color()).Mult(brightness)
						colorVec = colorVec.Add(lightColorVec)
					}
				}
			}
			pxColorVec = vec3.Min(pxColorVec.MultComp(colorVec), vec3.OfSize(1))
		}

	}

	return vec3.Vec3ToRGBA(pxColorVec, 255)
}

func CalculatePostProcessing(colorVecIn vec3.Vec3, marchRslt MarchResult, screenPos Point, renderer *Renderer) color.RGBA {
	pxColorVec := colorVecIn
	if renderer.scene.options.ao.enabled && marchRslt.HitObject != nil {
		ao := 1.0 - float64(marchRslt.Steps)/float64(MAX_STEPS-1)
		pxColorVec = pxColorVec.Mult(ao)
	}

	if renderer.scene.options.dropoff.enabled {
		dropoffDist := math.Min(renderer.scene.options.dropoff.distance, MAXIMUM_TRACE_DISTANCE)
		distFrac := math.Min((marchRslt.Distance)/float64(dropoffDist), 1)
		dropoff := 1 - math.Pow(distFrac, 2)
		blendColor := vec3.RGBAToVec3(renderer.scene.options.dropoff.color)
		pxColorVec = pxColorVec.Mult(dropoff).Add(blendColor.Mult(1 - dropoff))
	}

	if renderer.scene.options.vignette.enabled {
		maxVignettNorm := utils.NewVec2(float64(renderer.camera.SizeX), float64(renderer.camera.SizeY)).Norm() * math.Min(1, (1-math.Min(1, renderer.scene.options.vignette.strength)))
		vignettAmt := 1 - (utils.NewVec2(float64(screenPos.X-renderer.camera.centerOffset.X), float64(screenPos.Y-renderer.camera.centerOffset.Y)).Norm() / maxVignettNorm)
		pxColorVec = pxColorVec.Mult(vignettAmt)
	}

	pxColorVal := vec3.Vec3ToRGBA(pxColorVec, 255)
	return pxColorVal
}

func CalculateLightingForHit(marchRslt MarchResult, renderer *Renderer) LightingResult {
	colorVec := vec3.Unit()
	if marchRslt.HitObject != nil && renderer.scene.options.shadows {
		hitPoint := marchRslt.HitPos
		colorVec = vec3.Zero()
		for _, lSource := range renderer.scene.Lights {
			lightDir := vec3.DirFromPos(lSource.Pos(), hitPoint).Unit()
			surfaceNormal := SurfaceNormal(hitPoint, marchRslt.HitObject)
			bounceDeg := vec3.Angle(lightDir, surfaceNormal)
			if bounceDeg < 90 {
				ray := Ray{hitPoint, lightDir}
				rslt := RayMarch(ray, renderer.scene)
				if drawables.Equals(rslt.HitObject, lSource) {
					brightness := float64(rslt.HitObject.Color().A) / 255
					brightness = brightness * (90 - bounceDeg) / 90
					lightColorVec := vec3.RGBAToVec3(lSource.Color()).Mult(brightness)
					colorVec = colorVec.Add(lightColorVec)
				}
			}
		}
	}
	return LightingResult{
		ColorVec: colorVec,
	}
}

type LightingResult struct {
	ColorVec vec3.Vec3
}

func calculatePhongReflectanceVec(ambientI, hitPoint vec3.Vec3, obj drawables.Drawable, rnd *Renderer, recursion int64) vec3.Vec3 {
	if recursion > 100 {
		// log.Info().Msg("Max recursion depth reached")
		// return vec3.Zero()
		// return vec3.Unit()
		return vec3.RGBAToVec3(obj.Color())
	}
	refProps := obj.ReflectionProperties()
	objColor := vec3.RGBAToVec3(obj.Color())

	specular := objColor.Mult(refProps.Metalness).Add(vec3.OfSize(1 - refProps.Metalness))
	surfaceNormal := SurfaceNormal(hitPoint, obj).Unit()

	outColorVec := objColor.MultComp(ambientI).Mult(refProps.Ambient)
	for _, lSource := range rnd.scene.Lights {
		lightColor := vec3.RGBAToVec3(lSource.Color())
		lightDir := vec3.DirFromPos(obj.Pos(), lSource.Pos())
		angle := vec3.Angle(lightDir, surfaceNormal)
		if angle < 90 {

			ray := Ray{hitPoint, lightDir}
			rslt := RayMarch(ray, rnd.scene)
			if rslt.HitObject != nil && drawables.Equals(rslt.HitObject, lSource) {
				// Color component from incoming light
				intensityLambert := objColor.MultComp(lightColor)
				ldsNormal := vec3.Dot(lightDir, surfaceNormal)
				ldsNormalMax := math.Max(ldsNormal, 0)
				intensityLambert = intensityLambert.Mult(ldsNormalMax)
				intensityLambert = intensityLambert.Mult(refProps.Lambertian)
				outColorVec = outColorVec.Add(intensityLambert)

				// Color component from specular light
				viewingRay := vec3.DirFromPos(rnd.camera.Pos, hitPoint)
				reflVec := lightDir.Reverse().Add(surfaceNormal.Mult(2).Mult(vec3.Dot(lightDir, surfaceNormal))).Unit()
				intensitySpecular := specular.MultComp(lightColor).Mult(refProps.Specular)
				rdv := vec3.Dot(reflVec.Reverse(), viewingRay)
				rdvMax := math.Max(rdv, 0)
				powSmooth := math.Pow(rdvMax, refProps.Smoothness)
				intensitySpecular = intensitySpecular.Mult(powSmooth)
				// intensitySpecular = intensitySpecular
				outColorVec = outColorVec.Add(intensitySpecular)
			}
		}

	}

	if rnd.scene.options.reflections {
		for _, objLSource := range rnd.scene.Drawables {
			lightDir := vec3.DirFromPos(obj.Pos(), objLSource.Pos())
			angle := vec3.Angle(lightDir, surfaceNormal)
			if !drawables.Equals(objLSource, obj) && angle < 90 {
				ray := Ray{hitPoint, lightDir}
				rslt := RayMarch(ray, rnd.scene)
				if rslt.HitObject != nil && drawables.Equals(rslt.HitObject, objLSource) {
					lightColor := calculatePhongReflectanceVec(ambientI, rslt.HitPos, objLSource, rnd, recursion+1)
					// Color component from incoming light
					intensityLambert := objColor.MultComp(lightColor)
					ldsNormal := vec3.Dot(lightDir, surfaceNormal)
					ldsNormalMax := math.Max(ldsNormal, 0)
					intensityLambert = intensityLambert.Mult(ldsNormalMax)
					intensityLambert = intensityLambert.Mult(refProps.Lambertian)
					outColorVec = outColorVec.Add(intensityLambert)

					// Color component from specular light
					viewingRay := vec3.DirFromPos(rnd.camera.Pos, hitPoint)
					reflVec := lightDir.Reverse().Add(surfaceNormal.Mult(2).Mult(vec3.Dot(lightDir, surfaceNormal))).Unit()
					intensitySpecular := specular.MultComp(lightColor).Mult(refProps.Specular)
					rdv := vec3.Dot(reflVec.Reverse(), viewingRay)
					rdvMax := math.Max(rdv, 0)
					powSmooth := math.Pow(rdvMax, refProps.Smoothness)
					intensitySpecular = intensitySpecular.Mult(powSmooth)
					// intensitySpecular = intensitySpecular
					outColorVec = outColorVec.Add(intensitySpecular)
				}
			}
		}
	}
	outColorVec = vec3.Min(outColorVec, vec3.Unit())
	// log.Info().Msgf("outColorVec: %v", outColorVec)
	return outColorVec
}

func CalculatePhongReflectance(ambientI, hitPoint vec3.Vec3, obj drawables.Drawable, rnd *Renderer) color.RGBA {
	outColorVec := calculatePhongReflectanceVec(ambientI, hitPoint, obj, rnd, 0)
	return vec3.Vec3ToRGBA(outColorVec, 255)
}
