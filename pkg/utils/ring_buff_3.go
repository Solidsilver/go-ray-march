package utils

type RingBuffer3 struct {
	buffer []float64
	idx    int
}

func NewRingBuffer3(fillWith float64) *RingBuffer3 {

	rbf := RingBuffer3{
		buffer: make([]float64, 3),
		idx:    0,
	}
	for i := 0; i < 3; i++ {
		rbf.buffer[i] = fillWith
	}

	return &rbf
}

func (r *RingBuffer3) Push(val float64) {
	r.buffer[r.idx] = val
	r.idx = (r.idx + 1) % 3
}

func (r RingBuffer3) Get(idx int) float64 {
	getIdx := (r.idx - 1 + idx) % 3
	if getIdx < 0 {
		getIdx += 3
	}
	return r.buffer[getIdx]
}

func (r RingBuffer3) Index() int {
	return r.idx
}

func (rbf RingBuffer3) GetCurrentSlope() float64 {
	// sum := 0.0
	// for i := 0; i > -2; i-- {
	// 	slope := rbf.Get(i) - rbf.Get(i-1)
	// 	sum += slope
	// }
	i1 := rbf.Get(-1)

	return ((rbf.Get(0) - i1) + (i1 - rbf.Get(-2))) / 2.0
}
