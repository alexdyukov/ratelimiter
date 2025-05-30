package bottleneck_test

import (
	"testing"
	"time"

	"github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func TestValve(t *testing.T) {
	bn, err := bottleneck.NewValve(rps, burst)
	if err != nil {
		t.Fatal(err)
	}

	approxTotal := float64(totalRequests/rps) * float64(time.Second)
	approxAdditional := float64((additionalPool+totalRequests)/rps)*float64(time.Second) - approxTotal
	wrappedTestBottleneck(t, bn, time.Duration(approxTotal), time.Duration(approxAdditional))
}
