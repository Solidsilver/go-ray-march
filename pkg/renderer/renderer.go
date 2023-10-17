package renderer

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/vec"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"goki.dev/mat32/v2"
)

const MINIMUM_HIT_DISTANCE = 0.0001
const MAX_HIT_DISTANCE = 10.0
const MAXIMUM_TRACE_DISTANCE = 5000.0
const MAX_STEPS = 100000
const LOD = true

// var BG_COLOR = color.RGBA{255, 255, 255, 255}
var BG_COLOR = color.RGBA{198, 226, 253, 255}

type Ray struct {
	origin mat32.Vec3
	dir    mat32.Vec3
}

type Renderer struct {
	scene  *Scene
	camera *Camera
	isDone atomic.Bool
	Reset  atomic.Bool
}

type LightingOpts struct {
	shadows  bool
	vignette struct {
		enabled  bool
		strength float32
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
		distance float32
	}
}

func DefaultLightingOpts() LightingOpts {
	return LightingOpts{
		shadows: true,
		vignette: struct {
			enabled  bool
			strength float32
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
			distance float32
		}{
			enabled: true,
			color: color.RGBA{
				100, 100, 100, 255,
			},
			distance: MAXIMUM_TRACE_DISTANCE / 25,
		},
	}
}

func NewRenderer(scene *Scene, camera *Camera) Renderer {
	return Renderer{scene, camera, atomic.Bool{}, atomic.Bool{}}
}

func (r *Renderer) IsDone() bool {
	return r.isDone.Load()
}

func (r *Renderer) SetDone(val bool) {
	r.isDone.Store(val)
}

func (r *Renderer) GetCamera() *Camera {
	return r.camera
}

func (r *Renderer) GetScene() *Scene {
	return r.scene
}

func CalculateLighting(marchRslt MarchResult, renderer *Renderer) color.RGBA {
	pxColorVal := BG_COLOR
	if marchRslt.HitObject != nil {
		hitPoint := marchRslt.HitPos
		colorVec := mat32.Vec3Zero
		for _, lSource := range renderer.scene.Lights {
			lightDir := vec.DirFromPos(lSource.Pos(), hitPoint)
			surfaceNormal := SurfaceNormal(hitPoint, marchRslt.HitObject, marchRslt.Mhd)
			bounceDeg := mat32.RadToDeg(lightDir.AngleTo(surfaceNormal))
			if bounceDeg < 90 {
				ray := Ray{hitPoint, lightDir}
				rslt := RayMarch(ray, renderer)
				if drawables.Equals(rslt.HitObject, lSource) {
					brightness := float32(rslt.HitObject.Color().A) / 255.0
					brightness = brightness * (90 - bounceDeg) / 90
					lightColorVec := vec.RGBAToVec3(lSource.Color()).MulScalar(brightness)
					colorVec = colorVec.Add(lightColorVec)

				}
			}
		}

		pxColorVec := vec.RGBAToVec3(marchRslt.HitObject.Color())

		pxColorVec = pxColorVec.Mul(colorVec)
		pxColorVec.ClampScalar(0, 1)

		pxColorVal = vec.Vec3ToRGBA(pxColorVec, pxColorVal.A)

	}
	return pxColorVal
}

