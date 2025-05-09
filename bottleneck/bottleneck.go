// Package bottleneck provides an cheap example of bottlenecks for main ratelimiter package.
// See Algorithms section at https://en.wikipedia.org/wiki/Rate_limiting .
package bottleneck

import "errors"

var (
	// ErrRPSNegativeOrZero is returned by New* when invalid rps parameter.
	ErrRPSNegativeOrZero = errors.New("bottleneck: rps parameter should be greater 0")
	// ErrBurstNegative is returned by New* when invalid burst parameter.
	ErrBurstNegative = errors.New("bottleneck: burst parameter should be greater or equal 0")
)
