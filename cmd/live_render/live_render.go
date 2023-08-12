package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/Solidsilver/go-ray-march/pkg/renderer"
	"github.com/fstanis/screenresolution"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	offscreen    *ebiten.Image
	renderer     *renderer.Renderer
	windowHeight int
	windowWidth  int
}

func NewGame(opts renderer.RenderOpts) *Game {

	height, width := getWindowSize(opts)
	g := &Game{
		offscreen:    ebiten.NewImage(opts.DimX, opts.DimY),
		renderer:     renderer.NewDefaultRenderScene(opts),
		windowWidth:  width,
		windowHeight: height,
	}
	return g
}

func getWindowSize(opts renderer.RenderOpts) (height, width int) {
	renderHeight := float64(opts.DimY)
	renderWidth := float64(opts.DimX)
	resolution := screenresolution.GetPrimary()
	if opts.DimX == resolution.Width && opts.DimY == resolution.Height {
		return resolution.Height, resolution.Width
	}
	if opts.DimX <= resolution.Width && opts.DimY <= resolution.Height {
		return opts.DimX, opts.DimY
	}
	renderRatio := renderWidth / renderHeight

	hDiff := renderHeight - float64(resolution.Height)
	wDiff := renderWidth - float64(resolution.Width)

	if hDiff > wDiff {
		h := float64(resolution.Height) * 0.90
		w := h * renderRatio
		height = int(math.Round(h))
		width = int(math.Round(w))
	} else {
		w := float64(resolution.Width)
		h := w / renderRatio
		height = int(math.Round(h))
		width = int(math.Round(w))
	}
	return height, width
}

// func color()

func (gm *Game) updateOffscreen() {
	if !gm.renderer.GetStatus() {
		// atomic.StoreUint32(&gm.isUpdating, 1)
		gm.offscreen.WritePixels(gm.renderer.GetCamera().Image.Pix)
		// atomic.StoreUint32(&gm.isUpdating, 0)
	}
}

func (g *Game) Update() error {
	// if atomic.LoadUint32(&g.isUpdating) == 0 {
	g.updateOffscreen()
	// }
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.offscreen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.renderer.GetCamera().SizeX, g.renderer.GetCamera().SizeY
}

func main() {
	resolution := screenresolution.GetPrimary()
	defaultRes := fmt.Sprintf("%dx%d", resolution.Width, resolution.Height)
	workersOpt := flag.Int("t", 4, "The number of concurrent jobs being processed")
	dimensionsOpt := flag.String("d", defaultRes, "The dimensions of the image to render")
	fov := flag.Float64("fov", 20, "The field of view of the camera")
	outDir := flag.String("o", "./rend_out_0", "The directory to output the image to")
	scaling := flag.Int("s", 1, "Scale to render at")
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

	if *scaling != 1 {
		dimXInt = dimXInt * *scaling
		dimYInt = dimYInt * *scaling
	}

	rOps := renderer.RenderOpts{
		Workers: *workersOpt,
		DimX:    dimXInt,
		DimY:    dimYInt,
		Fov:     *fov,
		OutPath: *outDir,
	}

	log.Println("Rendering with options: ", rOps.String())

	game := NewGame(rOps)
	ebiten.SetWindowSize(game.windowWidth, game.windowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeOnlyFullscreenEnabled)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowTitle("Ray Marcher")
	go renderer.Render(game.renderer, rOps.Workers)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
