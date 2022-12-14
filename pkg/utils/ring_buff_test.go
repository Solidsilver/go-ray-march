package utils

import (
	"fmt"
	"testing"
)

func TestGetIDX(t *testing.T) {

	rbf := NewRingBuffer[float64](3)

	rbf.Push(10)
	rbf.Push(9)
	rbf.Push(3)

	for i := -3; i < 5; i++ {
		fmt.Printf("rbf[%d]=%f\n", i, rbf.Get(i))
	}
}
