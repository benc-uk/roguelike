package engine

import (
	"math/rand/v2"
)

// ============================================================================
// Random number generator
// This is a seedable RNG for deterministic levels and level generation
// ============================================================================

var rng *rand.Rand

func init() {
	s := rand.NewPCG(0, 0)
	rng = rand.New(s)
}

func seedRNG(seed uint64) {
	s := rand.NewPCG(seed, seed)
	rng = rand.New(s)
}

func randString(strings ...string) string {
	return strings[rng.IntN(len(strings))]
}