func CalculateLighting2(marchRslt MarchResult, screenPos mat32.Vec2, renderer *Renderer) color.RGBA {
	pxColorVal := renderer.scene.options.bg.color
	pxColorVec := vec.RGBAToVec3(renderer.scene.options.bg.color)
	if marchRslt.HitObject != nil {
		pxColorVec = vec.RGBAToVec3(marchRslt.HitObject.Color())
		if renderer.scene.options.shadows {
			hitPoint := marchRslt.HitPos
			colorVec := mat32.Vec3Zero
			for _, lSource := range renderer.scene.Lights {
				lightDir := vec.DirFromPos(lSource.Pos(), hitPoint).Normal()
				surfaceNormal := SurfaceNormal(hitPoint, marchRslt.HitObject, marchRslt.Mhd)
				bounceDeg := mat32.RadToDeg(mat32.Abs(lightDir.AngleTo(surfaceNormal)))
				if bounceDeg < 90 {
					ray := Ray{hitPoint, lightDir}
					rslt := RayMarch(ray, renderer)
					if drawables.Equals(rslt.HitObject, lSource) {
						brightness := float32(rslt.HitObject.Color().A) / 255.0
						brightness = brightness * (90 - bounceDeg) / 90
						lightColorVec := vec.RGBAToVec3(lSource.Color()).MulScalar(brightness)
						colorVec = colorVec.Add(lightColorVec)
					}
				}
			}

			pxColorVec = pxColorVec.Mul(colorVec)
			pxColorVec.ClampScalar(0, 1)
		}

	}
	if renderer.scene.options.ao.enabled && marchRslt.HitObject != nil {
		ao := 1.0 - float32(marchRslt.Steps)/float32(MAX_STEPS-1)
		pxColorVec = pxColorVec.MulScalar(ao)

	}

	if renderer.scene.options.dropoff.enabled {
		dropoffDist := mat32.Min(renderer.scene.options.dropoff.distance, MAXIMUM_TRACE_DISTANCE)
		distFrac := mat32.Min((marchRslt.Distance / dropoffDist), 1)
		dropoff := float32(1 - mat32.Pow(distFrac, 2))
		blendColor := vec.RGBAToVec3(renderer.scene.options.dropoff.color)
		pxColorVec = pxColorVec.MulScalar(dropoff).Add(blendColor.MulScalar(1 - dropoff))

	}

	if renderer.scene.options.vignette.enabled {
		maxVignettNorm := mat32.NewVec2(float32(renderer.camera.SizeX), float32(renderer.camera.SizeY)).Length() * float32(mat32.Min(1, (1-mat32.Min(1, renderer.scene.options.vignette.strength))))
		vignettAmt := 1 - (mat32.NewVec2(screenPos.X-renderer.camera.centerOffset.X, (screenPos.Y-renderer.camera.centerOffset.Y)).Length() / maxVignettNorm)
		pxColorVec = pxColorVec.MulScalar(vignettAmt)
	}

	pxColorVal = vec.Vec3ToRGBA(pxColorVec, pxColorVal.A)
	return pxColorVal
}

func RayMarchWorkerLighting(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := id; i <= renderer.camera.SizeX; i += workers {
		for j := 0; j <= renderer.camera.SizeY; j++ {
			pt := mat32.NewVec2(float32(i), float32(j))
			ray := renderer.camera.RayForPixel(pt)
			marchRslt := RayMarch(ray, renderer)
			pxColorVal := CalculateLighting(marchRslt, renderer)
			renderer.camera.Image.Set(i, j, pxColorVal)
		}
	}
}

func RayMarchWorkerLighting2(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i <= renderer.camera.SizeX; i++ {
		for j := id; j <= renderer.camera.SizeY; j += workers {
			j2 := (j + i) % renderer.camera.SizeY
			pt := mat32.NewVec2(float32(i), float32(j2))
			ray := renderer.camera.RayForPixel(pt)
			marchRslt := RayMarch(ray, renderer)
			pxColorVal := CalculateLighting(marchRslt, renderer)
			renderer.camera.Image.Set(int(pt.X), int(pt.Y), pxColorVal)
		}
	}
}

// func RayMarchWorkerLighting3(id int, workers int, renderer *Renderer, pb *progressbar.ProgressBar, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	points := make([]Point, renderer.camera.Size()/int64(workers))
// 	count := 0
// 	for i := int64(id); i < renderer.camera.Size(); i += int64(workers) {
// 		y := (i) % int64(renderer.camera.SizeY)
// 		x := i / int64(renderer.camera.SizeY)
// 		pt := Point{int(x), int(y)}
// 		points[count] = pt
// 		count++
// 	}

// 	rand.Shuffle(len(points), func(i, j int) { points[i], points[j] = points[j], points[i] })
// 	for _, pt := range points {
// 		if renderer.Reset.Load() {
// 			// fmt.Printf("Resetting worker %d\n", id)
// 			return
// 		}
// 		ray := renderer.camera.RayForPixel(pt)
// 		marchRslt := RayMarch(ray, renderer)
// 		pxColorVal := CalculateLighting2(marchRslt, pt, renderer)
// 		renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
// 		pb.Add(1)
// 	}

// }

func RayMarchWorkerLighting3(id int, workers int, renderer *Renderer, pb *progressbar.ProgressBar, wg *sync.WaitGroup) {
	defer wg.Done()

	rangeSize := renderer.camera.Size() / int64(workers)

	for _, value := range rand.Perm(int(rangeSize)) {
		if renderer.Reset.Load() {
			return
		}
		i := int64(value)*int64(workers) + int64(id)
		y := (i) % int64(renderer.camera.SizeY)
		x := i / int64(renderer.camera.SizeY)
		pt := mat32.NewVec2(float32(x), float32(y))
		ray := renderer.camera.RayForPixel(pt)
		marchRslt := RayMarch(ray, renderer)
		pxColorVal := CalculateLighting2(marchRslt, pt, renderer)
		renderer.camera.Image.Set(int(pt.X), int(pt.Y), pxColorVal)
		pb.Add(1)
	}
}

