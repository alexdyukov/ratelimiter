package bottleneck

// TokenBucket typed Bottleneck similar to ticker queue.
// Provides bottle neck which pass 1 and no more then each {time.Second/RPS} second.
// RPS 1000 means pass 1 each 1 millisecond.
type TokenBucket struct {
	Equalizer
}

// NewTokenBucket returns token bucket algo implementation of Bottleneck.
func NewTokenBucket(rps, burst int) (*TokenBucket, error) {
	bn, err := NewEqualizer(rps, burst)
	if err != nil {
		return nil, err
	}

	ret := TokenBucket{*bn}

	return &ret, nil
}
