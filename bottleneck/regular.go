package bottleneck

import (
	"time"
)

// Regular typed Bottleneck similar to attractions queue.
// Provides uniform distribution of requests to stabilize CPU load.
type Regular struct {
	data  []int64
	pos   int
	burst int
}

// NewRegular returns regular implementation of bottle neck.
func NewRegular(rps, burst int) (*Regular, error) {
	if rps <= 0 {
		return nil, ErrRPSNegativeOrZero
	}

	if burst < 0 {
		return nil, ErrBurstNegative
	}

	bottleneck := Regular{
		data:  make([]int64, rps),
		pos:   0,
		burst: burst,
	}

	someTimeBefore := time.Now().Add(-1 * time.Second).UnixNano()

	for i := range bottleneck.data {
		bottleneck.data[i] = someTimeBefore
	}

	return &bottleneck, nil
}

// BreakThrough passes through bottle neck. Waits and pass if it busy.
func (bottleneck *Regular) BreakThrough() {
	nextAvailable := bottleneck.data[bottleneck.pos] + time.Second.Nanoseconds()
	now := time.Now().UnixNano()

	time.Sleep(time.Duration(nextAvailable - now))

	bottleneck.data[bottleneck.pos] = time.Now().UnixNano()
	bottleneck.pos = (bottleneck.pos + 1) % len(bottleneck.data)
}

// MaxRate returns max rate for both simultaneously processed and in queue requests.
func (bottleneck *Regular) MaxRate() int {
	return len(bottleneck.data) + bottleneck.burst
}
