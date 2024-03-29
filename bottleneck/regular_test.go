package bottleneck_test

import (
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/bottleneck"
)

func TestRegular(t *testing.T) {
	t.Parallel()

	bn := bottleneck.NewRegular(rps, burst)
	approxTotal := float64(totalRequests/rps) * float64(time.Second)
	approxAdditional := time.Second
	wrappedTestBottleneck(t, bn, time.Duration(approxTotal), approxAdditional)
}
