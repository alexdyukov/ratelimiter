package ratelimiter_test

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/v2"
	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
	"go.uber.org/goleak"
)

const (
	overheadTestMinDuration = time.Duration(10)
	overheadTestRPS         = int(time.Second / overheadTestMinDuration)
	overheadTestBurst       = 5
	testRPS                 = 10
	testBurst               = 5
)

func TestContextCancel(t *testing.T) {
	defer detectLeak(t)()

	rateLimiter := ratelimiter.New(bottleneck.NewValve(testRPS, testBurst))

	successCount := atomic.Int32{}

	// overflow first to be sure of ctx.Done path of select statement
	var waitGroup sync.WaitGroup
	for i := 0; i < testRPS; i++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			retval := rateLimiter.Take(context.Background())
			if retval {
				successCount.Add(1)
			}
		}()
	}
	waitGroup.Wait()

	if actual := successCount.Load(); actual != int32(testRPS) {
		msgFormat := "RateLimiter.Take() should success %v times, but %v happens"
		t.Fatalf(msgFormat, testRPS, actual)
	}

	// overflow end

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if rateLimiter.Take(ctx) {
		t.Fatal("RateLimiter.Take() should return false on canceled contexts")
	}
}

func TestNoOverflow(t *testing.T) {
	defer detectLeak(t)()

	rateLimiter := ratelimiter.New(bottleneck.NewValve(testRPS, testBurst))

	startTime := time.Now()

	var waitGroup sync.WaitGroup
	for requestNumber := 1; requestNumber <= testRPS+testBurst; requestNumber++ {
		waitGroup.Add(1)

		go func(requestNumber int) {
			defer waitGroup.Done()

			if !rateLimiter.Take(context.Background()) {
				msgFormat := "RateLimiter.Take() should success %v's request out of total %v"
				t.Errorf(msgFormat, requestNumber, testRPS+testBurst)
			}
		}(requestNumber)
	}
	waitGroup.Wait()

	spend := time.Since(startTime)

	if spend <= time.Second {
		msgFormat := "RateLimiter.Take() should block bursted request for more then 1sec but took %v"
		t.Fatalf(msgFormat, spend)
	}

	if spend >= 2*time.Second {
		msgFormat := "RateLimiter.Take() should block bursted request no more then 2 sec but took %v"
		t.Fatalf(msgFormat, spend)
	}
}

func TestOverflow(t *testing.T) {
	defer detectLeak(t)()

	rateLimiter := ratelimiter.New(bottleneck.NewValve(testRPS, testBurst))

	successCount := atomic.Int32{}
	failCount := atomic.Int32{}

	var waitGroup sync.WaitGroup
	for i := 0; i < 2*(testRPS+testBurst); i++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			if rateLimiter.Take(context.Background()) {
				successCount.Add(1)
			} else {
				failCount.Add(1)
			}
		}()
	}
	waitGroup.Wait()

	if actualSuccessed := successCount.Load(); actualSuccessed < int32(testRPS+testBurst) {
		msgFormat := "RateLimiter.Take() should success at least rps's (%v) times, but %v"
		t.Fatalf(msgFormat, testRPS+testBurst, actualSuccessed)
	}

	if actualFailed := failCount.Load(); actualFailed <= 0 {
		msgFormat := "RateLimiter.Take() should fail at least once when we hang out of limits (%v)"
		t.Fatalf(msgFormat, testRPS+testBurst)
	}
}

func detectLeak(t *testing.T) func() {
	t.Helper()

	return func() {
		// let background shutdown function completes
		time.Sleep(time.Second)

		runtime.GC()

		goleak.VerifyNone(t)
	}
}
