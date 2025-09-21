# Internet Golf

This is a Go web server built on top of [Caddy](https://caddyserver.com/). Its primary new feature is the high-level admin API it adds to all of its public website deployments. You can update any site that this server serves by making a simple HTTP request to this admin API, which makes it very easy to use in build pipelines.

It does some interesting things, but is not finished yet.

## Setup

Install Go.

Install dependencies:

```
go get .
```

## Build

The build is controlled by [Mage](https://magefile.org/). Docker is required to run the OpenAPI code generation step.

```
go tool mage build
```

- This command generates just the OpenAPI spec (golf-openapi.yaml) and client SDK (./client-sdk/) without running the full build:

```
go tool mage generateclientsdk
```

## Usage

Run server by itself (during dev, you probably just want it to bind to localhost):

```
go run ./cmd --local
```

Run client by itself (-h gives you the available commands):

```
go run ./client-cmd -h
```

## Tests

> [!NOTE]  
> The first time you run tests, you must either do it with sudo/the administrator terminal, or manually add the hosts from `test/utils_test.go` to your system's `hosts` file so that they point to 127.0.0.1. These hosts are used for integration tests.

Run tests:

```
go test ./...
```

Or for more readable output:

```
go test ./... -json | go tool tparse -all
```
