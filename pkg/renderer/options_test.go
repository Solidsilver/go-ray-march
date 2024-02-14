package renderer

import (
	"testing"
)

func TestOptionsPrinting(t *testing.T) {
	defOpt := DefaultLightingOpts()
	t.Logf("Default options: %s", defOpt.JsonString())
}
