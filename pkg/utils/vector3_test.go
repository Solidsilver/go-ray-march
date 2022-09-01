package utils

import (
	"fmt"
	"testing"
)

func TestCrosProduct(t *testing.T) {
	v1 := Vec3{4, 5, -1}
	v2 := Vec3{-3, -20, 4.323}

	crossVal := new(Vec3).Cross(v1, v2)
	fmt.Println(crossVal)
}
