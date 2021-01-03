package itree

import (
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

// seedOnce controls the random number generator seed.
var seedOnce sync.Once

// expectPanic calls f, expecting it to panic.
// If f does not panic, an error is logged to t.
func expectPanic(t *testing.T, why string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic: %s", why)
		}
	}()
	f()
}

// randomInterval returns a random float64-based interval in the range [0,2).
func randomInterval() Interval {
	start := rand.Float64()
	end := start + 1 - rand.Float64()
	return Interval{Float64(start), Float64(end)}
}

// seedRand seeds the random number generator with the current date.
func seedRand() {
	str := os.Getenv("TEST_SEED")
	seed, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		seed = time.Now().UnixNano()
	}
	println("Seed:", seed)
	rand.Seed(seed)
}
