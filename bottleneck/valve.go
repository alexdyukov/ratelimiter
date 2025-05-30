package bottleneck

import (
	"time"
)

// Valve typed Bottleneck similar to heart valve.
// Holds all together and pass them at one time.
type Valve struct {
	lastCheckout int64
	currentRate  int
	rps          int
	burst        int
}

// NewValve returns valve implementation of Bottleneck.
func NewValve(rps, burst int) (*Valve, error) {
	if rps <= 0 {
		return nil, ErrRPSNegativeOrZero
	}

	if burst < 0 {
		return nil, ErrBurstNegative
	}

	bottleneck := &Valve{
		lastCheckout: time.Now().UnixNano(),
		currentRate:  0,
		rps:          rps,
		burst:        burst,
	}

	return bottleneck, nil
}

// BreakThrough passes through bottle neck. Waits and pass if it busy.
func (bottleneck *Valve) BreakThrough() {
	if bottleneck.currentRate < bottleneck.rps {
		bottleneck.currentRate++

		return
	}

	now := time.Now().UnixNano()

	if now-bottleneck.lastCheckout < int64(time.Second) {
		time.Sleep(time.Duration(bottleneck.lastCheckout + int64(time.Second) - now))
	}

	bottleneck.currentRate = 1
	bottleneck.lastCheckout = time.Now().UnixNano()
}

// MaxRate returns max rate for both simultaneously processed and in queue requests.
func (bottleneck *Valve) MaxRate() int {
	return bottleneck.rps + bottleneck.burst
}
