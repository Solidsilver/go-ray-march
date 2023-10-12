package main

import (
	"errors"
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/Solidsilver/go-ray-march/pkg/renderer"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()
	workersOpt := flag.Int("t", 4, "The number of concurrent jobs being processed")
	dimensionsOpt := flag.String("d", "1920x1080", "The dimensions of the image to render")
	fov := flag.Float64("fov", 20, "The field of view of the camera")
	outDir := flag.String("o", "./rend_out_0", "The directory to output the image to")
	flag.Parse()

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
	renderer.Render(r, rOps.Workers)
	r.GetCamera().FlushToDisk()
}
