// NB: This module name is intentionally not "go get"-able or "go install"-able.
// Users should clone the repo to explore the examples.
module github.com/takekazu/planny

go 1.23.0

toolchain go1.24.1

require (
	connectrpc.com/connect v1.18.1
	connectrpc.com/grpchealth v1.4.0
	connectrpc.com/grpcreflect v1.3.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.39.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/oklog/ulid/v2 v2.1.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
