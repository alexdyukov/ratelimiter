package ratelimiter_test

import (
	"context"
	"testing"

	"github.com/alexdyukov/ratelimiter/v2"
	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func BenchmarkRegularBottleneck(b *testing.B) {
	b.StopTimer()

	bn, err := bottleneck.NewRegular(overheadTestRPS, overheadTestBurst)
	if err != nil {
		b.Fatal(err)
	}

	rateLimiter := ratelimiter.New(bn)

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

	bn, err := bottleneck.NewValve(overheadTestRPS, overheadTestBurst)
	if err != nil {
		b.Fatal(err)
	}

	rateLimiter := ratelimiter.New(bn)

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

	bn, err := bottleneck.NewEqualizer(overheadTestRPS, overheadTestBurst)
	if err != nil {
		b.Fatal(err)
	}

	rateLimiter := ratelimiter.New(bn)

	ctx := context.Background()

	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = rateLimiter.Take(ctx)
		}
	})
}
