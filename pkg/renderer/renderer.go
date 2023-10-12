package renderer

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

const MINIMUM_HIT_DISTANCE = 0.00001
const MAX_HIT_DISTANCE = 0.5
const MAXIMUM_TRACE_DISTANCE = 10000.0
const MAX_STEPS = 10000
const LOD = true

// var BG_COLOR = color.RGBA{255, 255, 255, 255}
var BG_COLOR = color.RGBA{198, 226, 253, 255}

type Ray struct {
	origin vec3.Vec3
	dir    vec3.Vec3
}

type Renderer struct {
	scene  *Scene
	camera *Camera
	isDone bool
}

type LightingOpts struct {
	shadows  bool
	vignette struct {
		enabled  bool
		strength float64
	}
	bg struct {
		color color.RGBA
		show  bool
	}
	ao struct {
		enabled bool
	}
	dropoff struct {
		enabled  bool
		color    color.RGBA
		distance float64
	}
}

func DefaultLightingOpts() LightingOpts {
	return LightingOpts{
		shadows: true,
		vignette: struct {
			enabled  bool
			strength float64
		}{
			enabled:  false,
			strength: 0.05,
		},
		bg: struct {
			color color.RGBA
			show  bool
		}{
			color: BG_COLOR,
			show:  true,
		},
		ao: struct {
			enabled bool
		}{
			enabled: true,
		},
		dropoff: struct {
			enabled  bool
			color    color.RGBA
			distance float64
		}{
			enabled: true,
			color: color.RGBA{
				0, 0, 0, 0,
			},
			distance: MAXIMUM_TRACE_DISTANCE / 25,
		},
	}
}

func NewRenderer(scene *Scene, camera *Camera) Renderer {
	return Renderer{scene, camera, false}
}

func (r *Renderer) GetStatus() bool {
	return r.isDone
}

func (r Renderer) GetCamera() *Camera {
	return r.camera
}

func (r Renderer) GetScene() *Scene {
	return r.scene
}

func CalculateLighting(marchRslt MarchResult, renderer *Renderer) color.RGBA {
	pxColorVal := BG_COLOR
	if marchRslt.HitObject != nil {
		hitPoint := marchRslt.HitPos
		colorVec := vec3.Zero()
		for _, lSource := range renderer.scene.Lights {
			lightDir := vec3.DirFromPos(lSource.Pos(), hitPoint).Unit()
			surfaceNormal := SurfaceNormal(hitPoint, marchRslt.HitObject, marchRslt.Mhd)
			bounceDeg := vec3.Angle(lightDir, surfaceNormal)
			if bounceDeg < 90 {
				ray := Ray{hitPoint, lightDir}
				rslt := RayMarch(ray, renderer)
				if drawables.Equals(rslt.HitObject, lSource) {
					brightness := float64(rslt.HitObject.Color().A) / 255
					brightness = brightness * (90 - bounceDeg) / 90
					lightColorVec := vec3.RGBAToVec3(lSource.Color()).Mult(brightness)
					colorVec = colorVec.Add(lightColorVec)

				}
			}
		}

		pxColorVec := vec3.RGBAToVec3(marchRslt.HitObject.Color())

		pxColorVec = vec3.Min(pxColorVec.MultComp(colorVec), vec3.OfSize(1))

		pxColorVal = vec3.Vec3ToRGBA(pxColorVec, pxColorVal.A)

	}
	return pxColorVal
}

