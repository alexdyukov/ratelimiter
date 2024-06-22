package bottleneck_test

import (
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func TestRegular(t *testing.T) {
	regular := bottleneck.NewRegular(rps, burst)
	approxTotal := float64(totalRequests/rps) * float64(time.Second)
	approxAdditional := time.Second
	wrappedTestBottleneck(t, regular, time.Duration(approxTotal), approxAdditional)
}
