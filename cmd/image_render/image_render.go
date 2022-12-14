package main

import (
	"flag"

	"github.com/Solidsilver/go-ray-march/pkg/renderer"
)

func main() {
	// defer profile.Start(profile.ProfilePath(".")).Stop()
	workersOpt := flag.Int("t", 4, "The number of concurrent jobs being processed")
	// outDir := flag.String("o", "./", "Folder to output the tiff files")
	// zoomLvl := flag.Int("z", 17, "Zoom level to create tiles at")
	// inDir := flag.String("i", "./", "Folder with the tif files")
	// outDir := flag.String("o", "./tiles", "Folder to output the tiles files")
	// verboseOpt := flag.Int("v", 1, "Set the verbosity level:\n"+
	// 	" 0 - Only prints error messages\n"+
	// 	" 1 - Adds run specs and error details\n"+
	// 	" 2 - Adds general progress info\n"+
	// 	" 3 - Adds debug info and details more detail\n")
	flag.Parse()

	// utils.SetupLogByLevel(*verboseOpt)
	// proc_runners.Bulk2Tiles(*inDir, *outDir, *workersOpt, *zoomLvl)
	// os.RemoveAll("../../rend_out")
	// os.Mkdir("../../rend_out", os.ModePerm)
	renderer.RenderDefault(*workersOpt)
}
