package bottleneck_test

import (
	"errors"
	"testing"

	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func TestTokenBucket(t *testing.T) {
	_, err := bottleneck.NewTokenBucket(-1, testBurst)
	if err == nil || !errors.Is(err, bottleneck.ErrRPSNegativeOrZero) {
		t.Fatal(err)
	}

	_, err = bottleneck.NewTokenBucket(testRPS, -1)
	if err == nil || !errors.Is(err, bottleneck.ErrBurstNegative) {
		t.Fatal(err)
	}

	bn, err := bottleneck.NewTokenBucket(testRPS, testBurst)
	if err != nil {
		t.Fatal(err)
	}

	if maxRate := bn.MaxRate(); maxRate != testRPS+testBurst {
		t.Fatalf("invalid maxrate: want %d, but got %d", testRPS+testBurst, maxRate)
	}
}
