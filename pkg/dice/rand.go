package dice

import (
	"math/rand"
)

// Rand is an interface used to get random values.
type Rand interface {
	Intn(n int) int
}

type stdRand struct{}

// StdRand is an empty struct which calls rand.Uint64n for its implementation.
var StdRand Rand = stdRand{}

// Get a random int from [0, n).
func (r stdRand) Intn(n int) int {
	return rand.Intn(n)
}
