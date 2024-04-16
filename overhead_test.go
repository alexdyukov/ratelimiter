package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/v2"
	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
	"github.com/reugn/equalizer"
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

	rateLimiter := ratelimiter.New(bn)

	ctx := context.Background()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rateLimiter.Take(ctx)
		}
	})
}

func BenchmarkOverheadReugnEqualizerTokenBucket(b *testing.B) {
	b.StopTimer()

	rateLimiter, err := equalizer.NewTokenBucket(overheadTestRPS, time.Second)
	if err != nil {
		b.Error(err)
	}

	ctx := context.Background()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := rateLimiter.Acquire(ctx)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkOverheadReugnEqualizerSlider(b *testing.B) {
	b.StopTimer()

	rateLimiter, err := equalizer.NewSlider(time.Second, overheadTestMinDuration, overheadTestRPS)
	if err != nil {
		b.Error(err)
	}

	ctx := context.Background()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := rateLimiter.Acquire(ctx)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkOverheadRateLimiterValveBottleneck(b *testing.B) {
	b.StopTimer()

	bn := bottleneck.NewValve(overheadTestRPS, overheadTestBurst)

	rateLimiter := ratelimiter.New(bn)

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

	rateLimiter := ratelimiter.New(bn)

	ctx := context.Background()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rateLimiter.Take(ctx)
		}
	})
}
