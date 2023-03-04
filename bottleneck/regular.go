package bottleneck

import (
	"time"
)

// Regular typed Bottleneck similar to attractions queue.
// Provides uniform distribution of requests to stabilize CPU load.
type Regular struct {
	data  []time.Time
	pos   int
	burst int
}

// BreakThrough passes through bottle neck. Waits and pass if it busy.
func (bottleneck *Regular) BreakThrough() {
	nextAvailable := bottleneck.data[bottleneck.pos].Add(time.Second)
	time.Sleep(time.Until(nextAvailable))
	bottleneck.data[bottleneck.pos] = time.Now()
	bottleneck.pos = (bottleneck.pos + 1) % len(bottleneck.data)
}

// MaxRate returns max rate for both simultaneously processed and in queue requests.
func (bottleneck *Regular) MaxRate() int {
	return len(bottleneck.data) + bottleneck.burst
}

// NewRegular returns regular implementation of bottle neck.
func NewRegular(rps, burst int) *Regular {
	if rps <= 0 {
		panic("bottleneck: rps argument should be greater 0")
	}

	if burst < 0 {
		panic("bottleneck: burst argument should be greater or equal 0")
	}

	bottleneck := Regular{
		data:  make([]time.Time, rps),
		pos:   0,
		burst: burst,
	}

	someSecondsBefore := time.Now().Add(-1 * time.Second)

	for i := range bottleneck.data {
		bottleneck.data[i] = someSecondsBefore
	}

	return &bottleneck
}
