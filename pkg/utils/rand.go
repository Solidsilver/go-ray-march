package utils

func lgc(a, c, m, seed uint32) uint32 {
	return (a*seed + c) % m
}

func LGCRand(prev, max uint32) uint32 {
	const a = 1103515244
	const c = 12345
	return lgc(a, c, max, prev)
}

func LGCRandDec(prev, max uint32) LGCRandDecResult {
	next := LGCRand(prev, max)
	rndScaled := float64(next) / float64(max)
	return LGCRandDecResult{Rnd: rndScaled, Iter: next}
}

type LGCRandDecResult struct {
	Rnd  float64
	Iter uint32
}
