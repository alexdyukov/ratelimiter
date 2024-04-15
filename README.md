# ratelimiter
Go universal rate limiter with easy to use API.
====
[![GoDoc](https://godoc.org/github.com/alexdyukov/ratelimiter?status.svg)](https://godoc.org/github.com/alexdyukov/ratelimiter)
[![Tests](https://github.com/alexdyukov/ratelimiter/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/alexdyukov/ratelimiter/actions/workflows/tests.yml?query=branch%3Amaster)

Package provides methods to rate limit any type of requests (not only http). Use builtin Bottlenecks or write your own for custom bottle neck logic.

## Compares to most popular solutions

### [go.uber.org/ratelimit](https://pkg.go.dev/go.uber.org/ratelimit)
- `go.uber.org/ratelimit` is not FIFO
- `go.uber.org/ratelimit` does not support `burst` option
- `go.uber.org/ratelimit` cannot cancel request or shrink queue
- `go.uber.org/ratelimit` is faster. See benchmark block for `BenchmarkOverheadUber` and `BenchmarkOverheadRateLimiterEqualizerBottleneck`

### [x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)
- `x/time/rate` cannot cancel request or shrink queue

## Benchmarks

Overhead of each Take()/Wait() request with infinity rate:
```go
$ # RegularBottleneck lies on memory usage because uses slice of time.Time{} with length of required RPS
$ go test -bench=. -benchmem -benchtime=10000000x
goos: linux
goarch: amd64
pkg: github.com/alexdyukov/ratelimiter
cpu: AMD Ryzen 7 5800U with Radeon Graphics
BenchmarkOverheadXTimeRate-16                           10000000               173.0 ns/op             0 B/op          0 allocs/op
BenchmarkOverheadUber-16                                10000000                37.02 ns/op            0 B/op          0 allocs/op
BenchmarkOverheadReugnEqualizerTokenBucket-16           10000000               149.3 ns/op             0 B/op          0 allocs/op
BenchmarkOverheadReugnEqualizerSlider-16                10000000               135.2 ns/op             0 B/op          0 allocs/op
BenchmarkOverheadRateLimiterRegularBottleneck-16        10000000               256.6 ns/op             0 B/op          0 allocs/op
BenchmarkOverheadRateLimiterValveBottleneck-16          10000000               175.1 ns/op             0 B/op          0 allocs/op
BenchmarkOverheadRateLimiterEqualizerBottleneck-16      10000000               255.8 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/alexdyukov/ratelimiter       16.011s
```

## Example

```go
package main

import (
        "io/ioutil"
        "net/http"

        "github.com/alexdyukov/ratelimiter"
        "github.com/alexdyukov/ratelimiter/bottleneck"
)

func main() {
        rps := 100
        burst := 25

        bn := bottleneck.NewRegular(rps, burst)
        rl, shutdown := ratelimiter.New(bn)
        defer shutdown()

        echo := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if !rl.Take(r.Context()) {
                        w.Header().Add("Retry-After", "60")
                        http.Error(w, "Try again later", http.StatusServiceUnavailable)
                        return
                }
                data, _ := ioutil.ReadAll(r.Body)
                w.Write(data)
        })

        http.ListenAndServe(":8080", echo)
}
```

## License

MIT licensed. See the included LICENSE file for details.

