package bottleneck

import (
	"time"
)

// Valve typed Bottleneck similar to heart valve.
// Holds all together and pass them at one time.
type Valve struct {
	lastCheckout time.Time
	currentRate  int
	rps          int
	burst        int
}

// BreakThrough passes through bottle neck. Waits and pass if it busy.
func (bottleneck *Valve) BreakThrough() {
	if bottleneck.currentRate < bottleneck.rps {
		bottleneck.currentRate++

		return
	}

	if time.Since(bottleneck.lastCheckout) < time.Second {
		time.Sleep(time.Second)
	}

	bottleneck.currentRate = 1
	bottleneck.lastCheckout = time.Now()
}

// MaxRate returns max rate for both simultaneously processed and in queue requests.
func (bottleneck *Valve) MaxRate() int {
	return bottleneck.rps + bottleneck.burst
}

// NewValve returns valve implementation of Bottleneck.
func NewValve(rps, burst int) *Valve {
	if rps <= 0 {
		panic("bottleneck: rps argument should be greater 0")
	}

	if burst < 0 {
		panic("bottleneck: burst argument should be greater or equal 0")
	}

	return &Valve{
		lastCheckout: time.Now(),
		currentRate:  0,
		rps:          rps,
		burst:        burst,
	}
}
