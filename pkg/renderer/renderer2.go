package renderer

// import (
// 	"image"
// 	"image/color"
// 	"image/png"
// 	"os"
// 	"sync"

// 	"github.com/rs/zerolog/log"
// 	"solidsilver.dev/go-ray-marching/pkg/drawables"
// 	"solidsilver.dev/go-ray-marching/pkg/utils"
// )

// const MINIMUM_HIT_DISTANCE = 0.01
// const MAXIMUM_TRACE_DISTANCE = 1000.0

// type Ray struct {
// 	origin utils.Vec3
// 	dir    utils.Vec3
// }

// type RayJob struct {
// 	ray Ray
// 	pos Point
// }

// type RayResult struct {
// 	pos   Point
// 	color color.RGBA
// }

// type Scene struct {
// 	Drawables []drawables.Drawable
// }

// type Renderer struct {
// 	scene  Scene
// 	camera Camera
// }

// func RayMarch(ray Ray, scene Scene) color.RGBA {
// 	totalDistTraveled := 0.0
// 	curPos := new(utils.Vec3)

// 	for totalDistTraveled < MAXIMUM_TRACE_DISTANCE {
// 		curPos.Copy(ray.dir).Mult(totalDistTraveled)
// 		curPos.Add(*curPos, ray.origin)

// 		minDist := MAXIMUM_TRACE_DISTANCE
// 		for _, draw := range scene.Drawables {
// 			dist := draw.Dist(*curPos)
// 			if dist < minDist {
// 				// log.Info().Msgf("newMinDist: %v, ray: ", dist, ray)
// 				// log.Info().Msgf("Disst: %v", dist)
// 				minDist = dist
// 			}
// 		}

// 		if minDist < MINIMUM_HIT_DISTANCE {
// 			// log.Info().Msg("HIT!!")
// 			// return color.RGBA{1, 1, 1, 1}
// 			return color.RGBA{100, 200, 200, 0xff}
// 		}

// 		totalDistTraveled += minDist
// 	}
// 	return color.RGBA{0, 0, 0, 0}

// }

// func RayMarchWorker(rayJobs <-chan RayJob, scene Scene, rsltChan chan<- RayResult, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for job := range rayJobs {
// 		// log.Info().Msg("Worker got job")
// 		pxColorVal := RayMarch(job.ray, scene)
// 		// log.Info().Msg("Worker got job")
// 		rslt := RayResult{job.pos, pxColorVal}
// 		// log.Info().Msgf("Sending result %v", rslt)

// 		// image.Set(job.pos.X, job.pos.Y, pxColorVal)
// 		rsltChan <- rslt
// 	}
// 	log.Info().Msg("Worker: No more jobs, wg.Done()")
// }

// func RenderDefault(workers int) {
// 	drawable1 := drawables.Sphere{Center: utils.Vec3{X: 5, Y: 1, Z: 0}, Rad: 2.2}
// 	drawable2 := drawables.Sphere{Center: utils.Vec3{X: 50, Y: 0, Z: 0}, Rad: 20}
// 	cam := NewCamera(utils.Vec3{X: 0, Y: 0, Z: 0}, 500, 500)
// 	renderer := Renderer{
// 		Scene{[]drawables.Drawable{drawable1, drawable2}},
// 		cam,
// 	}

// 	jobCount := cam.SizeX * cam.SizeY
// 	jobChan := make(chan RayJob, jobCount)
// 	rsltChan := make(chan RayResult, jobCount+1)
// 	wg := new(sync.WaitGroup)

// 	for i := 0; i < workers; i++ {
// 		wg.Add(1)
// 		go RayMarchWorker(jobChan, renderer.scene, rsltChan, wg)
// 	}

// 	center := Point{cam.SizeX / 2, cam.SizeY / 2}
// 	for i := 0; i < cam.SizeX; i++ {
// 		for j := 0; j < cam.SizeX; j++ {
// 			camPosI := i - center.X
// 			camPosJ := j - center.Y
// 			// rayPos := cam.Pos.Add(cam.Pos, utils.Vec3{X: 0, Y: float64(camPosJ), Z: float64(camPosI)})
// 			rayPos := utils.Vec3{X: 0, Y: float64(camPosJ), Z: float64(camPosI)}
// 			// log.Info().Msgf("i: %v, j: %v, cpI: %v, cpJ: %v, rayPos: %v", i, j, camPosI, camPosJ, rayPos)
// 			ray := Ray{rayPos, cam.Dir}
// 			rayJob := RayJob{ray, Point{i, j}}
// 			jobChan <- rayJob
// 		}
// 	}

// 	close(jobChan)
// 	log.Info().Msg("Finished loading jobs, closing jobs & waiting for workers")

// 	wg.Wait()
// 	log.Info().Msg("Done wainting fo workers, closing results chan")
// 	close(rsltChan)

// 	log.Info().Msg("Drawing results to img")

// 	rsltImg := image.NewRGBA(image.Rect(0, 0, cam.SizeX, cam.SizeX))
// 	for rayRslt := range rsltChan {
// 		log.Info().Msgf("Setting color %v at %v, %v", rayRslt.color, rayRslt.pos.X, rayRslt.pos.Y)
// 		rsltImg.Set(rayRslt.pos.X, rayRslt.pos.Y, rayRslt.color)
// 	}

// 	log.Info().Msg("Workers done, encoding image to path")

// 	EncodePNGToPath("./render.png", rsltImg)
// }

// func EncodePNGToPath(imgPath string, img image.Image) error {
// 	out, err := os.Create(imgPath)
// 	if err != nil {
// 		log.Info().Msgf("Could not create output file: %v", imgPath)
// 		// log.Error().Msgf("Could not create output file: %v", imgPath)
// 		return err
// 	}
// 	defer out.Close()
// 	err = png.Encode(out, img)
// 	if err != nil {
// 		log.Info().Msg("Could not encode output image")
// 		// log.Error().Msg("Could not encode output image")
// 	}
// 	return err
// }
