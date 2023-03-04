package bottleneck_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	// you cannot sleep less then 1ms for Equalizer on win machines, thats why there is rps = 1000
	rps            int = 1000
	burst          int = 100
	totalRequests  int = 2990
	additionalPool int = 200
)

type testBottleneck interface {
	BreakThrough()
}

func wrappedTestBottleneck(t *testing.T, bn testBottleneck, approxTotal, approxAdditional time.Duration) {
	startTime := time.Now()

	for i := 0; i < totalRequests; i++ {
		bn.BreakThrough()
	}

	spend := time.Since(startTime)
	lower := time.Duration(0.95 * float64(approxTotal))
	higher := time.Duration(1.15 * float64(approxTotal))

	assert.True(t, lower < spend && spend < higher, "%v rps with %v total requests should spend at least %v and no more %v, but spend: %v", rps, totalRequests, lower, higher, spend)

	startTime = time.Now()

	for i := 0; i < additionalPool; i++ {
		bn.BreakThrough()
	}

	spend = time.Since(startTime)
	lower = time.Duration(0.95 * float64(approxAdditional))
	higher = time.Duration(1.15 * float64(approxAdditional))

	assert.True(t, lower < spend && spend < higher, "%v rps with %v total requests should spend at least %v and no more %v, but spend: %v", rps, totalRequests, lower, higher, spend)
}
