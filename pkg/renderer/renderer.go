package renderer

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/rs/zerolog/log"
	"solidsilver.dev/go-ray-marching/pkg/drawables"
	"solidsilver.dev/go-ray-marching/pkg/utils"
)

const MINIMUM_HIT_DISTANCE = 0.05
const MAXIMUM_TRACE_DISTANCE = 10000.0

// const COL_TRANSP = color.RGBA{0, 0, 0, 0}

type Ray struct {
	origin utils.Vec3
	dir    utils.Vec3
}

type Renderer struct {
	scene  *Scene
	camera *Camera
}

// func RayMarchWorker(rayJobs <-chan RayJob, scene Scene, image *image.RGBA, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for job := range rayJobs {
// 		// log.Info().Msg("Worker got job")
// 		pxColorVal := RayMarch(job.ray, scene)
// 		// log.Info().Msg("Worker got job")

// 		image.Set(job.pos.X, job.pos.Y, pxColorVal)
// 	}
// 	log.Info().Msg("Worker: No more jobs, wg.Done()")
// }

func RayMarchWorker2(rayJobs <-chan *Point, scene *Scene, cam *Camera, image *image.RGBA, wg *sync.WaitGroup) {
	defer wg.Done()
	ray := new(Ray)
	for job := range rayJobs {
		cam.RayForPixel2(job, ray)
		// pxColorVal := RayMarch(ray, scene)

		// image.Set(job.X, job.Y, pxColorVal)
	}
}

func RayMarchWorker3(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()
	ray := new(Ray)
	iter := uint32(id)
	for i := id; i <= renderer.camera.SizeX; i += workers {
		for j := 0; j <= renderer.camera.SizeY; j++ {
			pt := Point{i, j}
			renderer.camera.RayForPixel2(&pt, ray)
			rslt := utils.LGCRandDec(iter, 100000)
			iter = rslt.Iter
			pxColorVal := RayMarch(ray, renderer.scene, rslt.Rnd)
			// log.Info().Msg("Worker got job")

			renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		}
	}
}

func RayMarchWorker4(id int, workers int, renderer *Renderer, wg *sync.WaitGroup) {
	defer wg.Done()
	ray := new(Ray)
	for i := id; i <= renderer.camera.SizeX; i += workers {
		for j := 0; j <= renderer.camera.SizeY; j++ {
			pt := Point{i, j}
			renderer.camera.RayForPixel2(&pt, ray)
			marchRslt := RayMarch2(ray, renderer.scene)
			lightHits := make([]float64, 0)
			pxColorVal := color.RGBA{0, 0, 0, 255}
			if marchRslt.HitObject != nil {
				// if marchRslt.HitObject.ID() == "d2" {
				// 	print("hitd2")
				// }
				// pxColorVal = color.RGBA{255, 255, 255, 255}
				hitPoint := marchRslt.HitPos
				pxColorVal = marchRslt.HitObject.Color()
				for _, lSource := range renderer.scene.Lights {
					dir := new(utils.Vec3)
					dir.Sub(lSource.Pos(), hitPoint)
					dir.Unit()
					surfaceNormal := SurfaceNormal(hitPoint, marchRslt.HitObject)
					// bounceDeg := 180 - utils.Angle(*dir, utils.DirFromPos(hitPoint, renderer.camera.Dir))
					bounceDeg := math.Min(90.0, utils.Angle(*dir, surfaceNormal))
					ray := Ray{hitPoint, *dir}
					rslt := RayMarch2(&ray, renderer.scene)
					if drawables.Equals(rslt.HitObject, lSource) {
						// if bounceDeg > 40 {
						// 	print("hello")
						// }
						brightness := float64(rslt.HitObject.Color().A)
						brightness = brightness * (90 - bounceDeg) / 90
						lightHits = append(lightHits, brightness)
						// pxColorVal = color.RGBA{
						// 	R: min(255, pxColorVal.R *
						// }
					}
					// else {
					// 	brightness := 20.0
					// 	lightHits = append(lightHits, brightness)
					// }

					//march ray to light
				}

				sumBright := 0.0
				for _, lHit := range lightHits {
					sumBright += lHit
				}
				avgBright := sumBright // float64(len(lightHits))

				pxColorVal = color.RGBA{
					R: uint8(math.Min(255, float64(pxColorVal.R)*(avgBright/255))),
					G: uint8(math.Min(255, float64(pxColorVal.G)*(avgBright/255))),
					B: uint8(math.Min(255, float64(pxColorVal.B)*(avgBright/255))),
					A: pxColorVal.A,
				}

			}

			renderer.camera.Image.Set(pt.X, pt.Y, pxColorVal)
		}
	}
}

func Render(renderer *Renderer, workers int) {
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		// go RayMarchWorker2(jobChan, &renderer.scene, renderer.camera, cam.Image, wg)
		go RayMarchWorker4(i, workers, renderer, wg)
	}

	log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")

	wg.Wait()

	log.Info().Msg("Workers done, encoding image to path")

	renderer.camera.FlushToDisk()
}

