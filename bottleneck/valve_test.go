package bottleneck_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func TestValve(t *testing.T) {
	_, err := bottleneck.NewValve(-1, testBurst)
	if err == nil || !errors.Is(err, bottleneck.ErrRPSNegativeOrZero) {
		t.Fatal(err)
	}

	_, err = bottleneck.NewValve(testRPS, -1)
	if err == nil || !errors.Is(err, bottleneck.ErrBurstNegative) {
		t.Fatal(err)
	}

	bn, err := bottleneck.NewValve(testRPS, testBurst)
	if err != nil {
		t.Fatal(err)
	}

	if maxRate := bn.MaxRate(); maxRate != testRPS+testBurst {
		t.Fatalf("invalid maxrate: want %d, but got %d", testRPS+testBurst, maxRate)
	}

	approxTotal := float64(totalRequests/testRPS) * float64(time.Second)
	approxAdditional := float64((additionalPool+totalRequests)/testRPS)*float64(time.Second) - approxTotal
	wrappedTestBottleneck(t, bn, time.Duration(approxTotal), time.Duration(approxAdditional))
}
