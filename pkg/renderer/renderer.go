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

func calculatePixelColor(marchRslt MarchResult, screenPos Point, renderer *Renderer) color.RGBA {
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
		pxColorVal := calculatePixelColor(marchRslt, pt, renderer)
		renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		pb.Add(1)
	}
}

var (
	rmwRangeMap = make(map[int][]int)
	rngMtx      sync.RWMutex
)

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
		marchRslt := RayMarch(ray, renderer, false)
		pxColorVal := calculatePixelColor(marchRslt, pt, renderer)
		// calculatePixelColor(marchRslt, pt, renderer)

		// marchRslt := RayMarchP(ray, renderer, false)
		// pxColorVal := calculatePixelColor(marchRslt, pt, renderer)
		// // calculatePixelColor(marchRslt, pt, renderer)
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
		pxColorVal := calculatePixelColor(marchRslt, pt, renderer)
		//  calculatePixelColor(marchRslt, pt, renderer)

		// marchRslt := RayMarchP(ray, renderer, false)
		// pxColorVal := calculatePixelColor(marchRslt, pt, renderer)
		// calculatePixelColor(marchRslt, pt, renderer)
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
