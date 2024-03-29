## This is a microservice written in Golang: maillingList

- it demonstrates usage of
  - JSON API Server,
  - gRPC Server/Client,
  - Protocol buffers
  - SQLite
  - Goroutines
  - CLI

# Setup

This project requires a `gcc` compiler installed and the `protobuf` code generation tools.

## Install protobuf compiler

Install the `protoc` tool using the instructions available at [https://grpc.io/docs/protoc-installation/](https://grpc.io/docs/protoc-installation/).

Alternatively you can download a pre-built binary from [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases) and placing the extracted binary somewhere in your `$PATH`.

## Install Go protobuf codegen tools

`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

## Generate Go code from .proto files

```
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/mail.proto
```

### How to start

- `go mod tidy`
- `go run ./server` - to start the server
- `go run ./client` - to execute the client

### Configuration

- you may set the environmental variables:
  - `MAILINGLIST_DB=` - filepath of a sqlite database (default is **_`./_data/list.db`_**)
  - `MAILINGLIST_BIND_JSON=` - address and port to listen (default is **_`localhost:3000`_**)
  - `MAILINGLIST_BIND_GRPC=` - address and port to listen (default is **_`localhost:3001`_**)

### Kudos

- The app's idea came from [Golang course, Zero To Mastery](https://academy.zerotomastery.io/courses/1600953/lectures/38731793)
