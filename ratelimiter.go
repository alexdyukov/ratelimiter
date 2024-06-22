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
	// panic by writes to closed channel impossible, because
	// channel closed by runtime.SetFinalizer which can be invoked only if `ratelimiter` link lost
	case ratelimiter.notify <- struct{}{}:
		return true
	}
}

// New returns an initialized RateLimiter with provided Bottleneck.
func New(bottleneck Bottleneck) *RateLimiter {
	if bottleneck == nil {
		panic("ratelimiter: bottleneck argument cannot be nil")
	}

	notify := make(chan struct{})

	ratelimiter := &RateLimiter{
		curRate: int64(0),
		maxRate: int64(bottleneck.MaxRate()),
		notify:  notify,
	}

	runtime.SetFinalizer(ratelimiter, func(_ *RateLimiter) { close(notify) })

	go func() {
		for {
			bottleneck.BreakThrough()

			if _, ok := <-notify; !ok {
				return
			}
		}
	}()

	return ratelimiter
}
