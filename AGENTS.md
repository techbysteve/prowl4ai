# Repository Guidelines

## Project Structure & Module Organization

- `cmd/prowl4ai/main.go`: CLI entrypoint (`crawl` command).
- `internal/browser/`: browser adapter interfaces and Playwright implementation.
- `internal/prowler/`: crawl orchestration/service logic.
- `internal/extract/`: HTML cleaning and Markdown conversion pipeline.
- `internal/config/`: default browser/runtime configuration.
- `internal/model/`, `internal/stderrors/`: shared data models and error helpers.
- `prowl4ai.go`: public library API.
- `plans/`: planning docs; `README.md` and `roadmap.md`: product/context docs.

## Build, Test, and Development Commands

- `go mod download`: install module dependencies.
- `go run github.com/playwright-community/playwright-go/cmd/playwright@latest install chromium`: install Playwright Chromium runtime.
- `go run ./cmd/prowl4ai crawl https://example.com`: run CLI locally.
- `go build ./...`: compile all packages.
- `go test ./...`: run all tests (currently reports no test files).
- `go test -race ./...`: optional race check for concurrent code changes.

## Coding Style & Naming Conventions

- Follow standard Go formatting: run `gofmt -w .` before commits.
- Use idiomatic Go naming: exported `CamelCase`, unexported `camelCase`, package names short/lowercase.
- Keep packages focused by responsibility (browser, extract, config, model).
- Prefer small functions and explicit error wrapping with context.
