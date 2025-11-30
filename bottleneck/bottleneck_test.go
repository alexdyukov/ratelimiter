package bottleneck_test

import (
	"testing"
	"time"
)

type testBottleneck interface {
	BreakThrough()
}

const (
	// you cannot sleep less then 1ms for Equalizer on win machines, thats why there is rps = 1000.
	testRPS        int = 1000
	testBurst      int = 100
	totalRequests  int = 2990
	additionalPool int = 200
)

func wrappedTestBottleneck(t *testing.T, bottleneck testBottleneck, approxTotal, approxAdditional time.Duration) {
	t.Helper()

	startTime := time.Now()

	for i := 0; i < totalRequests; i++ {
		bottleneck.BreakThrough()
	}

	spend := time.Since(startTime)
	lower := time.Duration(0.95 * float64(approxTotal))
	higher := time.Duration(1.90 * float64(approxTotal))

	if lower >= spend || spend >= higher {
		msgFormat := "main pool: %v total requests should spend at least %v and no more %v, but spend: %v"
		t.Fatalf(msgFormat, totalRequests, lower, higher, spend)
	}

	startTime = time.Now()

	for i := 0; i < additionalPool; i++ {
		bottleneck.BreakThrough()
	}

	spend = time.Since(startTime)
	lower = time.Duration(0.95 * float64(approxAdditional))
	higher = time.Duration(1.90 * float64(approxAdditional))

	if lower >= spend || spend >= higher {
		msgFormat := "additional pool: %v total requests should spend at least %v and no more %v, but spend: %v"
		t.Fatalf(msgFormat, totalRequests, lower, higher, spend)
	}
}
