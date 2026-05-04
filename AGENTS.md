# XLog - Agent Instructions

Static site generator for digital gardening written in Go. Serves and edits markdown files with extension system.

## Architecture

- **Main binary**: `cmd/xlog/xlog.go` (imports all extensions via `extensions/all`)
- **Extensions**: 37 extensions in `extensions/` loaded via blank imports
- **Vendored code**: `markdown/` contains vendored goldmark parser (excluded from linting)
- **Entry point**: `xlog.Start(context.Background())` in main

To add/remove extensions, modify `extensions/all/all.go` import list.

## Development Commands

```bash
# Run tests (standard)
go test ./...

# Run tests with race detector (CI also runs this)
go test -race ./...

# Lint (exact CI command)
golangci-lint run

# Run the server locally
go run ./cmd/xlog

# Build binary
go build -o xlog ./cmd/xlog

# Build static site
go run ./cmd/xlog -build <output-dir>
```

## Testing Conventions

- Use table-driven tests for multiple cases (see `page_test.go:19-36` for example)
- Tests that need filesystem state: use `t.TempDir()` and `os.Chdir()` pattern
- Always defer restoration of working directory in tests that change it
- Helper pattern: `setupTestEnv(t)` returns cleanup func (see extension tests)

## Linting Rules

- `.golangci.yml` excludes `markdown/` directory (vendored goldmark)
- Tests excluded from: `gosec`, `errcheck`
- Min cyclomatic complexity: 15
- CI runs `golangci-lint-action@v6` with latest version

## Key Constraints

- Go 1.25 minimum (go.mod)
- CI uses Go 1.24 (will need update to match go.mod)
- Main branch: `master`
- Docker: builds to `ghcr.io/emad-elsaid/xlog`, serves on port 3000, mounts `~/.xlog:/files`

## CLI Flags

Most extensions add their own flags. Check with `go run ./cmd/xlog -h`.
Notable core flags:
- `-bind`: server address (default `127.0.0.1:3000`)
- `-build <dir>`: generate static site
- `-source <dir>`: markdown files location (default `.`)

## Release Process

- Tag `v*` triggers multi-platform binary builds (Linux/Windows/Darwin, 386/amd64/arm64)
- Binaries built from `./cmd/xlog`
- Docker image pushed to GitHub Container Registry
