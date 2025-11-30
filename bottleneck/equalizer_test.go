package bottleneck_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func TestEqualizer(t *testing.T) {
	_, err := bottleneck.NewEqualizer(-1, testBurst)
	if err == nil || !errors.Is(err, bottleneck.ErrRPSNegativeOrZero) {
		t.Fatal(err)
	}

	_, err = bottleneck.NewEqualizer(testRPS, -1)
	if err == nil || !errors.Is(err, bottleneck.ErrBurstNegative) {
		t.Fatal(err)
	}

	bn, err := bottleneck.NewEqualizer(testRPS, testBurst)
	if err != nil {
		t.Fatal(err)
	}

	if maxRate := bn.MaxRate(); maxRate != testRPS+testBurst {
		t.Fatalf("invalid maxrate: want %d, but got %d", testRPS+testBurst, maxRate)
	}

	overheadMultiplier := float64(1.1)

	approxTotal := float64(totalRequests) / float64(testRPS) * float64(time.Second) * overheadMultiplier
	approxAdditional := float64(additionalPool) / float64(testRPS) * float64(time.Second) * overheadMultiplier
	wrappedTestBottleneck(t, bn, time.Duration(approxTotal), time.Duration(approxAdditional))
}
