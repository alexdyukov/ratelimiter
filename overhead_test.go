package ratelimiter_test

import (
	"context"
	"testing"

	"github.com/alexdyukov/ratelimiter"
	"github.com/alexdyukov/ratelimiter/bottleneck"
	"go.uber.org/ratelimit"
	"golang.org/x/time/rate"
)

func BenchmarkOverheadXTimeRate(b *testing.B) {
	b.StopTimer()

	rateLimiter := rate.NewLimiter(rate.Limit(overheadTestRPS), overheadTestBurst)

	ctx := context.Background()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := rateLimiter.Wait(ctx)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkOverheadUber(b *testing.B) {
	b.StopTimer()

	// no burst option for uber's rate limiter
	rateLimiter := ratelimit.New(overheadTestRPS)

	// no need context, cause uber's rate limiter doesn not support cancels
	// ctx := context.Background()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rateLimiter.Take()
		}
	})
}

func BenchmarkOverheadRateLimiterRegularBottleneck(b *testing.B) {
	b.StopTimer()

	bn := bottleneck.NewRegular(overheadTestRPS, overheadTestBurst)

	rateLimiter, shutdown := ratelimiter.New(bn)
	defer shutdown()

	ctx := context.Background()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rateLimiter.Take(ctx)
		}
	})
}

func BenchmarkOverheadRateLimiterValveBottleneck(b *testing.B) {
	b.StopTimer()

	bn := bottleneck.NewValve(overheadTestRPS, overheadTestBurst)

	rateLimiter, shutdown := ratelimiter.New(bn)
	defer shutdown()

	ctx := context.Background()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rateLimiter.Take(ctx)
		}
	})
}

func BenchmarkOverheadRateLimiterEqualizerBottleneck(b *testing.B) {
	b.StopTimer()

	bn := bottleneck.NewEqualizer(overheadTestRPS, overheadTestBurst)

	rateLimiter, shutdown := ratelimiter.New(bn)
	defer shutdown()

	ctx := context.Background()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rateLimiter.Take(ctx)
		}
	})
}
