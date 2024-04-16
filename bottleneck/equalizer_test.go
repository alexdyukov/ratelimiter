package bottleneck_test

import (
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func TestEqualizer(t *testing.T) {
	bn := bottleneck.NewEqualizer(rps, burst)

	overheadMultiplier := float64(1.1)

	approxTotal := float64(totalRequests) / float64(rps) * float64(time.Second) * overheadMultiplier
	approxAdditional := float64(additionalPool) / float64(rps) * float64(time.Second) * overheadMultiplier
	wrappedTestBottleneck(t, bn, time.Duration(approxTotal), time.Duration(approxAdditional))
}
