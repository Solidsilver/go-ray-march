package renderer

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

const MINIMUM_HIT_DISTANCE = 0.001
const MAXIMUM_TRACE_DISTANCE = 10000.0
const MAX_STEPS = 10000

// var BG_COLOR = color.RGBA{255, 255, 255, 255}
// var BG_COLOR = color.RGBA{100, 100, 100, 255}
var BG_COLOR = color.RGBA{0, 0, 0, 255}

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
	reflections bool
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
			enabled: false,
		},
		dropoff: struct {
			enabled  bool
			color    color.RGBA
			distance float64
		}{
			enabled: false,
			color: color.RGBA{
				0, 0, 0, 0,
			},
			distance: MAXIMUM_TRACE_DISTANCE / 25,
		},
		reflections: true,
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

func CalculateLightingOld(marchRslt MarchResult, renderer *Renderer) color.RGBA {
	pxColorVal := BG_COLOR
	if marchRslt.HitObject != nil {
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

		pxColorVec := vec3.RGBAToVec3(marchRslt.HitObject.Color())

		pxColorVec = vec3.Min(pxColorVec.MultComp(colorVec), vec3.OfSize(1))

		pxColorVal = vec3.Vec3ToRGBA(pxColorVec, pxColorVal.A)

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
			pxColorVal := CalculateLightingOld(marchRslt, renderer)
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
			marchRslt := RayMarch(ray, renderer.scene)
			pxColorVal := CalculateLightingOld(marchRslt, renderer)
			renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		}
	}
}

func RayMarchWorkerLighting3(id int, workers int, renderer *Renderer, pb *progressbar.ProgressBar, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info().Int("tid", id).Msg("Preparing Render")
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
	log.Info().Int("tid", id).Msg("Starting Render")
	for _, pt := range points {
		ray := renderer.camera.RayForPixel(pt)
		marchRslt := RayMarch(ray, renderer.scene)

		pxColorVal := BG_COLOR
		if marchRslt.HitObject != nil {
			iambient := vec3.RGBAToVec3(renderer.scene.options.bg.color)
			pxColorVal = CalculatePhongReflectance(iambient, marchRslt.HitPos, marchRslt.HitObject, renderer)
		}
		pxColorVal = CalculatePostProcessing(vec3.RGBAToVec3(pxColorVal), marchRslt, pt, renderer)

		// reflectionColor := CalculateReflectionColor(marchRslt, renderer)
		// log.Info().Int("tid", id).Msg("Calculating Lighting")

		// log.Info().Int("tid", id).Msg("Done with Lighting")

		renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		pb.Add(1)
		// log.Info().Int("tid", id).Msg("Done with pixel")
	}

}

func GetReflections(orgResult MarchResult, renderer *Renderer) []MarchResult {
	latestResult := orgResult
	reflections := []MarchResult{}
	for !(latestResult.ReachedMaxDistance || latestResult.ReachedMaxSteps) {
		reflectResult := CalculateReflectionsForHit(latestResult, renderer)
		reflections = append(reflections, reflectResult)
		if !reflectResult.DidHit || reflectResult == latestResult || reflectResult.HitObject.Reflectivity() == 0 {
			break
		}
		latestResult = reflectResult
	}
	return reflections
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
		drawables.NewNamedSphere("s1", vec3.Vec3{X: 3, Y: 1.5, Z: -1.5}, 1, color.RGBA{70, 150, 205, 255}, false, drawables.ReflectionProperties{
			Ambient:    0,
			Lambertian: 0.5,
			Specular:   0,
			Metalness:  0,
			Smoothness: 1,
			// Reflection: 0.5,
		}),
		drawables.NewNamedSphere("s2", vec3.Zero(), 1, color.RGBA{240, 167, 49, 255}, false),

		// drawables.NewMandelB("m1", 60, 1.5, 12, vec3.Zero(), color.RGBA{240, 167, 49, 255}, false, 1),
		// drawables.NewMandelB("m2", 60, 1.5, 12, vec3.Zero(), color.RGBA{25, 35, 45, 255}, false),
		//drawables.NewNamedCube("b2", vec3.Vec3{X: 10, Y: -4, Z: 2}, .65, color.RGBA{237, 66, 22, 255}),
		// drawables.NewNamedTorus("t1", vec3.Vec3{X: 10, Y: -4, Z: -2}, 4, 0.25, color.RGBA{130, 156, 154, 255}),
		//drawables.NewNamedCube("b1", vec3.Vec3{X: -4, Y: -2, Z: -1.5}, 1, color.RGBA{255, 255, 255, 255}),
	)

	scene.AddLights(
		drawables.NewNamedLight("l1", vec3.Vec3{X: -135, Y: -1, Z: -1}, 1, color.RGBA{255, 255, 255, 255}, false),
		drawables.NewNamedLight("l2", vec3.Vec3{X: -135, Y: 1, Z: 1}, 1, color.RGBA{255, 255, 255, 128}, false),
		// drawables.NewNamedLight("l5", vec3.Vec3{X: 15, Y: -8, Z: -8}, 1, color.RGBA{255, 255, 255, 100}, false),
		// drawables.NewNamedSphere("l2", vec3.Vec3{X: -15, Y: 8, Z: 8}, 1, color.RGBA{0, 255, 0, 255}, false),
		//drawables.NewNamedSphere("l3", vec3.Vec3{X: -15, Y: -8, Z: 8}, 0.5, color.RGBA{0, 0, 255, 255}, false),
		//drawables.NewNamedSphere("l3", vec3.Vec3{X: -10, Y: -10, Z: 10}, 0.5, color.RGBA{69, 79, 79, 255}),
	)

	cam := NewCameraFOV(vec3.Vec3{X: -15, Y: 0, Z: 0}, opts.DimX, opts.DimY, opts.Fov, opts.OutPath)

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
