package bottleneck

import (
	"time"
)

// Equalizer typed Bottleneck similar to ticker queue.
// Provides bottle neck which pass 1 and no more then each {time.Second/RPS} second.
// RPS 1000 means pass 1 each 1 millisecond.
type Equalizer struct {
	lastCheckout int64
	diffDuration int64
	burst        int
}

// BreakThrough passes through bottle neck. Waits and pass if it busy.
func (bottleneck *Equalizer) BreakThrough() {
	newTime := bottleneck.lastCheckout + bottleneck.diffDuration
	now := time.Now().UnixNano()

	time.Sleep(time.Duration(newTime - now))

	bottleneck.lastCheckout = time.Now().UnixNano()
}

// MaxRate returns max rate for both simultaneously processed and in queue requests.
func (bottleneck *Equalizer) MaxRate() int {
	return int(int64(time.Second)/bottleneck.diffDuration) + bottleneck.burst
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
		lastCheckout: time.Now().Add(-1 * time.Second).UnixNano(),
		diffDuration: int64(time.Second) / int64(rps),
		burst:        burst,
	}
}
