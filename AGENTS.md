# Mate Go libraries

This module collects small reusable Go libraries for Mate projects. Its module
path is `github.com/MateWorkspace/mate-golang`. Each library lives in a public
package directly below the module root; consumers import only the packages they
need. Keep one root `go.mod` and `go.sum` while the libraries share a release
cycle. There is no service entrypoint or application layer here.

## Layout

- `pgxdt/`: PostgreSQL execution and transaction helpers using pgx v5.
- `go.mod` and `go.sum`: module identity and reproducible dependency checksums.

Add future libraries as sibling packages. Use `internal/` only for helpers that
must remain private to this module. Do not add a `pkg/` or `libraries/` wrapper:
the package directories are already public Go import paths.

## Style and checks

Use short lowercase package names, explicit constructors, small interfaces,
context-first I/O methods, and early returns. Keep dependencies limited to what
the libraries use. Avoid service-specific models, logging, configuration, or
domain contracts. Document exported APIs with Go comments and format with
`gofmt`. Preserve caller contexts and underlying errors.

Run from this directory:

```sh
go build ./...
go test ./...
go vet ./...
```

Consumers import `github.com/MateWorkspace/mate-golang/pgxdt`. Until the module
is published, a consumer can use a local module replacement:

```sh
go mod edit -require=github.com/MateWorkspace/mate-golang@v0.0.0
go mod edit -replace=github.com/MateWorkspace/mate-golang=../mate-golang
go mod tidy
```

Add the import before running `go mod tidy`; unused requirements are removed.
Once published, consumers should depend on a tagged module version. Release
tags apply to the whole collection, such as `v0.1.0`.