func CalculateLighting2(marchRslt MarchResult, screenPos Point, renderer *Renderer) color.RGBA {
	pxColorVal := renderer.scene.options.bg.color
	pxColorVec := vec3.RGBAToVec3(renderer.scene.options.bg.color)
	if marchRslt.HitObject != nil {
		pxColorVec = vec3.RGBAToVec3(marchRslt.HitObject.Color())
		if renderer.scene.options.shadows {
			hitPoint := marchRslt.HitPos
			colorVec := vec3.Zero()
			for _, lSource := range renderer.scene.Lights {
				lightDir := vec3.DirFromPos(lSource.Pos(), hitPoint).Unit()
				surfaceNormal := SurfaceNormal(hitPoint, marchRslt.HitObject, marchRslt.Mhd)
				bounceDeg := vec3.Angle(lightDir, surfaceNormal)
				if bounceDeg < 90 {
					ray := Ray{hitPoint, lightDir}
					rslt := RayMarch(ray, renderer)
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

	pxColorVal = vec3.Vec3ToRGBA(pxColorVec, pxColorVal.A)
	return pxColorVal
}

func RayMarchWorkerLighting(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := id; i <= renderer.camera.SizeX; i += workers {
		for j := 0; j <= renderer.camera.SizeY; j++ {
			pt := Point{i, j}
			ray := renderer.camera.RayForPixel(pt)
			marchRslt := RayMarch(ray, renderer)
			pxColorVal := CalculateLighting(marchRslt, renderer)
			renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		}
	}
}

func RayMarchWorkerLighting2(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i <= renderer.camera.SizeX; i++ {
		for j := id; j <= renderer.camera.SizeY; j += workers {
			j2 := (j + i) % renderer.camera.SizeY
			pt := Point{i, j2}
			ray := renderer.camera.RayForPixel(pt)
			marchRslt := RayMarch(ray, renderer)
			pxColorVal := CalculateLighting(marchRslt, renderer)
			renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		}
	}
}

func RayMarchWorkerLighting3(id int, workers int, renderer *Renderer, pb *progressbar.ProgressBar, wg *sync.WaitGroup) {
	defer wg.Done()
	points := make([]Point, renderer.camera.Size()/int64(workers))
	count := 0
	for i := int64(id); i < renderer.camera.Size(); i += int64(workers) {
		y := (i) % int64(renderer.camera.SizeY)
		x := i / int64(renderer.camera.SizeY)
		pt := Point{int(x), int(y)}
		points[count] = pt
		count++
	}

	rand.Shuffle(len(points), func(i, j int) { points[i], points[j] = points[j], points[i] })
	for _, pt := range points {
		ray := renderer.camera.RayForPixel(pt)
		marchRslt := RayMarch(ray, renderer)
		pxColorVal := CalculateLighting2(marchRslt, pt, renderer)
		renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		pb.Add(1)
	}

}

func RenderOut(renderer *Renderer, workers int) {
	Render(renderer, workers)
	renderer.camera.FlushToDisk()
}

func Render(renderer *Renderer, workers int) {
	wg := new(sync.WaitGroup)
	pb := progressbar.NewOptions64(renderer.camera.Size(),
		progressbar.OptionSetDescription("Rendering Image..."),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("px"),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionUseANSICodes(true),
	)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go RayMarchWorkerLighting3(i, workers, renderer, pb, wg)
	}

	log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")
	wg.Wait()
	time.Sleep(time.Second)
	renderer.isDone = true
}

func NewDefaultRenderScene(opts RenderOpts) *Renderer {

	// Setup Scene
	scene := NewBlankScene()
	scene.AddDrawables(
		// drawables.NewNamedSphere("s2", vec3.Vec3{X: 10, Y: 5, Z: 1}, 1, color.RGBA{70, 150, 205, 255}, true),
		drawables.NewMandelB("m1", 60, 1.5, 12, vec3.Zero(), color.RGBA{255, 255, 255, 255}, true),
		// drawables.NewMandelB("m2", 60, 1.5, 12, vec3.Zero(), color.RGBA{25, 35, 45, 255}, false),
		//drawables.NewNamedCube("b2", vec3.Vec3{X: 10, Y: -4, Z: 2}, .65, color.RGBA{237, 66, 22, 255}),
		// drawables.NewNamedTorus("t1", vec3.Vec3{X: 10, Y: -4, Z: -2}, 4, 0.25, color.RGBA{130, 156, 154, 255}),
		//drawables.NewNamedCube("b1", vec3.Vec3{X: -4, Y: -2, Z: -1.5}, 1, color.RGBA{255, 255, 255, 255}),
	)

	scene.AddLights(
		drawables.NewNamedSphere("l1", vec3.Vec3{X: -15, Y: -1, Z: -1}, 1, color.RGBA{240, 240, 240, 255}, false),
		// drawables.NewNamedSphere("l2", vec3.Vec3{X: -15, Y: 1, Z: 1}, 1, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewNamedSphere("l5", vec3.Vec3{X: -15, Y: -8, Z: -8}, 1, color.RGBA{200, 200, 200, 255}, false),
		// drawables.NewNamedSphere("l2", vec3.Vec3{X: -15, Y: 8, Z: 8}, 1, color.RGBA{0, 255, 0, 255}, false),
		//drawables.NewNamedSphere("l3", vec3.Vec3{X: -15, Y: -8, Z: 8}, 0.5, color.RGBA{0, 0, 255, 255}, false),
		//drawables.NewNamedSphere("l3", vec3.Vec3{X: -10, Y: -10, Z: 10}, 0.5, color.RGBA{69, 79, 79, 255}),
	)

	cam := NewCameraFOV(vec3.Vec3{X: -12, Y: -0.20, Z: -0.8}, opts.DimX, opts.DimY, opts.Fov, opts.OutPath)

	renderer := Renderer{
		scene,
		cam,
		false,
	}
	return &renderer

}

type RenderOpts struct {
	Workers int
	OutPath string
	DimX    int
	DimY    int
	Fov     float64
}

func DefaultRenderOpts() RenderOpts {
	return RenderOpts{
		Workers: 1,
		OutPath: "./rend_out_0",
		DimX:    1920,
		DimY:    1080,
		Fov:     20,
	}
}

func (opts RenderOpts) String() string {
	return fmt.Sprintf("Threads: %d, OutPath: %s, Dim: %dx%d, Fov: %0.2f", opts.Workers, opts.OutPath, opts.DimX, opts.DimY, opts.Fov)
}
