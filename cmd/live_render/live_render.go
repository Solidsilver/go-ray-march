package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Solidsilver/go-ray-march/pkg/renderer"
	"github.com/fstanis/screenresolution"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pkg/profile"
)

type Game struct {
	offscreen        *ebiten.Image
	renderer         *renderer.Renderer
	windowHeight     int
	windowWidth      int
	renderWG         *sync.WaitGroup
	rops             renderer.RenderOpts
	keys             []ebiten.Key
	isProcessingMove atomic.Bool
}

func NewGame(opts renderer.RenderOpts) *Game {

	height, width := getWindowSize(opts)
	g := &Game{
		offscreen:    ebiten.NewImage(opts.DimX, opts.DimY),
		renderer:     renderer.NewDefaultRenderScene(opts),
		windowWidth:  width,
		windowHeight: height,
		renderWG:     &sync.WaitGroup{},
		rops:         opts,
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
	// if !gm.renderer.GetStatus() {
	// atomic.StoreUint32(&gm.isUpdating, 1)
	gm.offscreen.WritePixels(gm.renderer.GetCamera().Image.Pix)
	// atomic.StoreUint32(&gm.isUpdating, 0)
	// }
}

func isMvtKeyInList(keys []ebiten.Key) bool {
	for _, key := range keys {
		if key == ebiten.KeyArrowUp || key == ebiten.KeyArrowDown || key == ebiten.KeyArrowLeft || key == ebiten.KeyArrowRight || key == ebiten.KeyW || key == ebiten.KeyS || key == ebiten.KeyA || key == ebiten.KeyD {
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendJustPressedKeys(g.keys[:0])

	if !g.isProcessingMove.Load() /* && !g.renderer.Reset.Load() */ && isMvtKeyInList(g.keys) {
		g.isProcessingMove.Store(true)

		if !g.renderer.IsDone() {
			g.renderer.Reset.Store(true)
			g.renderWG.Wait()
		}
		g.renderer.GetCamera().Reset()

		moveAmt := 0.05

		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			moveAmt = 0.005
		}

		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.renderer.GetCamera().MoveUp(moveAmt)
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.renderer.GetCamera().MoveDown(moveAmt)
		}

		// if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		// 	g.renderer.GetCamera().RotateLeft(moveAmt)
		// }
		// if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		// 	g.renderer.GetCamera().RotateRight(moveAmt)
		// }
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.renderer.GetCamera().MoveLeft(moveAmt)
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.renderer.GetCamera().MoveRight(moveAmt)
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.renderer.GetCamera().MoveForward(moveAmt)
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.renderer.GetCamera().MoveBackward(moveAmt)
		}
		// g.renderer.Reset = false

		go g.renderer.Render2(g.rops.Workers, g.renderWG)
		g.isProcessingMove.Store(false)

	} else {
		g.updateOffscreen()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(g.offscreen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.renderer.GetCamera().SizeX, g.renderer.GetCamera().SizeY
}

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()
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
	go game.renderer.Render2(rOps.Workers, game.renderWG)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