func RayMarchWorkerLighting6(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()

	rangeSize := int(renderer.camera.Size() / int64(workers))

	for _, value := range rand.Perm(rangeSize) {
		if renderer.Reset.Load() {
			return
		}
		i := value*workers + id
		y := i % renderer.camera.SizeY
		x := i / renderer.camera.SizeY
		pt := mat32.NewVec2(float32(x), float32(y))
		ray := renderer.camera.RayForPixel(pt)
		marchRslt := RayMarch(ray, renderer)
		pxColorVal := CalculateLighting2(marchRslt, pt, renderer)
		renderer.camera.Image.Set(int(pt.X), int(pt.Y), pxColorVal)
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
	// time.Sleep(time.Second)
	renderer.isDone.Store(true)
	renderer.Reset.Store(false)
}

func (renderer *Renderer) Render2(workers int, wg *sync.WaitGroup) {
	renderer.isDone.Store(false)
	startTime := time.Now()

	wg.Add(1)
	defer wg.Done()

	var wg2 sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg2.Add(1)
		go RayMarchWorkerLighting6(i, workers, renderer, &wg2)
	}

	log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")
	wg2.Wait()

	renderer.Reset.Store(false)
	renderer.isDone.Store(true)
	renderDuration := time.Since(startTime)
	fmt.Printf("Rendered frame in %s\n", renderDuration.String())

}

func NewDefaultRenderScene(opts RenderOpts) *Renderer {

	// Setup Scene
	scene := NewBlankScene()
	scene.AddDrawables(
		// drawables.NewNamedSphere("s2", mat32.Vec3{X: 0, Y: 0, Z: 0}, 1, color.RGBA{185, 134, 247, 255}, true),
		drawables.NewMandelB("m1", 90, 5, -8, mat32.Vec3Zero, color.RGBA{135, 134, 247, 255}, false),
		// drawables.NewMandelB("m2", 60, 1.5, 12, mat32.Zero(), color.RGBA{25, 35, 45, 255}, false),
		// drawables.NewNamedCube("b2", mat32.Vec3{X: 10, Y: -4, Z: 2}, .65, color.RGBA{237, 66, 22, 255}),
		// drawables.NewNamedTorus("t1", mat32.Vec3{X: 10, Y: -4, Z: -2}, 4, 0.25, color.RGBA{130, 156, 154, 255}),
		//drawables.NewNamedCube("b1", mat32.Vec3{X: -4, Y: -2, Z: -1.5}, 1, color.RGBA{255, 255, 255, 255}),
	)

	scene.AddLights(
		drawables.NewNamedSphere("l1", mat32.Vec3{X: -15, Y: -1, Z: -1}, 1, color.RGBA{150, 150, 150, 255}, false),
		// drawables.NewNamedSphere("l2", mat32.Vec3{X: -15, Y: 1, Z: 1}, 1, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewNamedSphere("l5", mat32.Vec3{X: -15, Y: -8, Z: -8}, 1, color.RGBA{200, 200, 200, 255}, false),
		// drawables.NewNamedSphere("l2", mat32.Vec3{X: -15, Y: 8, Z: 8}, 1, color.RGBA{0, 255, 0, 255}, false),
		// drawables.NewNamedSphere("l3", mat32.Vec3{X: -15, Y: -8, Z: 8}, 0.5, color.RGBA{0, 0, 255, 255}, false),
		//drawables.NewNamedSphere("l3", mat32.Vec3{X: -10, Y: -10, Z: 10}, 0.5, color.RGBA{69, 79, 79, 255}),
		// drawables.NewNamedSphere("l1", mat32.Vec3{X: -1, Y: -1, Z: -15}, 1, color.RGBA{240, 240, 240, 255}, false),
		// drawables.NewNamedSphere("l5", mat32.Vec3{X: -8, Y: -8, Z: -15}, 1, color.RGBA{200, 200, 200, 255}, false),
		// drawables.NewNamedSphere("l1", mat32.Vec3{X: 1, Y: 1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", mat32.Vec3{X: -1, Y: -1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", mat32.Vec3{X: 1, Y: -1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", mat32.Vec3{X: -1, Y: 1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", mat32.Vec3{X: 0, Y: 0, Z: 17}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l5", mat32.Vec3{X: 8, Y: 8, Z: 15}, 1, color.RGBA{200, 200, 200, 255}, false),
	)

	cam := NewCameraFOV(mat32.Vec3{X: -15, Y: 0, Z: 0}, opts.DimX, opts.DimY, opts.Fov, opts.OutPath)

	cam.up = mat32.Vec3Z
	cam.Dir = mat32.Vec3X //.MulScalar(-1)

	renderer := Renderer{
		scene,
		cam,
		atomic.Bool{},
		atomic.Bool{},
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
