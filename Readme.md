# Internet Golf

This is a Go web server built on top of [Caddy](https://caddyserver.com/). Its primary new feature is the high-level admin API it adds to all of its public website deployments. You can update any site that this server serves by making a simple HTTP request to this admin API, which makes it very easy to use in build pipelines.

It does some interesting things, but is not finished yet.

## Starting the Server with Docker

Download the "docker-usage" folder from this repository. From that folder, run `docker compose up -d` to start the server. Then, use `./docker-client.sh [your args here]` (Linux) or `./docker-client.ps1 [your args here]` (Windows) to run Client CLI commands. Start with `./docker-client -h` to see the available commands.

## Deploying Stuff from Github Actions

This section under construction.

## Development

Install Go.

Install dependencies:

```
go get .
```

The build is controlled by [Mage](https://magefile.org/).

- This command makes sure all generated artifacts are up to date (which requires Docker, to run the OpenAPI code generation step), and then builds the result:

```
go tool mage buildwithcodegen
```

- This command just builds what's there:

```
go tool mage build
```

- This command generates just the OpenAPI spec (golf-openapi.yaml) and client SDK (./client-sdk/) without running the full build, which is useful to run whenever the API changes during development (but does also require Docker):

```
go tool mage generateclientsdk
```

## Running Stuff for Development

Run server by itself (during dev, you probably just want it to bind to localhost; also, without this option, it will try to force https for everything, which is hard to make work locally):

```
go run ./cmd --local
```

Run client by itself (-h gives you the available commands):

```
go run ./client-cmd -h
```

## Tests

> [!NOTE]  
> The first time you run tests, you must either do it with sudo/admin rights, or manually add the hosts from `test/utils_test.go` to your system's `hosts` file so that they point to 127.0.0.1. These hosts are used for integration tests.

Run tests:

```
go test ./...
```

Or for more readable output:

```
go test ./... -json | go tool tparse -all
```
