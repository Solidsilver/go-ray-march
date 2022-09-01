package renderer

import (
	"image"
	"image/color"
	"sync"

	"github.com/rs/zerolog/log"
	"solidsilver.dev/go-ray-marching/pkg/drawables"
	"solidsilver.dev/go-ray-marching/pkg/utils"
)

const MINIMUM_HIT_DISTANCE = 0.05
const MAXIMUM_TRACE_DISTANCE = 1000.0

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

func RayMarch(ray Ray, scene Scene) color.Color {
	totalDistTraveled := 0.0
	curPos := utils.NewCopy(ray.origin)
	totalMin := MAXIMUM_TRACE_DISTANCE
	var closest drawables.Drawable

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

		if minDist < 0 {
			return color.RGBA{0, 0, 0, 255}
		}

		if minDist < MINIMUM_HIT_DISTANCE {
			return closest.Color()
			// return color.RGBA{100, 200, 200, 255}
		}

		distP := minDist * 0.95
		totalDistTraveled += distP
		if minDist < totalMin {
			totalMin = minDist
		}
		curPos.Add(*curPos, *utils.NewCopy(ray.dir).Mult(distP))
	}
	// distPercet := MINIMUM_HIT_DISTANCE / totalMin
	// clr := 255 * distPercet
	// fmt.Println(clr)
	// return color.RGBA{100, 200, 200, uint8(clr)}
	return color.RGBA{0, 0, 0, 0}

}

func RayMarchWorker(rayJobs <-chan RayJob, scene Scene, image *image.RGBA, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range rayJobs {
		// log.Info().Msg("Worker got job")
		pxColorVal := RayMarch(job.ray, scene)
		// log.Info().Msg("Worker got job")

		image.Set(job.pos.X, job.pos.Y, pxColorVal)
	}
	log.Info().Msg("Worker: No more jobs, wg.Done()")
}

func Render(renderer Renderer, workers int) {
	cam := renderer.camera
	jobCount := cam.SizeX * cam.SizeY
	jobChan := make(chan RayJob, jobCount)
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go RayMarchWorker(jobChan, renderer.scene, cam.Image, wg)
	}

	center := Point{cam.SizeX / 2, cam.SizeY / 2}
	for i := 0; i <= cam.SizeX; i++ {
		for j := 0; j < cam.SizeX; j++ {
			camPosI := i - center.X
			camPosJ := j - center.Y
			// rayPos := cam.Pos.Add(cam.Pos, utils.Vec3{X: 0, Y: float64(camPosJ), Z: float64(camPosI)})
			rayPos := utils.Vec3{X: 0, Y: float64(camPosJ), Z: float64(camPosI)}
			rayPos.Add(rayPos, cam.Pos)
			// rayPos := utils.Vec3{X: 0, Y: float64(j), Z: float64(i)}
			// log.Info().Msgf("i: %v, j: %v, cpI: %v, cpJ: %v, rayPos: %v", i, j, camPosI, camPosJ, rayPos)
			ray := Ray{rayPos, cam.Dir}
			rayJob := RayJob{ray, Point{i, j}}
			jobChan <- rayJob
		}
	}
	log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")
	close(jobChan)

	wg.Wait()

	log.Info().Msg("Workers done, encoding image to path")

	cam.FlushToDisk()
}

func RenderDefault(workers int) {
	drawable1 := drawables.NewSphere(utils.Vec3{X: 350, Y: 0, Z: 0}, 300, color.RGBA{100, 200, 200, 255})
	drawable2 := drawables.NewSphere(utils.Vec3{X: 450, Y: -100, Z: 20}, 400, color.RGBA{252, 102, 11, 204})
	cam := NewCamera(utils.Vec3{X: 0, Y: 0, Z: 0}, 1920, 1080)
	renderer := Renderer{
		Scene{[]drawables.Drawable{drawable1, drawable2}},
		cam,
	}

	Render(renderer, workers)

	// for i := 0; i < 20; i++ {
	// 	cam.Pos = *cam.Pos.Add(cam.Pos, utils.Vec3UnitX())
	// 	Render(renderer, workers)
	// }

}
