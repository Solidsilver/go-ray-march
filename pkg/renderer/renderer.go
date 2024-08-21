package renderer

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
	"github.com/rs/zerolog/log"
	pb "github.com/schollz/progressbar/v3"
)

// const MINIMUM_HIT_DISTANCE = 0.0001
// const MAX_HIT_DISTANCE = 10.0
// const MAXIMUM_TRACE_DISTANCE = 5000.0
// const MAX_STEPS = 100000

// var MAX_AMBIENT_STEPS = math.Sqrt(MINIMUM_HIT_DISTANCE*10)/10 + 150

var BG_COLOR = color.RGBA{0, 0, 0, 255}

// var BG_COLOR = color.RGBA{198, 226, 253, 255}

type Ray struct {
	origin vec3.Vec3
	dir    vec3.Vec3
}

type Renderer struct {
	scene  *Scene
	camera *Camera
	isDone atomic.Bool
	Reset  atomic.Bool
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

func (r *Renderer) UpdateCamera(f func(c *Camera)) {
	f(r.camera)
}

func (r *Renderer) GetScene() *Scene {
	return r.scene
}

func CalculateLighting(marchRslt MarchResult, renderer *Renderer) color.RGBA {
	pxColorVal := BG_COLOR
	if marchRslt.HitObject != nil {
		hitPoint := marchRslt.HitPos
		colorVec := vec3.Zero
		for _, lSource := range renderer.scene.Lights {
			lightDir := vec3.DirFromPos(lSource.Pos(), hitPoint)
			surfaceNormal := SurfaceNormal(marchRslt, false)
			bounceDeg := vec3.Angle(lightDir, surfaceNormal)
			if bounceDeg < 90 {
				ray := Ray{hitPoint, lightDir}
				rslt := RayMarch(ray, renderer, true)
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
			colorVec := vec3.Zero
			for _, lSource := range renderer.scene.Lights {
				lightDir := vec3.DirFromPos(lSource.Pos(), hitPoint)
				surfaceNormal := SurfaceNormal(marchRslt, renderer.scene.options.trace.fastMath)
				bounceDeg := vec3.Angle(lightDir, surfaceNormal)
				if bounceDeg < 90 {
					ray := Ray{hitPoint, lightDir}
					rslt := RayMarch(ray, renderer, true)
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
		if renderer.scene.options.ao.enabled {
			var ao float64
			if renderer.scene.options.ao.inverted {
				ao = math.Min(float64(marchRslt.Steps)/float64(renderer.scene.options.ao.maxSteps-1), 0.95) - 1.0
			} else {
				ao = 1.0 - math.Min(float64(marchRslt.Steps)/float64(renderer.scene.options.ao.maxSteps-1), 0.95)
			}

			pxColorVec = pxColorVec.Mult(ao)

		}

	}
	// if renderer.scene.options.ao.enabled && marchRslt.HitObject != nil {
	// 	ao := 1.0 - math.Min(float64(marchRslt.Steps)/float64(MAX_AMBIENT_STEPS-1), 1)
	// 	pxColorVec = pxColorVec.Mult(ao)

	// }

	if renderer.scene.options.dropoff.enabled {
		dropoffDist := math.Min(renderer.scene.options.dropoff.distance, renderer.scene.options.trace.maxDist)
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

func CalculateLightingTest(marchRslt MarchResult, screenPos Point, renderer *Renderer) color.RGBA {
	pxColorVal := renderer.scene.options.bg.color
	opts := renderer.scene.options
	if (opts.ao.enabled && marchRslt.Steps >= int(opts.ao.maxSteps)) || (opts.dropoff.enabled && marchRslt.Distance >= opts.dropoff.distance) {
		return pxColorVal
	}
	pxColorVec := vec3.RGBAToVec3P(opts.bg.color)
	if marchRslt.HitObject != nil {
		pxColorVec = vec3.NewCp(marchRslt.HitObject.ColorVec())
		if opts.shadows {
			colorVec := vec3.NewP(0, 0, 0)
			for _, lSource := range renderer.scene.Lights {
				lightDir := vec3.DirFromPos(lSource.Pos(), marchRslt.HitPos)
				surfaceNormal := SurfaceNormal(marchRslt, opts.trace.fastMath)
				brightness := vec3.Dot(surfaceNormal, lightDir)
				if brightness > 0 {
					ray := Ray{marchRslt.HitPos, lightDir}
					rslt := RayMarchP(ray, renderer, true)
					if drawables.Equals(rslt.HitObject, lSource) {
						lightColorVec := vec3.NewCp(lSource.ColorVec())
						lightColorVec.MultSet(brightness)
						colorVec.AddSet(lightColorVec)
					}
				}
			}

			pxColorVec.MultCompSet(colorVec)
			pxColorVec.MinSet(vec3.NewOfSizeP(1))
		}

		if opts.ao.enabled {
			var aoStrength float64
			if opts.ao.inverted {
				aoStrength = math.Min(float64(marchRslt.Steps)/(opts.ao.maxSteps-1.0), 0.95) - 1.0
			} else {
				aoStrength = 1.0 - math.Min(float64(marchRslt.Steps)/(opts.ao.maxSteps-1.0), 0.95)
			}

			pxColorVec.MultSet(aoStrength)
		}

	}

	if opts.dropoff.enabled {
		dropoffDist := math.Min(opts.dropoff.distance, opts.trace.maxDist)
		distFrac := math.Min((marchRslt.Distance)/dropoffDist, 1)
		dropoff := 1 - math.Pow(distFrac, 2)
		blendColor := vec3.RGBAToVec3P(opts.dropoff.color)
		pxColorVec.MultSet(dropoff)
		blendColor.MultSet(1 - dropoff)
		pxColorVec.AddSet(blendColor)
	}

	if opts.vignette.enabled {
		maxVignettNorm := renderer.camera.Dim().Norm() * math.Min(1, (1-math.Min(1, opts.vignette.strength)))
		vignettAmt := 1 - (utils.NewVec2(screenPos.X-renderer.camera.centerOffset.X, screenPos.Y-renderer.camera.centerOffset.Y).Norm() / maxVignettNorm)
		pxColorVec.MultSet(vignettAmt)
	}

	pxColorVal = vec3.Vec3ToRGBA(*pxColorVec, pxColorVal.A)
	return pxColorVal
}

func RayMarchWorkerLighting(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := id; i <= renderer.camera.SizeX; i += workers {
		for j := 0; j <= renderer.camera.SizeY; j++ {
			pt := Point{i, j}
			ray := renderer.camera.RayForPixel(pt)
			marchRslt := RayMarch(ray, renderer, true)
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
			marchRslt := RayMarch(ray, renderer, true)
			pxColorVal := CalculateLighting(marchRslt, renderer)
			renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		}
	}
}

func RayMarchWorkerLighting3(id int, workers int, renderer *Renderer, pb *pb.ProgressBar, wg *sync.WaitGroup) {
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
		if renderer.Reset.Load() {
			return
		}
		ray := renderer.camera.RayForPixel(pt)
		marchRslt := RayMarch(ray, renderer, false)
		pxColorVal := CalculateLighting2(marchRslt, pt, renderer)
		renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		pb.Add(1)
	}

}

var rmwRangeMap = make(map[int][]int)
var rngMtx sync.RWMutex

func setRangePerm(id int, val []int) {
	rngMtx.Lock()
	rmwRangeMap[id] = val
	rngMtx.Unlock()
}

func getRangePerm(id int) (val []int, ok bool) {
	rngMtx.RLock()
	val, ok = rmwRangeMap[id]
	rngMtx.RUnlock()
	return
}

func RayMarchWorkerLighting6(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()

	perm, ok := getRangePerm(id)
	if !ok {
		rangeSize := int(renderer.camera.Size() / int64(workers))
		perm = rand.Perm(rangeSize)
		setRangePerm(id, perm)
	}

	for _, value := range perm {
		if renderer.Reset.Load() {
			return
		}
		i := int64(value)*int64(workers) + int64(id)
		y := (i) % int64(renderer.camera.SizeY)
		x := i / int64(renderer.camera.SizeY)
		pt := Point{int(x), int(y)}
		ray := renderer.camera.RayForPixel(pt)
		// marchRslt := RayMarch(ray, renderer, false)
		// pxColorVal := CalculateLighting2(marchRslt, pt, renderer)
		//  CalculateLighting2(marchRslt, pt, renderer)

		marchRslt := RayMarchP(ray, renderer, false)
		pxColorVal := CalculateLightingTest(marchRslt, pt, renderer)
		// CalculateLightingTest(marchRslt, pt, renderer)
		renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
	}
}

func RayMarchWorkerLightingStatic(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()

	rangeSize := renderer.camera.Size() / int64(workers)

	for value := 0; value < int(rangeSize); value++ {
		if renderer.Reset.Load() {
			return
		}
		i := int64(value)*int64(workers) + int64(id)
		y := (i) % int64(renderer.camera.SizeY)
		x := i / int64(renderer.camera.SizeY)
		pt := Point{int(x), int(y)}
		ray := renderer.camera.RayForPixel(pt)
		marchRslt := RayMarch(ray, renderer, false)
		pxColorVal := CalculateLighting2(marchRslt, pt, renderer)
		//  CalculateLighting2(marchRslt, pt, renderer)

		// marchRslt := RayMarchP(ray, renderer, false)
		// pxColorVal := CalculateLightingTest(marchRslt, pt, renderer)
		// CalculateLightingTest(marchRslt, pt, renderer)
		renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
	}
}

func RenderOut(renderer *Renderer, workers int) {
	Render(renderer, workers)
	renderer.camera.FlushToDisk()
}

func Render(renderer *Renderer, workers int) {
	wg := new(sync.WaitGroup)
	pb := pb.NewOptions64(renderer.camera.Size(),
		pb.OptionSetDescription("Rendering Image..."),
		pb.OptionThrottle(65*time.Millisecond),
		pb.OptionShowIts(),
		pb.OptionSetItsString("px"),
		pb.OptionSpinnerType(14),
		pb.OptionFullWidth(),
		pb.OptionSetRenderBlankState(true),
		pb.OptionSetPredictTime(true),
		pb.OptionShowElapsedTimeOnFinish(),
		pb.OptionUseANSICodes(true),
	)

	for i := range workers {
		wg.Add(1)
		go RayMarchWorkerLighting3(i, workers, renderer, pb, wg)
	}

	log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")
	wg.Wait()
	// time.Sleep(time.Second)
	renderer.isDone.Store(true)
	renderer.Reset.Store(false)
}

// TODO: Rename render functions to be more clear
func (renderer *Renderer) Render2(workers int, wg *sync.WaitGroup) {
	renderer.isDone.Store(false)
	// startTime := time.Now()

	wg.Add(1)
	defer wg.Done()

	var wg2 sync.WaitGroup

	for i := range workers {
		wg2.Add(1)
		go RayMarchWorkerLighting6(i, workers, renderer, &wg2)
	}

	// log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")
	wg2.Wait()

	renderer.Reset.Store(false)
	renderer.isDone.Store(true)
	// renderDuration := time.Since(startTime)
	// fmt.Printf("Rendered frame in %s\n", renderDuration.String())

}

func (renderer *Renderer) RenderStatic(workers int, wg *sync.WaitGroup) {
	// renderer.isDone.Store(false)
	// startTime := time.Now()

	wg.Add(1)
	defer wg.Done()

	var wg2 sync.WaitGroup

	for i := range workers {
		wg2.Add(1)
		go RayMarchWorkerLightingStatic(i, workers, renderer, &wg2)
	}

	// log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")
	wg2.Wait()

	// renderer.Reset.Store(false)
	// renderer.isDone.Store(true)
	// renderDuration := time.Since(startTime)
	// fmt.Printf("Rendered frame in %s\n", renderDuration.String())

}

func NewDefaultRenderScene(opts RenderOpts) *Renderer {

	// Setup Scene
	scene := NewBlankScene()
	scene.AddDrawables(
		// drawables.NewNamedSphere("s2", vec3.Vec3{X: 0, Y: 0, Z: 0}, 1.5, color.RGBA{255, 255, 255, 255}, false, false),
		drawables.NewMandelB(60, 1.5, 8, vec3.Zero, color.RGBA{255, 255, 255, 255}, true),
		// drawables.NewMandelB("m2", 60, 1.5, 12, vec3.Zero, color.RGBA{25, 35, 45, 255}, false),
		// drawables.NewCube(vec3.Vec3{X: 0, Y: 0, Z: 0}, 10, color.RGBA{237, 66, 22, 255}),
		// drawables.NewNamedTorus("t1", vec3.Vec3{X: 10, Y: -4, Z: -2}, 4, 0.25, color.RGBA{130, 156, 154, 255}),
		//drawables.NewNamedCube("b1", vec3.Vec3{X: -4, Y: -2, Z: -1.5}, 1, color.RGBA{255, 255, 255, 255}),
	)

	scene.AddLights(
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: -15, Y: -1, Z: -1}, 1, color.RGBA{240, 240, 240, 255}, false),
		// drawables.NewNamedSphere("l2", vec3.Vec3{X: -2, Y: 2, Z: 0}, 1, color.RGBA{255, 255, 255, 255}, true, false),
		// drawables.NewLight(vec3.Vec3{X: 20, Y: -8, Z: -8}, 0.001, color.RGBA{100, 200, 200, 255}, false),
		// drawables.NewLight(vec3.Vec3{X: 1, Y: 20, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -8, Y: -20, Z: -8}, 0.001, color.RGBA{200, 19, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 3, Y: -7, Z: 8}, 0.001, color.RGBA{200, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 101, Y: -20, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -8, Y: -100, Z: -8}, 0.001, color.RGBA{70, 80, 90, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 220, Y: -8, Z: -8}, 0.001, color.RGBA{100, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 10, Y: 20, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -80, Y: -20, Z: -80}, 0.001, color.RGBA{200, 19, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 30, Y: -7, Z: 80}, 0.001, color.RGBA{200, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 1010, Y: -20, Z: 1}, 0.001, color.RGBA{199, 9, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -80, Y: -100, Z: -8}, 0.001, color.RGBA{70, 80, 90, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 20, Y: -8, Z: -8}, 0.001, color.RGBA{4, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 1, Y: -20, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -8, Y: -20, Z: -8}, 0.001, color.RGBA{200, 19, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 3, Y: -7, Z: 8}, 0.001, color.RGBA{200, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 101, Y: -20, Z: 10}, 0.001, color.RGBA{199, 3, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -8, Y: -100, Z: -8}, 0.001, color.RGBA{70, 80, 90, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 200, Y: -8, Z: -800}, 0.001, color.RGBA{100, 200, 1, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 10, Y: 120, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -80, Y: -20, Z: -80}, 0.001, color.RGBA{200, 19, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 30, Y: -7, Z: 80}, 0.001, color.RGBA{200, 200, 9, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 1010, Y: -20, Z: 1}, 0.001, color.RGBA{199, 55, 19, 255}, false),
		// drawables.NewLight(vec3.Vec3{X: -180, Y: -100, Z: -8}, 0.001, color.RGBA{70, 80, 90, 255}, false),
		// drawables.NewNamedSphere("l2", vec3.Vec3{X: -15, Y: 8, Z: 8}, 1, color.RGBA{0, 255, 0, 255}, false),
		// drawables.NewLight(vec3.Vec3{X: -5, Y: -2, Z: 1}, 0.005, color.RGBA{255, 255, 255, 255}, false),
		//drawables.NewNamedSphere("l3", vec3.Vec3{X: -10, Y: -10, Z: 10}, 0.5, color.RGBA{69, 79, 79, 255}),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: -1, Y: -1, Z: -15}, 1, color.RGBA{240, 240, 240, 255}, false),
		// drawables.NewNamedSphere("l5", vec3.Vec3{X: -8, Y: -8, Z: -15}, 1, color.RGBA{200, 200, 200, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: 1, Y: 1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: -1, Y: -1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: 1, Y: -1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: -1, Y: 1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: 0, Y: 0, Z: 17}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l5", vec3.Vec3{X: 8, Y: 8, Z: 15}, 1, color.RGBA{200, 200, 200, 255}, false),
	)

	// cam := NewCameraFOV(vec3.Vec3{X: 0, Y: 0, Z: 15}, opts.DimX, opts.DimY, opts.Fov, opts.OutPath)
	// cam := NewCameraFOV(vec3.Vec3{X: -15, Y: 0, Z: 0}, opts.DimX, opts.DimY, opts.Fov, opts.OutPath)
	cam := NewCameraFOV(vec3.Vec3{X: -15, Y: 0, Z: 0}, opts.DimX, opts.DimY, opts.Fov, opts.OutPath)

	// cam.up = vec3.UnitX()
	// cam.Dir = vec3.UnitZ().Mult(-1)

	cam.up = vec3.UnitZ
	cam.Dir = vec3.UnitX

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
