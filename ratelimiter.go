// Package ratelimiter provides an easy to use rate limiter with context.Context support,
// that can be used to limit the rate of arbitrary things mainly in your http middleware.
// See https://en.wikipedia.org/wiki/Rate_limiting .
package ratelimiter

import (
	"context"
	"runtime"
	"sync/atomic"
)

type (
	// Bottleneck represents gate keeper for rate limiter.
	// See bottleneck subpackage for examples.
	Bottleneck interface {
		BreakThrough()
		MaxRate() int
	}

	// RateLimiter implement rate limiter with wait option.
	RateLimiter struct {
		notify  chan struct{}
		curRate int64
		maxRate int64
	}
)

// Take returns true until rate and burst reached, false overwise.
// Requests over rate (rate < i <= rate+burst) holds in queue before next spot released.
// On context cancel or deadline expires request leaves the queue fast and returns false.
func (ratelimiter *RateLimiter) Take(ctx context.Context) bool {
	defer atomic.AddInt64(&ratelimiter.curRate, -1)

	if atomic.AddInt64(&ratelimiter.curRate, 1) > ratelimiter.maxRate {
		return false
	}

	select {
	case <-ctx.Done():
		return false
	case <-ratelimiter.notify:
		return true
	}
}

// New returns an initialized RateLimiter and shutdown function.
// Shutdown function completes RateLimiter's instance background tasks.
func New(bottleneck Bottleneck) (*RateLimiter, func()) {
	if bottleneck == nil {
		panic("ratelimiter: bottleneck argument cannot be nil")
	}

	needShutdown := atomic.Bool{}

	ratelimiter := &RateLimiter{
		curRate: int64(0),
		maxRate: int64(bottleneck.MaxRate()),
		notify:  make(chan struct{}),
	}

	runtime.SetFinalizer(ratelimiter, func(_ *RateLimiter) { needShutdown.Store(true) })

	go func() {
		for !needShutdown.Load() {
			bottleneck.BreakThrough()
			ratelimiter.notify <- struct{}{}
		}
		close(ratelimiter.notify)
	}()

	return ratelimiter, func() {}
}
