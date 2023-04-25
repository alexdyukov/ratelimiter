module github.com/alexdyukov/ratelimiter

go 1.18

replace github.com/alexdyukov/ratelimiter => ./

retract (
	v1.1.0 // invalid implementation
	v1.0.5 // invalid implementation
	v1.0.4 // invalid implementation
	v1.0.0 // invalid implementation
)

require (
	github.com/stretchr/testify v1.8.2
	go.uber.org/ratelimit v0.2.0
	golang.org/x/time v0.3.0
)

require (
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
