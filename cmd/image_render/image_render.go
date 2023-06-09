package main

import (
	"flag"

	"github.com/Solidsilver/go-ray-march/pkg/renderer"
)

func main() {
	// defer profile.Start(profile.ProfilePath(".")).Stop()
	workersOpt := flag.Int("t", 4, "The number of concurrent jobs being processed")
	flag.Parse()

	renderer.RenderDefault(*workersOpt)
}
