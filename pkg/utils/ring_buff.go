package utils

type RingBuffer[T any] struct {
	buffer []T
	Size   int
	idx    int
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
	rbf := new(RingBuffer[T])
	rbf.buffer = make([]T, size)
	rbf.Size = size
	rbf.idx = 0
	return rbf
}

func NewRingBufferFilled[T any](size int, filledWith T) *RingBuffer[T] {
	rbf := NewRingBuffer[T](size)
	for i := 0; i < size; i++ {
		rbf.buffer[i] = filledWith
	}
	return rbf
}

func (r *RingBuffer[T]) Push(val T) {
	r.buffer[r.idx] = val
	r.idx = (r.idx + 1) % r.Size
}

func (r *RingBuffer[T]) Get(idx int) T {
	getIdx := (r.idx - 1 + idx) % r.Size
	if getIdx < 0 {
		getIdx += r.Size
	}
	return r.buffer[getIdx]
}

func (r *RingBuffer[T]) Index() int {
	return r.idx
}
