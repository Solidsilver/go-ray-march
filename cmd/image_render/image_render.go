package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"

	"github.com/Solidsilver/go-ray-march/pkg/renderer"
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
	// renderer.Render3(r, rOps.Workers)
	// startTime := time.Now()
	r.Render2(rOps.Workers, &sync.WaitGroup{})
	// renderDuration := time.Since(startTime)
	// log.Println("Rendered in: ", renderDuration.String())
	r.GetCamera().FlushToDisk()
	// log.Println("Flushed to disk in: ", time.Since(startTime).String())
}
