package utils_test

import (
	"testing"

	"solidsilver.dev/go-ray-marching/pkg/utils"
)

func TestLocalSig(t *testing.T) {

	steps := 5
	// dist := 300.432
	for i := 200.1; i < 220; i += 1.1 {
		sigL := utils.SigLocal(i/float64(steps), steps, 100)
		println(sigL)
		if steps == 5 && i > 210 {
			steps = 6
		}
	}

}
