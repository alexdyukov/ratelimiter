# ratelimiter
Go universal rate limiter with easy to use API.
====
[![Go Reference](https://pkg.go.dev/badge/image)](https://pkg.go.dev/github.com/alexdyukov/ratelimiter)
[![Go Report](https://goreportcard.com/badge/github.com/alexdyukov/ratelimiter)](https://goreportcard.com/report/github.com/alexdyukov/ratelimiter)
[![Go Coverage](https://github.com/alexdyukov/ratelimiter/wiki/coverage.svg)](https://raw.githack.com/wiki/alexdyukov/ratelimiter/coverage.html)

Package provides methods to rate limit any type of requests (not only http). Use builtin Bottlenecks or write your own for custom bottle neck logic.

## Compares to most popular solutions

### [go.uber.org/ratelimit](https://pkg.go.dev/go.uber.org/ratelimit)
- `go.uber.org/ratelimit` does not support `burst` option
- `go.uber.org/ratelimit` cannot cancel request or shrink queue

### [x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)
- `x/time/rate` cannot cancel request or shrink queue

### [github.com/reugn/equalizer](https://pkg.go.dev/github.com/reugn/equalizer)
- `github.com/reugn/equalizer` does not support `burst` option

## Benchmarks

Overhead of each Take()/Wait() request with infinity rate:
```go
$ # RegularBottleneck lies on memory usage because uses slice of int64 with length of required RPS
$ go version && go test -bench=. -benchmem -benchtime=10000000x
go version go1.24.2 linux/amd64
goos: linux
goarch: amd64
pkg: github.com/alexdyukov/ratelimiter/v2
cpu: AMD Ryzen 7 8845H w/ Radeon 780M Graphics
BenchmarkRegularBottleneck-16           10000000               271.9 ns/op             0 B/op          0 allocs/op
BenchmarkValveBottleneck-16             10000000               215.9 ns/op             0 B/op          0 allocs/op
BenchmarkEqualizerBottleneck-16         10000000               284.2 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/alexdyukov/ratelimiter/v2    13.769s
```

## Example

```go
package main

import (
        "io/ioutil"
        "net/http"

        "github.com/alexdyukov/ratelimiter/v2"
        "github.com/alexdyukov/ratelimiter/v2/bottleneck"
)

func main() {
        rps := 100
        burst := 25

        bn := bottleneck.NewRegular(rps, burst)
        rl := ratelimiter.New(bn)

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
