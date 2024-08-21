package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	_ "net/http/pprof"

	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/renderer"
	"github.com/fstanis/screenresolution"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pkg/profile"
	"golang.org/x/exp/slices"
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

var validKeys = []ebiten.Key{
	ebiten.KeyArrowUp,
	ebiten.KeyArrowDown,
	ebiten.KeyArrowLeft,
	ebiten.KeyArrowRight,
	ebiten.KeyW,
	ebiten.KeyS,
	ebiten.KeyA,
	ebiten.KeyD,
}

func isMvtKeyInList(keys []ebiten.Key) bool {
	for _, key := range keys {
		if slices.Contains(validKeys, key) {
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendJustPressedKeys(g.keys[:0])
	// if ebiten.IsKeyPressed(ebiten.KeyP) {
	// 	go g.renderer.GetCamera().FlushToDisk()
	// }
	if slices.Contains(g.keys, ebiten.KeyP) {
		go g.renderer.GetCamera().FlushToDisk()
	}

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

		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			scene := g.renderer.GetScene()
			bulb := scene.Drawables[0].(drawables.MandelBulb)
			bulb.Power += 0.1
			fmt.Printf("Power: %f\n", bulb.Power)
			scene.Drawables[0] = bulb
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			scene := g.renderer.GetScene()
			bulb := scene.Drawables[0].(drawables.MandelBulb)
			bulb.Power -= 0.1
			fmt.Printf("Power: %f\n", bulb.Power)
			scene.Drawables[0] = bulb
		}
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
	defer profile.Start(profile.MemProfile).Stop()
	resolution := screenresolution.GetPrimary()
	defaultRes := fmt.Sprintf("%dx%d", resolution.Width, resolution.Height)
	workersOpt := flag.Int("t", runtime.NumCPU(), "The number of concurrent jobs being processed")
	dimensionsOpt := flag.String("d", defaultRes, "The dimensions of the image to render")
	fov := flag.Float64("fov", 20, "The field of view of the camera")
	outDir := flag.String("o", "./rend_out_0", "The directory to output the image to")
	scaling := flag.Int("s", 1, "Scale to render at")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		fmt.Println("Writing CPU profile to: ", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal()
		}
		pprof.StartCPUProfile(f)

		defer pprof.StopCPUProfile()
	}
	go func() {
		http.ListenAndServe(":8080", nil)
	}()

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
