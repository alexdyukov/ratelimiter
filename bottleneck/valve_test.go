package bottleneck_test

import (
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/bottleneck"
)

func TestValve(t *testing.T) {
	t.Parallel()

	bn := bottleneck.NewValve(rps, burst)
	approxTotal := float64(totalRequests/rps) * float64(time.Second)
	approxAdditional := float64((additionalPool+totalRequests)/rps)*float64(time.Second) - approxTotal
	wrappedTestBottleneck(t, bn, time.Duration(approxTotal), time.Duration(approxAdditional))
}
