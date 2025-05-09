package bottleneck_test

import (
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func TestRegular(t *testing.T) {
	bn, err := bottleneck.NewRegular(rps, burst)
	if err != nil {
		t.Fatal(err)
	}

	approxTotal := float64(totalRequests/rps) * float64(time.Second)
	approxAdditional := time.Second
	wrappedTestBottleneck(t, bn, time.Duration(approxTotal), approxAdditional)
}
