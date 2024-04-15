package bottleneck_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testBottleneck interface {
	BreakThrough()
}

const (
	// you cannot sleep less then 1ms for Equalizer on win machines, thats why there is rps = 1000.
	rps            int = 1000
	burst          int = 100
	totalRequests  int = 2990
	additionalPool int = 200
)

func wrappedTestBottleneck(t *testing.T, bn testBottleneck, approxTotal, approxAdditional time.Duration) {
	t.Helper()

	startTime := time.Now()

	for i := 0; i < totalRequests; i++ {
		bn.BreakThrough()
	}

	spend := time.Since(startTime)
	lower := time.Duration(0.95 * float64(approxTotal))
	higher := time.Duration(1.90 * float64(approxTotal))

	msgFormat := "%v rps with %v total requests should spend at least %v and no more %v, but spend: %v"
	assert.True(t, lower < spend && spend < higher, msgFormat, rps, totalRequests, lower, higher, spend)

	startTime = time.Now()

	for i := 0; i < additionalPool; i++ {
		bn.BreakThrough()
	}

	spend = time.Since(startTime)
	lower = time.Duration(0.95 * float64(approxAdditional))
	higher = time.Duration(1.90 * float64(approxAdditional))

	msgFormat = "%v rps with %v total requests should spend at least %v and no more %v, but spend: %v"
	assert.True(t, lower < spend && spend < higher, msgFormat, rps, totalRequests, lower, higher, spend)
}
