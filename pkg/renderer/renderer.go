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

const MINIMUM_HIT_DISTANCE = 0.001
const MAXIMUM_TRACE_DISTANCE = 10000.0
const MAX_STEPS = 10000

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
}

func DefaultLightingOpts() LightingOpts {
	return LightingOpts{
		shadows: true,
		vignette: struct {
			enabled  bool
			strength float64
		}{
			enabled:  true,
			strength: 0.05,
		},
		bg: struct {
			color color.RGBA
			show  bool
		}{
			color: BG_COLOR,
			show:  true,
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

func CalculateLighting2(marchRslt MarchResult, screenPos Point, renderer *Renderer, lightingOpts LightingOpts) color.RGBA {
	pxColorVal := lightingOpts.bg.color
	pxColorVec := vec3.RGBAToVec3(lightingOpts.bg.color)
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

		pxColorVec = vec3.RGBAToVec3(marchRslt.HitObject.Color())

		pxColorVec = vec3.Min(pxColorVec.MultComp(colorVec), vec3.OfSize(1))

	}
	if lightingOpts.vignette.enabled {
		maxVignettNorm := utils.NewVec2(float64(renderer.camera.SizeX), float64(renderer.camera.SizeY)).Norm() * math.Min(1, (1-math.Min(1, lightingOpts.vignette.strength)))
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
			marchRslt := RayMarch(ray, renderer.scene)
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
			marchRslt := RayMarch(ray, renderer.scene)
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

	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Shuffle(len(points), func(i, j int) { points[i], points[j] = points[j], points[i] })
	lopts := DefaultLightingOpts()
	for _, pt := range points {
		ray := renderer.camera.RayForPixel(pt)
		marchRslt := RayMarch(ray, renderer.scene)
		// vignetteAmt := 1 - (utils.NewVec2(float64(pt.X-renderer.camera.centerOffset.X), float64(pt.Y-renderer.camera.centerOffset.Y)).Norm() / maxVignettNorm)
		pxColorVal := CalculateLighting2(marchRslt, pt, renderer, lopts)
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
	// log.Info().Msg("Workers done, encoding image to path")
	time.Sleep(time.Second)
	renderer.isDone = true
}

func DefaultScene(opts RenderOpts) *Renderer {

	// Setup Scene
	drawable1 := drawables.NewMandelB("m1", 60, 1.5, 8, vec3.Zero(), color.RGBA{240, 167, 49, 255})
	drawable2 := drawables.NewNamedSphere("s2", vec3.Vec3{X: 10, Y: -4, Z: -2}, 1, color.RGBA{237, 66, 22, 255})
	drawable3 := drawables.NewTorus(vec3.Vec3{X: 10, Y: -4, Z: -2}, 4, 0.5, color.RGBA{130, 156, 154, 255})
	// drawable4 := drawables.NewNamedCube("b1", vec3.Vec3{X: -4, Y: -2, Z: -1.5}, 1, color.RGBA{255, 255, 255, 255})

	// Setup lighting
	// light1 := drawables.NewNamedSphere("l1", vec3.Vec3{X: -10, Y: -10, Z: -10}, 1, color.RGBA{255, 0, 0, 255})
	// light2 := drawables.NewNamedSphere("l2", vec3.Vec3{X: -10, Y: 10, Z: 0}, 1, color.RGBA{0, 255, 0, 255})
	// light3 := drawables.NewNamedSphere("l3", vec3.Vec3{X: -10, Y: -10, Z: 10}, 0.5, color.RGBA{0, 0, 255, 255})
	light1 := drawables.NewNamedSphere("l1", vec3.Vec3{X: -10, Y: -10, Z: -10}, 1, color.RGBA{254, 255, 255, 255})
	light2 := drawables.NewNamedSphere("l2", vec3.Vec3{X: -10, Y: 10, Z: 0}, 1, color.RGBA{69, 79, 79, 255})
	// light3 := drawables.NewNamedSphere("l3", vec3.Vec3{X: -10, Y: -10, Z: 10}, 0.5, color.RGBA{69, 79, 79, 255})

	// Setup Scene
	scene := NewScene([]drawables.Drawable{drawable2, drawable3, drawable1}, []drawables.Drawable{light1, light2})

	// Setup Camera
	// cam := NewCameraFOV(vec3.Vec3{X: -50, Y: 0, Z: 0}, 5000, 5000, 2.75, "../../rend_out_3") // 4k

	// Standard Camera Resolutions
	// cam := NewCameraFOV(vec3.Vec3{X: -15, Y: 0, Z: 0}, 1920, 1080, 20, "./rend_out_0") // 1080p
	// cam := NewCameraFOV(vec3.Vec3{X: -15, Y: 0, Z: 0}, 3840, 2160, 20, "./rend_out_0") // 4k
	// cam := NewCameraFOV(vec3.Vec3{X: -10, Y: 0, Z: -1}, 7680, 4320, 10) // 8k
	// cam := NewCameraFOV(vec3.Vec3{X: -15, Y: 0, Z: 0}, opts.DimX, opts.DimY, opts.Fov, opts.OutPath) // 16k
	// cam := NewCameraFOV(vec3.Vec3{X: -25, Y: 0, Z: 0}, 30720, 17280, 45) // 32k

	cam := NewCameraFOV(vec3.Vec3{X: -15, Y: -1, Z: -0.1}, opts.DimX, opts.DimY, opts.Fov, opts.OutPath)

	renderer := Renderer{
		scene,
		cam,
		false,
	}
	return &renderer

	// Render(&renderer, opts.Workers)

	// renderer.camera.FlushToDisk()

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
