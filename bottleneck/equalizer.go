package bottleneck

import (
	"time"
)

// Equalizer typed Bottleneck similar to ticker queue.
// Provides bottle neck which pass 1 and no more then each {time.Second/RPS} second.
// RPS 1000 means pass 1 each 1 millisecond.
type Equalizer struct {
	lastCheckout time.Time
	diffDuration time.Duration
	burst        int
}

// BreakThrough passes through bottle neck. Waits and pass if it busy.
func (bottleneck *Equalizer) BreakThrough() {
	newTime := bottleneck.lastCheckout.Add(bottleneck.diffDuration)
	time.Sleep(time.Until(newTime))
	bottleneck.lastCheckout = time.Now()
}

// MaxRate returns max rate for both simultaneously processed and in queue requests.
func (bottleneck *Equalizer) MaxRate() int {
	return int(time.Second/bottleneck.diffDuration) + bottleneck.burst
}

// NewEqualizer returns equalized implementation of Bottleneck.
func NewEqualizer(rps, burst int) *Equalizer {
	if rps <= 0 {
		panic("bottleneck: rps argument should be greater 0")
	}

	if burst < 0 {
		panic("bottleneck: burst argument should be greater or equal 0")
	}

	return &Equalizer{
		lastCheckout: time.Now().Add(-1 * time.Second),
		diffDuration: time.Second / time.Duration(rps),
		burst:        burst,
	}
}
