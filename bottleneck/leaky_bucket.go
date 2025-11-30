package bottleneck

// LeakyBucket typed Bottleneck implements leaky bucket algorithm.
type LeakyBucket struct {
	Regular
}

// NewLeakyBucket returns leaky bucket algo implementation of Bottleneck.
func NewLeakyBucket(size int) (*LeakyBucket, error) {
	bn, err := NewRegular(size, 0)
	if err != nil {
		return nil, err
	}

	ret := LeakyBucket{*bn}

	return &ret, nil
}
