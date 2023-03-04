package bottleneck_test

import (
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/bottleneck"
)

func TestEqualizer(t *testing.T) {
	bn := bottleneck.NewEqualizer(rps, burst)
	approxTotal := float64(totalRequests) / float64(rps) * float64(time.Second)
	approxAdditional := float64(additionalPool) / float64(rps) * float64(time.Second)
	wrappedTestBottleneck(t, bn, time.Duration(approxTotal), time.Duration(approxAdditional))
}
