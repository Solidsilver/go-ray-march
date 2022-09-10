package renderer

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"sync"

	"github.com/rs/zerolog/log"
	"solidsilver.dev/go-ray-marching/pkg/drawables"
	"solidsilver.dev/go-ray-marching/pkg/utils"
)

const MINIMUM_HIT_DISTANCE = 0.01
const MAXIMUM_TRACE_DISTANCE = 10000.0

type Ray struct {
	origin utils.Vec3
	dir    utils.Vec3
}

type RayJob struct {
	ray Ray
	pos Point
}

type Scene struct {
	Drawables []drawables.Drawable
}

type Renderer struct {
	scene  Scene
	camera *Camera
}

func RayMarch(ray *Ray, scene *Scene) color.RGBA {
	totalDistTraveled := 0.0
	curPos := utils.NewCopy(ray.origin)
	totalMin := MAXIMUM_TRACE_DISTANCE
	var closest drawables.Drawable
	steps := 0

	// if curPos.Equals(utils.Vec3{X: 0, Y: -5, Z: -29}) {
	// 	print("hello")
	// }

	for totalDistTraveled < MAXIMUM_TRACE_DISTANCE {

		minDist := MAXIMUM_TRACE_DISTANCE
		for _, obj := range scene.Drawables {
			dist := obj.Dist(*curPos)
			// if dist < MINIMUM_HIT_DISTANCE*2 {
			// 	print("hello")
			// }
			if dist < minDist {
				minDist = dist
				closest = obj
			}
		}
		distP := minDist * (1 - MINIMUM_HIT_DISTANCE)
		curPos.Add(*curPos, *utils.NewCopy(ray.dir).Mult(distP))
		steps++

		if minDist < 0 {
			return color.RGBA{0, 0, 0, 255}
		}

		if minDist < MINIMUM_HIT_DISTANCE {
			// _, distF := math.Modf(totalDistTraveled / float64(steps))
			col := closest.Color()
			noise := (rand.Float64() - 0.5) * 3
			stepsFl := float64(steps) + /*(distF)*/ +noise
			// print(stepsFl)
			darkP := (1 / math.Sqrt(stepsFl+16)) * 4
			r := uint8(darkP * float64(col.R))
			g := uint8(darkP * float64(col.G))
			b := uint8(darkP * float64(col.B))
			// newA := darkP * float64(a)

			return color.RGBA{r, g, b, col.A}
			// return color.RGBA{100, 200, 200, 255}
		}

		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist
		}

	}
	// distPercet := MINIMUM_HIT_DISTANCE / totalMin
	// clr := 255 * distPercet
	// fmt.Println(clr)
	// return color.RGBA{100, 200, 200, uint8(clr)}
	return color.RGBA{0, 0, 0, 0}

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
		// log.Info().Msg("Worker got job")
		// if job.X == 486 && job.Y == 302 {
		// 	print("hello")
		// }
		cam.RayForPixel2(job, ray)
		pxColorVal := RayMarch(ray, scene)
		// log.Info().Msg("Worker got job")

		image.Set(job.X, job.Y, pxColorVal)
	}
	// log.Info().Msg("Worker: No more jobs, wg.Done()")
}

func RayMarchWorker3(id int, workers int, scene *Scene, cam *Camera, image *image.RGBA, wg *sync.WaitGroup) {
	defer wg.Done()
	ray := new(Ray)
	for i := id; i <= cam.SizeX; i += workers {
		for j := 0; j <= cam.SizeY; j++ {
			pt := Point{i, j}
			cam.RayForPixel2(&pt, ray)
			pxColorVal := RayMarch(ray, scene)
			// log.Info().Msg("Worker got job")

			image.Set(pt.X, pt.Y, pxColorVal)
		}
	}
}

func Render(renderer Renderer, workers int) {
	cam := renderer.camera
	// jobCount := 10000
	// jobChan := make(chan RayJob, jobCount)
	// jobChan := make(chan *Point, jobCount)
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		// go RayMarchWorker2(jobChan, &renderer.scene, renderer.camera, cam.Image, wg)
		go RayMarchWorker3(i, workers, &renderer.scene, renderer.camera, cam.Image, wg)
	}

	// for i := 0; i <= cam.SizeX; i++ {
	// 	for j := 0; j < cam.SizeY; j++ {
	// 		pt := Point{i, j}
	// 		jobChan <- &pt
	// 	}
	// }
	log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")
	// close(jobChan)

	wg.Wait()

	log.Info().Msg("Workers done, encoding image to path")

	cam.FlushToDisk()
}

type WorkerPrepJob struct {
	i, j   int
	center Point
	camDir utils.Vec3
	camPos utils.Vec3
}

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
	drawable1 := drawables.NewSphere(utils.Vec3{X: 650, Y: 200, Z: 0}, 200, color.RGBA{100, 200, 200, 255})
	drawable2 := drawables.NewSphere(utils.Vec3{X: 850, Y: -100, Z: -200}, 300, color.RGBA{252, 102, 11, 255})
	drawable3 := drawables.NewSphere(utils.Vec3{X: 2000, Y: 0, Z: 0}, 400, color.RGBA{76, 96, 218, 255})
	drawable4 := drawables.NewSphere(utils.Vec3{X: 200, Y: 400, Z: 400}, 60, color.RGBA{1, 123, 6, 255})
	// cam := NewCamera(utils.Vec3{X: -1000, Y: 0, Z: 0}, 1920, 1080) // 1080p
	// cam := NewCamera(utils.Vec3{X: -100, Y: 0, Z: 0}, 3840, 2160) // 4k
	// cam := NewCamera(utils.Vec3{X: -100, Y: 0, Z: 0}, 7680, 4320) // 8k
	cam := NewCamera(utils.Vec3{X: -1000, Y: 0, Z: 0}, 15360, 8640) // 16k

	// cam := NewCamera(utils.Vec3{X: -1000, Y: 0, Z: 0}, 3600, 2400)

	renderer := Renderer{
		Scene{[]drawables.Drawable{drawable1, drawable2, drawable3, drawable4}},
		cam,
	}

	Render(renderer, workers)

	// for i := 0; i < 20; i++ {
	// 	cam.Pos = *cam.Pos.Add(cam.Pos, utils.Vec3UnitX())
	// 	Render(renderer, workers)
	// }

}
