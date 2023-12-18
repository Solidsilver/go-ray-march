package renderer_test

import (
	"testing"

	"github.com/Solidsilver/go-ray-march/pkg/renderer"
)

func TestOptionsPrinting(t *testing.T) {
	defOpt := renderer.DefaultLightingOpts()
	t.Logf("Default options: %s", defOpt.JsonString())
}
