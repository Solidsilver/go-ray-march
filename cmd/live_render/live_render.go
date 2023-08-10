// Copyright 2017 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"flag"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/Solidsilver/go-ray-march/pkg/renderer"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	offscreen    *ebiten.Image
	renderer     *renderer.Renderer
	windowHeight int
	windowWidth  int
	isUpdating   uint32
}

func NewGame(opts renderer.RenderOpts) *Game {
	ratio := float64(opts.DimX) / float64(opts.DimY)
	fixedHeight := math.Min(1080, float64(opts.DimY))
	g := &Game{
		offscreen:    ebiten.NewImage(opts.DimX, opts.DimY),
		renderer:     renderer.DefaultScene(opts),
		windowWidth:  int(math.Round(fixedHeight * ratio)),
		windowHeight: int(fixedHeight),
	}
	return g
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

	game := NewGame(rOps)
	ebiten.SetWindowSize(game.windowWidth, game.windowHeight)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowTitle("Ray Marcher")
	go renderer.Render(game.renderer, rOps.Workers)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
