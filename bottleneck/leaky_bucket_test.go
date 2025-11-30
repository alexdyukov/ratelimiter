package bottleneck_test

import (
	"errors"
	"testing"

	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func TestLeakyBucket(t *testing.T) {
	_, err := bottleneck.NewLeakyBucket(-1)
	if err == nil || !errors.Is(err, bottleneck.ErrRPSNegativeOrZero) {
		t.Fatal(err)
	}

	bn, err := bottleneck.NewLeakyBucket(testRPS)
	if err != nil {
		t.Fatal(err)
	}

	if maxRate := bn.MaxRate(); maxRate != testRPS {
		t.Fatalf("invalid maxrate: want %d, but got %d", testRPS+testBurst, maxRate)
	}
}
