package ratelimiter_test

import (
	"context"
	"testing"

	"github.com/alexdyukov/ratelimiter/v2"
	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func BenchmarkRegularBottleneck(b *testing.B) {
	b.StopTimer()

	rateLimiter := ratelimiter.New(bottleneck.NewRegular(overheadTestRPS, overheadTestBurst))

	ctx := context.Background()

	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = rateLimiter.Take(ctx)
		}
	})
}

func BenchmarkValveBottleneck(b *testing.B) {
	b.StopTimer()

	rateLimiter := ratelimiter.New(bottleneck.NewValve(overheadTestRPS, overheadTestBurst))

	ctx := context.Background()

	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = rateLimiter.Take(ctx)
		}
	})
}

func BenchmarkEqualizerBottleneck(b *testing.B) {
	b.StopTimer()

	rateLimiter := ratelimiter.New(bottleneck.NewEqualizer(overheadTestRPS, overheadTestBurst))

	ctx := context.Background()

	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = rateLimiter.Take(ctx)
		}
	})
}
