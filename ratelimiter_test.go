package ratelimiter_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter"
	"github.com/alexdyukov/ratelimiter/bottleneck"
	"github.com/stretchr/testify/assert"
)

const (
	overheadTestRPS   = int(time.Second / time.Duration(10))
	overheadTestBurst = 5
	testRPS           = 10
	testBurst         = 5
)

func TestContextCancel(t *testing.T) {
	t.Parallel()

	bn := bottleneck.NewValve(testRPS, testBurst)

	rateLimiter, shutdown := ratelimiter.New(bn)
	defer shutdown()

	successCount := atomic.Int32{}

	// overflow first to be sure of ctx.Done path of select statement
	var wg sync.WaitGroup
	for i := 0; i < testRPS; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			retval := rateLimiter.Take(context.Background())
			if retval {
				successCount.Add(1)
			}
		}()
	}
	wg.Wait()

	msgFormat := "RateLimiter.Take() should success rps's (%v) times"
	assert.Equal(t, int32(testRPS), successCount.Load(), msgFormat, testRPS)

	// overflow end

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	retval := rateLimiter.Take(ctx)

	assert.False(t, retval, "RateLimiter.Take() should return false on canceled contexts")
}

func TestNoOverflow(t *testing.T) {
	t.Parallel()

	bn := bottleneck.NewValve(testRPS, testBurst)

	rateLimiter, shutdown := ratelimiter.New(bn)
	defer shutdown()

	startTime := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < testRPS+testBurst; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			retval := rateLimiter.Take(context.Background())

			msgFormat := "RateLimiter.Take() should success %v's request out of total %v"
			assert.True(t, retval, msgFormat, i, testRPS+testBurst)
		}(i)
	}
	wg.Wait()

	spend := time.Since(startTime)
	msgFormat := "RateLimiter.Take() should block bursted request for more then 1sec but took %v"
	assert.Greater(t, spend, time.Second, msgFormat, spend)

	msgFormat = "RateLimiter.Take() should block bursted request no more then 2 sec but took %v"
	assert.Less(t, spend, 2*time.Second, msgFormat, spend)
}

func TestOverflow(t *testing.T) {
	t.Parallel()

	bn := bottleneck.NewValve(testRPS, testBurst)

	rateLimiter, shutdown := ratelimiter.New(bn)
	defer shutdown()

	successCount := atomic.Int32{}
	failCount := atomic.Int32{}

	var wg sync.WaitGroup
	for i := 0; i < 2*(testRPS+testBurst); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if rateLimiter.Take(context.Background()) {
				successCount.Add(1)
			} else {
				failCount.Add(1)
			}
		}()
	}
	wg.Wait()

	msgFormat := "RateLimiter.Take() should success at least rps's (%v) times"
	assert.LessOrEqual(t, int32(testRPS+testBurst), successCount.Load(), msgFormat, testRPS)

	msgFormat = "RateLimiter.Take() should fail at least once when we push over queue, but %v"
	assert.Less(t, int32(0), failCount.Load(), msgFormat, failCount.Load())
}
