module github.com/alexdyukov/ratelimiter/v2

go 1.19

replace github.com/alexdyukov/ratelimiter/v2 => ./

require (
	github.com/reugn/equalizer v0.2.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/goleak v1.3.0
	go.uber.org/ratelimit v0.3.0
	golang.org/x/time v0.3.0
)

require (
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
