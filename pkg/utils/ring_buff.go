package utils

type RingBuffer[T any] struct {
	buffer []T
	size   int
	idx    int
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
	rbf := new(RingBuffer[T])
	rbf.buffer = make([]T, size)
	rbf.size = size
	rbf.idx = 0
	return rbf
}

func (r *RingBuffer[T]) Push(val T) {
	r.buffer[r.idx] = val
	r.idx = (r.idx + 1) % r.size
}

func (r RingBuffer[T]) Get(idx int) T {
	getIdx := (r.idx - 1 + idx) % r.size
	if getIdx < 0 {
		getIdx += r.size
	}
	return r.buffer[getIdx]
}

func (r RingBuffer[T]) Index() int {
	return r.idx
}

func (r RingBuffer[T]) Size() int {
	return r.size
}