// type WorkerPrepJob struct {
// 	i, j   int
// 	center Point
// 	camDir utils.Vec3
// 	camPos utils.Vec3
// }

// func QueueJobs(rayJobs chan<- RayJob, center Point, cam Camera, workers int) {
// 	wg := new(sync.WaitGroup)
// 	workerChan := make(chan WorkerPrepJob, 500)
// 	for i := 0; i < workers; i++ {
// 		wg.Add(1)
// 		go addJobWorker(workerChan, rayJobs, wg)
// 	}

// 	for i := 0; i <= cam.SizeX; i++ {
// 		for j := 0; j < cam.SizeX; j++ {
// 			workerChan <- WorkerPrepJob{i, j, center, cam.Dir, cam.Pos}
// 			// wg.Add(1)
// 			// go addJobWorker(rayJobs, center, cam, i, j, wg)
// 		}
// 	}
// 	close(workerChan)
// 	wg.Wait()
// 	close(rayJobs)
// }

// func addJobWorker(jobChan <-chan WorkerPrepJob, rayJobChan chan<- RayJob, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for job := range jobChan {

// 		camPosI := job.i - job.center.X
// 		camPosJ := job.j - job.center.Y
// 		// rayPos := cam.Pos.Add(cam.Pos, utils.Vec3{X: 0, Y: float64(camPosJ), Z: float64(camPosI)})
// 		rayPos := utils.Vec3{X: 0, Y: float64(camPosJ), Z: float64(camPosI)}
// 		rayPos.Add(rayPos, job.camPos)
// 		// rayPos := utils.Vec3{X: 0, Y: float64(j), Z: float64(i)}
// 		// log.Info().Msgf("i: %v, j: %v, cpI: %v, cpJ: %v, rayPos: %v", i, j, camPosI, camPosJ, rayPos)
// 		ray := Ray{rayPos, job.camDir}
// 		rayJob := RayJob{ray, Point{job.i, job.j}}
// 		rayJobChan <- rayJob
// 	}
// }

func RenderDefault(workers int) {
	drawable1 := drawables.NewNamedSphere("d1", utils.Vec3{X: 6, Y: 2, Z: 0}, 2.5, color.RGBA{100, 200, 200, 255})
	drawable2 := drawables.NewNamedSphere("d2", utils.Vec3{X: 8, Y: -1, Z: -2}, 3.5, color.RGBA{252, 102, 11, 255})
	drawable3 := drawables.NewNamedSphere("d3", utils.Vec3{X: 40, Y: 0, Z: 0}, 4.5, color.RGBA{76, 96, 218, 255})
	drawable4 := drawables.NewNamedSphere("d4", utils.Vec3{X: 2, Y: 4, Z: 4}, 0.9, color.RGBA{1, 123, 6, 255})
	cam := NewCamera(utils.Vec3{X: -25, Y: 0, Z: 0}, 1920, 1080) // 1080p
	// cam := NewCamera(utils.Vec3{X: -25, Y: 0, Z: 0}, 3840, 2160) // 4k
	// cam := NewCamera(utils.Vec3{X: -10, Y: 0, Z: 0}, 7680, 4320) // 8k
	// cam := NewCamera(utils.Vec3{X: -10, Y: 0, Z: 0}, 15360, 8640) // 16k
	// cam := NewCamera(utils.Vec3{X: -1000, Y: 0, Z: 0}, 30720, 17280) // 32k

	light1 := drawables.NewNamedSphere("l1", utils.Vec3{X: -30, Y: 0, Z: 0}, 1, color.RGBA{0, 0, 0, 255})
	light2 := drawables.NewNamedSphere("l2", utils.Vec3{X: 30, Y: 30, Z: 30}, 1, color.RGBA{0, 0, 0, 255})

	// cam := NewCamera(utils.Vec3{X: -1000, Y: 0, Z: 0}, 3600, 2400)
	scene := NewScene([]drawables.Drawable{drawable1, drawable2, drawable3, drawable4}, []drawables.Drawable{light1, light2})

	renderer := Renderer{
		scene,
		cam,
	}

	Render(&renderer, workers)
	// right := utils.Vec3UnitY()
	// right.Mult(0.5)
	// back := utils.Vec3UnitX()
	// back.Mult(0.5)

	// for i := 0; i < 50; i++ {
	// 	cam.Pos.Add(cam.Pos, right)
	// 	Render(&renderer, workers)
	// }

	// for i := 0; i < 50; i++ {
	// 	cam.Pos.Sub(cam.Pos, back)
	// 	Render(&renderer, workers)
	// }

	// for i := 0; i < 50; i++ {
	// 	cam.Pos.Sub(cam.Pos, right)
	// 	Render(&renderer, workers)
	// }

	// for i := 0; i < 50; i++ {
	// 	cam.Pos.Add(cam.Pos, back)
	// 	Render(&renderer, workers)
	// }

}
