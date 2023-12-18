package vec

func (v Vec[T]) Add(v2 Vec[T]) Vec[T] {
	for i := 0; i < len(v.vals); i++ {
		v.vals[i] += v2.vals[i]
	}
	return v
}
