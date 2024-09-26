package engine

import "math/rand/v2"

var rng *rand.Rand

func init() {
	s := rand.NewPCG(0, 0)
	rng = rand.New(s)
}

func seedRNG(seed uint64) {
	s := rand.NewPCG(seed, seed)
	rng = rand.New(s)
}
