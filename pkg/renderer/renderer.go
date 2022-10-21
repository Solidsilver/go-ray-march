package renderer

import (
	"image/color"
	"math"
	"sync"

	"github.com/rs/zerolog/log"
	"solidsilver.dev/go-ray-marching/pkg/drawables"
	"solidsilver.dev/go-ray-marching/pkg/vec3"
)

const MINIMUM_HIT_DISTANCE = 0.05
const MAXIMUM_TRACE_DISTANCE = 10000.0

var BG_COLOR = color.RGBA{0, 0, 0, 255}

type Ray struct {
	origin vec3.Vec3
	dir    vec3.Vec3
}

type Renderer struct {
	scene  *Scene
	camera *Camera
}

func CalculateLighting(marchRslt MarchResult, renderer *Renderer) color.RGBA {
	pxColorVal := BG_COLOR
	if marchRslt.HitObject != nil {
		hitPoint := marchRslt.HitPos
		pxColorVal = marchRslt.HitObject.Color()
		totalBrightness := 0.0
		for _, lSource := range renderer.scene.Lights {
			lightDir := vec3.DirFromPos(lSource.Pos(), hitPoint).Unit()
			surfaceNormal := SurfaceNormal(hitPoint, marchRslt.HitObject)
			bounceDeg := vec3.Angle2(lightDir, surfaceNormal)
			if bounceDeg < 90 {
				ray := Ray{hitPoint, lightDir}
				rslt := RayMarch(ray, renderer.scene)
				if drawables.Equals(rslt.HitObject, lSource) {
					brightness := float64(rslt.HitObject.Color().A)
					brightness = brightness * (90 - bounceDeg) / 90
					totalBrightness += brightness
				}
			}
		}
		brightScaled := totalBrightness / 255
		pxColorVal = color.RGBA{
			R: uint8(math.Min(255, float64(pxColorVal.R)*(brightScaled))),
			G: uint8(math.Min(255, float64(pxColorVal.G)*(brightScaled))),
			B: uint8(math.Min(255, float64(pxColorVal.B)*(brightScaled))),
			A: pxColorVal.A,
		}
	}
	return pxColorVal
}

func RayMarchWorkerLighting(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := id; i <= renderer.camera.SizeX; i += workers {
		for j := 0; j <= renderer.camera.SizeY; j++ {
			pt := Point{i, j}
			ray := renderer.camera.RayForPixel(pt)
			marchRslt := RayMarch(ray, renderer.scene)
			pxColorVal := CalculateLighting(marchRslt, renderer)
			renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		}
	}
}

func Render(renderer *Renderer, workers int) {
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go RayMarchWorkerLighting(i, workers, renderer, wg)
	}

	log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")
	wg.Wait()
	log.Info().Msg("Workers done, encoding image to path")
	renderer.camera.FlushToDisk()
}

func RenderDefault(workers int) {
	drawable1 := drawables.NewNamedSphere("s1", vec3.Vec3{X: 0, Y: 0, Z: 0}, 2.5, color.RGBA{100, 200, 200, 255})
	drawable2 := drawables.NewNamedSphere("s2", vec3.Vec3{X: 8, Y: -1, Z: -2}, 3.5, color.RGBA{252, 102, 11, 255})
	drawable3 := drawables.NewNamedSphere("s3", vec3.Vec3{X: 40, Y: 0, Z: 0}, 4.5, color.RGBA{76, 96, 218, 255})
	// drawable4 := drawables.NewNamedSphere("d4", vec3.Vec3{X: 2, Y: 4, Z: 4}, 1, color.RGBA{1, 123, 6, 255})
	drawable4 := drawables.NewNamedCube("b1", vec3.Vec3{X: 2, Y: -4, Z: -4}, 1, color.RGBA{1, 123, 6, 255})
	// cam := NewCameraFOV(vec3.Vec3{X: -25, Y: 0, Z: 0}, 1920, 1080, 45) // 1080p
	cam := NewCameraFOV(vec3.Vec3{X: -25, Y: 0, Z: 0}, 3840, 2160, 25) // 4k
	// cam := NewCameraFOV(vec3.Vec3{X: -10, Y: 0, Z: 0}, 7680, 4320, 45) // 8k
	// cam := NewCameraFOV(vec3.Vec3{X: -20, Y: 0, Z: 0}, 15360, 8640, 45) // 16k
	// cam := NewCameraFOV(vec3.Vec3{X: -25, Y: 0, Z: 0}, 30720, 17280, 45) // 32k

	light1 := drawables.NewNamedSphere("l1", vec3.Vec3{X: -50, Y: -50, Z: -50}, 1, color.RGBA{0, 0, 0, 255})
	light2 := drawables.NewNamedSphere("l2", vec3.Vec3{X: 50, Y: 50, Z: 50}, 1, color.RGBA{0, 0, 0, 255})
	light3 := drawables.NewNamedSphere("l3", vec3.Vec3{X: -50, Y: 0, Z: 0}, 1, color.RGBA{0, 0, 0, 100})

	scene := NewScene([]drawables.Drawable{drawable1, drawable2, drawable3, drawable4}, []drawables.Drawable{light1, light2, light3})

	renderer := Renderer{
		scene,
		cam,
	}

	Render(&renderer, workers)

}
