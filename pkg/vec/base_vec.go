package vec

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Integer | constraints.Float
}

type Vec[T Number] struct {
	vals []T
}

func (v *Vec[T]) GetOrDefault(i int, val T) T {
	if i < len(v.vals) {
		return v.vals[i]
	}
	return val
}

func (v *Vec[T]) GetOrZero(i int) T {
	if i < len(v.vals) {
		return v.vals[i]
	}
	return 0
}

func (v *Vec[T]) Dim() int {
	return len(v.vals)
}

func New[T Number](vals ...T) *Vec[T] {
	vec := new(Vec[T])
	vec.vals = vals
	return vec
}

func Zero[T Number](n int) Vec[T] {
	return Vec[T]{make([]T, n)}
}

func One[T Number](n int) Vec[T] {
	vals := make([]T, n)
	for i := 0; i < n; i++ {
		vals[i] = 1
	}
	return Vec[T]{vals}
}

func OfSize[T Number](n int, val T) Vec[T] {
	vals := make([]T, n)
	for i := 0; i < n; i++ {
		vals[i] = val
	}
	return Vec[T]{vals}
}
