package bottleneck_test

import (
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/bottleneck"
)

func TestEqualizer(t *testing.T) {
	bn := bottleneck.NewEqualizer(rps, burst)

	overhead_multiplier := float64(1.1)

	approxTotal := float64(totalRequests) / float64(rps) * float64(time.Second) * overhead_multiplier
	approxAdditional := float64(additionalPool) / float64(rps) * float64(time.Second) * overhead_multiplier
	wrappedTestBottleneck(t, bn, time.Duration(approxTotal), time.Duration(approxAdditional))
}
