package renderer

import (
	"image/color"
	"sync"
	"time"

	"github.com/Solidsilver/go-ray-marching/pkg/drawables"
	"github.com/Solidsilver/go-ray-marching/pkg/vec3"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

const MINIMUM_HIT_DISTANCE = 0.001
const MAXIMUM_TRACE_DISTANCE = 10000.0
const MAX_STEPS = 10000

// var BG_COLOR = color.RGBA{255, 255, 255, 255}
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
			// if pt.X == 1920/2 && pt.Y == 1080/2 {
			// 	print("0,0")
			// }
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
			// if pt.X == 1920/2 && pt.Y == 1080/2 {
			// 	print("0,0")
			// }
			ray := renderer.camera.RayForPixel(pt)
			marchRslt := RayMarch(ray, renderer.scene)
			pxColorVal := CalculateLighting(marchRslt, renderer)
			renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		}
	}
}

func RayMarchWorkerLighting3(id int, workers int, renderer *Renderer, pb *progressbar.ProgressBar, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i <= renderer.camera.SizeX; i++ {
		for j := id; j <= renderer.camera.SizeY; j += workers {
			j2 := (j + i) % renderer.camera.SizeY
			pt := Point{i, j2}
			// if pt.X == 1920/2 && pt.Y == 1080/2 {
			// 	print("0,0")
			// }
			ray := renderer.camera.RayForPixel(pt)
			marchRslt := RayMarch(ray, renderer.scene)
			pxColorVal := CalculateLighting(marchRslt, renderer)
			renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
			pb.Add(1)
		}
		// pb.Add(renderer.camera.SizeY)
	}
}

func Render(renderer *Renderer, workers int) {
	wg := new(sync.WaitGroup)
	// pb := progressbar.Default(renderer.camera.Size(), "Rendering Image...")
	pb := progressbar.NewOptions64(renderer.camera.Size(),
		progressbar.OptionSetDescription("Rendering Image..."),
		progressbar.OptionThrottle(65*time.Millisecond),
		// progressbar.OptionShowCount(),
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
	log.Info().Msg("Workers done, encoding image to path")
	renderer.camera.FlushToDisk()
}

func RenderDefault(workers int) {
	// drawable1 := drawables.NewNamedSphere("s1", vec3.Vec3{X: 0, Y: 0, Z: 0}, 2.5, color.RGBA{100, 200, 200, 255})
	// drawable2 := drawables.NewNamedSphere("s2", vec3.Vec3{X: 20, Y: 4, Z: 4}, 3.5, color.RGBA{252, 102, 11, 255})
	// drawable3 := drawables.NewNamedSphere("s3", vec3.Vec3{X: 40, Y: 0, Z: 0}, 4.5, color.RGBA{76, 96, 218, 255})
	// drawable4 := drawables.NewNamedSphere("d4", vec3.Vec3{X: 2, Y: 4, Z: 4}, 1, color.RGBA{1, 123, 6, 255})
	// drawable4 := drawables.NewNamedCube("b1", vec3.Vec3{X: 10, Y: -4, Z: -4}, 2, color.RGBA{255, 255, 255, 255})
	// drawable1 := drawables.NewMandelB("m1", 60, 1.5, 8, vec3.Zero(), color.RGBA{100, 200, 200, 255})
	drawable1 := drawables.NewMandelB("m1", 60, 1.5, 8, vec3.Zero(), color.RGBA{16, 51, 74, 255})
	// cam := NewCameraFOV(vec3.Vec3{X: -50, Y: 0, Z: 0}, 1920, 1080, 5) // 1080p
	// cam := NewCameraFOV(vec3.Vec3{X: -50, Y: 0, Z: 0}, 3840, 2160, 3) // 4k
	// cam := NewCameraFOV(vec3.Vec3{X: -10, Y: 0, Z: -1}, 7680, 4320, 10) // 8k
	// cam := NewCameraFOV(vec3.Vec3{X: -20, Y: 0, Z: 0}, 15360, 8640, 45) // 16k
	// cam := NewCameraFOV(vec3.Vec3{X: -25, Y: 0, Z: 0}, 30720, 17280, 45) // 32k
	light1 := drawables.NewNamedSphere("l1", vec3.Vec3{X: -10, Y: -10, Z: -10}, 1, color.RGBA{255, 255, 255, 255})
	light2 := drawables.NewNamedSphere("l2", vec3.Vec3{X: -10, Y: 10, Z: 0}, 1, color.RGBA{255, 255, 255, 255})
	light3 := drawables.NewNamedSphere("l3", vec3.Vec3{X: -10, Y: -10, Z: 10}, 1, color.RGBA{255, 255, 255, 255})
	// light4 := drawables.NewNamedSphere("l4", vec3.Vec3{X: -100, Y: 0, Z: 0}, 1, color.RGBA{255, 255, 255, 100})
	scene := NewScene([]drawables.Drawable{drawable1}, []drawables.Drawable{light1, light2, light3})

	cam := NewCameraFOV(vec3.Vec3{X: -50, Y: 0, Z: 0}, 2000, 2000, 2.75) // 4k
	// cam := NewCameraFOV(vec3.Vec3{X: -30, Y: 0, Z: 0}, 2000, 2000, 10) // 4k
	// scene := NewScene([]drawables.Drawable{drawable1, drawable2, drawable3, drawable4}, []drawables.Drawable{light1, light2})

	renderer := Renderer{
		scene,
		cam,
	}

	Render(&renderer, workers)

	// iter := 100.0
	//
	// for i := 0.0; i < iter; i++ {
	//
	// 	pow := 4 * (math.Cos(2*math.Pi*i/iter) + 1)
	//
	// 	drawable1 = drawables.NewMandelB("m1", 60, 1.5, pow, vec3.Zero(), color.RGBA{255, 255, 255, 255})
	//
	// 	scene = NewScene([]drawables.Drawable{drawable1}, []drawables.Drawable{light1, light2, light3})
	//
	// 	renderer = Renderer{
	// 		scene,
	// 		cam,
	// 	}
	// 	// 	// cam.Dir.Add(cam.Dir, right)
	// 	// 	// cam.Pos.Sub(cam.Pos, right)
	// 	// 	rad := 50.0
	// 	// 	val := 2 * math.Pi * i / iter
	// 	// 	cam.Pos.X = math.Cos(val) * rad
	// 	// 	cam.Pos.Y = math.Sin(val) * rad
	// 	// 	cam.Dir = (vec3.DirFromPos(vec3.Vec3{X: 0, Y: 0, Z: 0}, cam.Pos))
	// 	// 	// fmt.Printf("Pos: %v\n", cam.Pos)
	// 	// 	// fmt.Printf("Dir: %v\n", cam.Dir)
	// 	// 	// drawable4 = drawables.NewNamedCube("b1", *utils.NewAdd(drawable4.Pos(), right), 1, color.RGBA{1, 123, 6, 255})
	// 	// 	// scene = NewScene([]drawables.Drawable{drawable1, drawable2, drawable3, drawable4}, []drawables.Drawable{light1, light2, light3})
	// 	// 	// renderer.scene = scene
	// 	Render(&renderer, workers)
	// }

}
