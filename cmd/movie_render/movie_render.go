package main

import (
	"errors"
	"flag"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Solidsilver/go-ray-march/pkg/renderer"
	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
	"github.com/schollz/progressbar/v3"
)

func main() {

	workersOpt := flag.Int("t", runtime.NumCPU(), "The number of concurrent jobs being processed")
	dimensionsOpt := flag.String("d", "1920x1080", "The dimensions of the image to render")
	fov := flag.Float64("fov", 20, "The field of view of the camera")
	outDir := flag.String("o", "./rend_out_0", "The directory to output the image to")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal()
		}
		pprof.StartCPUProfile(f)

		defer pprof.StopCPUProfile()
	}

	dims := strings.Split(*dimensionsOpt, "x")
	dimX := dims[0]
	dimY := dimX
	if len(dims) == 2 {
		dimY = dims[1]
	}

	dimXInt, err1 := strconv.Atoi(dimX)
	dimYInt, err2 := strconv.Atoi(dimY)
	if err1 != nil || err2 != nil {
		log.Fatal(errors.Join(err1, err2))
	}

	rOps := renderer.RenderOpts{
		Workers: *workersOpt,
		DimX:    dimXInt,
		DimY:    dimYInt,
		Fov:     *fov,
		OutPath: *outDir,
	}

	log.Println("Rendering with options: ", rOps.String())

	r := renderer.NewDefaultRenderScene(rOps)

	degrees := 360.0
	radius := 5.0
	degIncrements := 0.5

	pb := progressbar.NewOptions64(int64(degrees/degIncrements),
		progressbar.OptionSetDescription("Rendering Frames..."),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("frames"),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionUseANSICodes(true),
	)

	// rotate camera around center at a radius of 10, taking 360 frames
	for i := 0.0; i < degrees; i += degIncrements {

		// log.Println("Rendering frame: ", i)
		// startTime := time.Now()

		posX := radius * math.Cos(utils.DegToRad(i))
		posY := radius * math.Sin(utils.DegToRad(i))
		pos := vec3.New(posX, posY, 0)
		dirVec := vec3.DirFromPos(vec3.Zero, pos)

		// log.Printf("Setting camera pos: (%0.2f,%0.2f) | dir %v\n", posX, posY, dirVec)

		r.UpdateCamera(func(c *renderer.Camera) {
			c.Pos = pos
			c.Dir = dirVec
		})

		r.Render2(rOps.Workers, &sync.WaitGroup{})
		pb.Add(1)
		// log.Printf("Frame %f took: %s\n", i, time.Since(startTime).String())
		r.GetCamera().FlushToDisk()

	}

	log.Println("Done rendering")

}
